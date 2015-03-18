package protocol

import (
	"bytes"
	"errors"
	"io"
)

const varPart = uint32(0x7F)

var (
	// ErrVarIntTooLarge is returned when a read varint was too large
	// (more than 5 bytes)
	ErrVarIntTooLarge = errors.New("VarInt too large")
)

func writeVarInt(w io.ByteWriter, i VarInt) error {
	ui := uint32(i)
	for {
		if (ui & ^varPart) == 0 {
			err := w.WriteByte(byte(ui))
			return err
		}
		err := w.WriteByte(byte((ui & varPart) | 0x80))
		if err != nil {
			return err
		}
		ui >>= 7
	}
}

func readVarInt(r io.ByteReader) (VarInt, error) {
	var size uint
	var val uint32
	for {
		b, err := r.ReadByte()
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

func writeString(w *bytes.Buffer, str string) error {
	b := []byte(str)
	err := writeVarInt(w, VarInt(len(b)))
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func readString(r *bytes.Reader) (string, error) {
	l, err := readVarInt(r)
	if err != nil {
		return "", nil
	}
	buf := make([]byte, int(l))
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}

func writeBool(w *bytes.Buffer, b bool) error {
	if b {
		return w.WriteByte(1)
	}
	return w.WriteByte(0)
}

func readBool(r *bytes.Reader) (bool, error) {
	b, err := r.ReadByte()
	if b == 0 {
		return false, err
	}
	return true, err
}
