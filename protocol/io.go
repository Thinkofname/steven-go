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

package protocol

import (
	"errors"
	"io"

	"github.com/thinkofdeath/phteven/encoding/nbt"
)

const varPart = uint32(0x7F)
const varPartLong = uint64(0x7F)

var (
	// ErrVarIntTooLarge is returned when a read varint was too large
	// (more than 5 bytes)
	ErrVarIntTooLarge = errors.New("VarInt too large")
	// ErrVarLongTooLarge is returned when a read varint was too large
	// (more than 10 bytes)
	ErrVarLongTooLarge = errors.New("VarLong too large")
)

func varIntSize(i VarInt) int {
	size := 0
	ui := uint32(i)
	for {
		size++
		if (ui & ^varPart) == 0 {
			return size
		}
		ui >>= 7
	}

}

// WriteVarInt encodes the passed VarInt into the writer.
func WriteVarInt(w io.Writer, i VarInt) error {
	ui := uint32(i)
	for {
		if (ui & ^varPart) == 0 {
			err := WriteByte(w, byte(ui))
			return err
		}
		err := WriteByte(w, byte((ui&varPart)|0x80))
		if err != nil {
			return err
		}
		ui >>= 7
	}
}

// ReadVarInt reads a VarInt encoded integer from the reader.
func ReadVarInt(r io.Reader) (VarInt, error) {
	var size uint
	var val uint32
	for {
		b, err := ReadByte(r)
		if err != nil {
			return VarInt(val), err
		}

		val |= (uint32(b) & varPart) << (size * 7)
		size++
		if size > 5 {
			return VarInt(val), ErrVarIntTooLarge
		}

		if (b & 0x80) == 0 {
			break
		}
	}
	return VarInt(val), nil
}

// WriteVarLong encodes the passed VarLong into the writer.
func WriteVarLong(w io.Writer, i VarLong) error {
	ui := uint64(i)
	for {
		if (ui & ^varPartLong) == 0 {
			err := WriteByte(w, byte(ui))
			return err
		}
		err := WriteByte(w, byte((ui&varPartLong)|0x80))
		if err != nil {
			return err
		}
		ui >>= 7
	}
}

// ReadVarLong reads a VarLong encoded 64 bit integer from the reader.
func ReadVarLong(r io.Reader) (VarLong, error) {
	var size uint
	var val uint64
	for {
		b, err := ReadByte(r)
		if err != nil {
			return VarLong(val), err
		}

		val |= (uint64(b) & varPartLong) << (size * 7)
		size++
		if size > 10 {
			return VarLong(val), ErrVarLongTooLarge
		}

		if (b & 0x80) == 0 {
			break
		}
	}
	return VarLong(val), nil
}

// WriteString writes a VarInt prefixed utf-8 string to the
// writer.
func WriteString(w io.Writer, str string) error {
	b := []byte(str)
	err := WriteVarInt(w, VarInt(len(b)))
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// ReadString reads a VarInt prefixed utf-8 string to the
// reader.
func ReadString(r io.Reader) (string, error) {
	l, err := ReadVarInt(r)
	if err != nil {
		return "", nil
	}
	buf := make([]byte, int(l))
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}

// WriteBool writes a bool to the writer as a single byte.
func WriteBool(w io.Writer, b bool) error {
	if b {
		return WriteByte(w, 1)
	}
	return WriteByte(w, 0)
}

// ReadBool reads a single byte from the reader as a bool.
func ReadBool(r io.Reader) (bool, error) {
	b, err := ReadByte(r)
	if b == 0 {
		return false, err
	}
	return true, err
}

// WriteByte writes a single byte to the writer. If the
// Writer is a ByteWriter then that will be used instead.
func WriteByte(w io.Writer, b byte) error {
	if bw, ok := w.(io.ByteWriter); ok {
		return bw.WriteByte(b)
	}
	var buf [1]byte
	buf[0] = b
	_, err := w.Write(buf[:1])
	return err
}

// ReadByte reads a single byte from the Reader. If the
// Reader is a ByteReader then that will be used instead.
func ReadByte(r io.Reader) (byte, error) {
	if br, ok := r.(io.ByteReader); ok {
		return br.ReadByte()
	}
	var buf [1]byte
	_, err := r.Read(buf[:1])
	return buf[0], err
}

// ReadNBT reads an nbt tag from the reader.
// Returns nil if there is no tag.
func ReadNBT(r io.Reader) (*nbt.Compound, error) {
	b, err := ReadByte(r)
	if err != nil || b == 0 { // 0 == No tag
		return nil, err
	}
	n := nbt.NewCompound()
	err = n.Deserialize(r)
	return n, err
}

// WriteNBT writes an nbt tag to the wrtier.
// nil can be used to specify that there isn't a tag.
func WriteNBT(w io.Writer, n *nbt.Compound) error {
	if n == nil {
		return WriteByte(w, 0)
	}
	return n.Serialize(w)
}
