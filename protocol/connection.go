// Package protocol provides definitions for Minecraft packets as well as
// methods for reading and writing them
package protocol

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// Conn is a connection from or to a Minecraft client.
//
// The Minecraft protocol as multiple states that it
// switches between during login/status pinging, the
// state may be set using the State field.
type Conn struct {
	conn      byteConn
	net       net.Conn
	direction int
	State     State

	host string
	port uint16
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
		conn: byteConn{
			ReadWriter: c,
		},
		net:       c,
		direction: serverbound,
		host:      parts[0],
		port:      uint16(port),
	}, nil
}

// BUG(Think) Compression and encryption is missing

// WritePacket serializes the packet to the underlying
// connection, optionally encrypting and/or compressing
func (c *Conn) WritePacket(packet Packet) error {
	buf := &bytes.Buffer{}

	// Contents of the packet (ID + Data)
	if err := writeVarInt(buf, VarInt(packet.id())); err != nil {
		return err
	}
	if err := packet.write(buf); err != nil {
		return err
	}

	// Write the length prefix followed by the buffer
	if err := writeVarInt(c.conn, VarInt(buf.Len())); err != nil {
		return err
	}
	_, err := buf.WriteTo(c.conn)
	return err
}

// ReadPacket deserializes a packet from the underlying
// connection, optionally decrypting and/or decompressing
func (c *Conn) ReadPacket() (Packet, error) {
	// Length prefix
	size, err := readVarInt(c.conn)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(c.conn, buf); err != nil {
		return nil, err
	}

	r := bytes.NewReader(buf)
	// Packet ID
	id, err := readVarInt(r)
	if err != nil {
		return nil, err
	}
	// Direction is swapped as this is coming from the other way
	packets := packetCreator[c.State][(c.direction+1)&1]
	if id < 0 || int(id) > len(packets) || packets[id] == nil {
		return nil, fmt.Errorf("Unknown packet %02X", id)
	}
	packet := packets[id]()
	if err := packet.read(r); err != nil {
		return packet, err
	}
	// If we haven't fully read the whole buffer then something went wrong.
	// Mostly likely our packet definitions are out of date or incorrect
	if r.Len() != 0 {
		return packet, fmt.Errorf("Have %d byte(s) left to read", r.Len())
	}
	return packet, nil
}

func (c *Conn) Close() error {
	return c.net.Close()
}

// Provides helper byte reading/writing methods
type byteConn struct {
	io.ReadWriter
}

func (c byteConn) WriteByte(b byte) error {
	var buf [1]byte
	buf[0] = b
	_, err := c.Write(buf[:1])
	return err
}

func (c byteConn) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := c.Read(buf[:1])
	return buf[0], err
}
