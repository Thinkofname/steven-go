// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package protocol provides definitions for Minecraft packets as well as
// methods for reading and writing them
package protocol

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

// Conn is a connection from or to a Minecraft client.
//
// The Minecraft protocol as multiple states that it
// switches between during login/status pinging, the
// state may be set using the State field.
type Conn struct {
	r                    io.Reader
	w                    io.Writer
	net                  net.Conn
	direction            int
	State                State
	compressionThreshold int

	host string
	port uint16

	zlibReader io.ReadCloser
	zlibWriter *zlib.Writer
}

// Dial creates a connection to a Minecraft server at
// the passed address. The address is in same format as
// the vanilla client takes it
//     host:port
//     // or
//     host
// If the port isn't provided then a SRV lookup is
// performed and if successful it will continue
// connecting using the returned address. If the lookup
// fails then the port is assumed to be 25565.
func Dial(address string) (*Conn, error) {
	if !strings.ContainsRune(address, ':') {
		// Attempt a srv lookup first (like vanilla)
		_, srvs, err := net.LookupSRV("minecraft", "tcp", address)
		if err == nil {
			address = fmt.Sprintf("%s:%d", srvs[0].Target, srvs[0].Port)
		} else {
			// Fallback to the default port
			address = address + ":25565"
		}
	}
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(address, ":", 2)
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	return &Conn{
		r:                    c,
		w:                    c,
		net:                  c,
		direction:            serverbound,
		host:                 parts[0],
		port:                 uint16(port),
		compressionThreshold: -1,
	}, nil
}

// WritePacket serializes the packet to the underlying
// connection, optionally encrypting and/or compressing
func (c *Conn) WritePacket(packet Packet) error {
	// 15 second timeout
	c.net.SetWriteDeadline(time.Now().Add(15 * time.Second))

	buf := &bytes.Buffer{}

	// Contents of the packet (ID + Data)
	if err := writeVarInt(buf, VarInt(packet.id())); err != nil {
		return err
	}
	if err := packet.write(buf); err != nil {
		return err
	}

	uncompessedSize := 0
	extra := 0
	// Only compress if compression is enabled and the packet is large enough
	if c.compressionThreshold >= 0 && buf.Len() > c.compressionThreshold {
		var err error
		nBuf := &bytes.Buffer{}
		if c.zlibWriter == nil {
			c.zlibWriter, _ = zlib.NewWriterLevel(nBuf, zlib.BestSpeed)
		} else {
			c.zlibWriter.Reset(nBuf)
		}
		uncompessedSize = buf.Len()

		if _, err = buf.WriteTo(c.zlibWriter); err != nil {
			return err
		}
		if err = c.zlibWriter.Close(); err != nil {
			return err
		}
		buf = nBuf
	}

	// Account for the compression header if enabled
	if c.compressionThreshold >= 0 {
		extra = varIntSize(VarInt(uncompessedSize))
	}

	// Write the length prefix followed by the buffer
	if err := writeVarInt(c.w, VarInt(buf.Len()+extra)); err != nil {
		return err
	}

	// Write the uncompressed packet size
	if c.compressionThreshold >= 0 {
		if err := writeVarInt(c.w, VarInt(uncompessedSize)); err != nil {
			return err
		}
	}

	_, err := buf.WriteTo(c.w)
	return err
}

// ReadPacket deserializes a packet from the underlying
// connection, optionally decrypting and/or decompressing
func (c *Conn) ReadPacket() (Packet, error) {
	// 15 second timeout
	c.net.SetReadDeadline(time.Now().Add(15 * time.Second))
	// Length prefix
	size, err := readVarInt(c.r)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(c.r, buf); err != nil {
		return nil, err
	}

	var r *bytes.Reader
	r = bytes.NewReader(buf)

	// If compression is enabled then we may need to decompress the packet
	if c.compressionThreshold >= 0 {
		// With compression enabled an extra length prefix is added
		// which is the length of the packet when uncompressed.
		uncompSize, err := readVarInt(r)
		if err != nil {
			return nil, err
		}
		// A uncompressed size of 0 means the packet wasn't compressed
		// and when can continue normally.
		if uncompSize != 0 {
			// Reuse the old reader to save on allocations
			if c.zlibReader == nil {
				c.zlibReader, err = zlib.NewReader(r)
				if err != nil {
					return nil, err
				}
			} else {
				err = c.zlibReader.(zlib.Resetter).Reset(r, nil)
				if err != nil {
					return nil, err
				}
			}

			// Read the whole packet at once instead of in tiny steps
			data := make([]byte, uncompSize)
			_, err := io.ReadFull(c.zlibReader, data)
			if err != nil {
				return nil, err
			}
			r = bytes.NewReader(data)
		}
	}

	// Packet ID
	id, err := readVarInt(r)
	if err != nil {
		return nil, err
	}
	// Direction is swapped as this is coming from the other way
	packets := packetCreator[c.State][(c.direction+1)&1]
	if id < 0 || int(id) > len(packets) || packets[id] == nil {
		return nil, fmt.Errorf("Unknown packet %s:%02X", c.State, id)
	}
	packet := packets[id]()
	if err := packet.read(r); err != nil {
		return packet, fmt.Errorf("packet(%s:%02X): %s", c.State, id, err)
	}
	// If we haven't fully read the whole buffer then something went wrong.
	// Mostly likely our packet definitions are out of date or incorrect
	if r.Len() > 0 {
		return packet, fmt.Errorf("Didn't finish reading packet %s:%02X, have %d bytes left", c.State, id, r.Len())
	}
	return packet, nil
}

// EnableEncryption enables cfb8 encryption on the protocol using the passed
// key.
func (c *Conn) EnableEncryption(key []byte) error {
	cip, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	c.r = cipher.StreamReader{
		R: c.net,
		S: newCFB8(cip, key, true),
	}

	c.w = cipher.StreamWriter{
		W: c.net,
		S: newCFB8(cip, key, false),
	}
	return nil
}

// SetCompression changes the threshold at which packets are compressed.
func (c *Conn) SetCompression(threshold int) {
	c.compressionThreshold = threshold
}

// Close closes the underlying connection
func (c *Conn) Close() error {
	return c.net.Close()
}
