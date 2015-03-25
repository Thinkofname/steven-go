package protocol

import (
	"errors"
	"io"

	"github.com/thinkofdeath/steven/encoding/nbt"
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

func writeVarInt(w io.Writer, i VarInt) error {
	ui := uint32(i)
	for {
		if (ui & ^varPart) == 0 {
			err := writeByte(w, byte(ui))
			return err
		}
		err := writeByte(w, byte((ui&varPart)|0x80))
		if err != nil {
			return err
		}
		ui >>= 7
	}
}

func readVarInt(r io.Reader) (VarInt, error) {
	var size uint
	var val uint32
	for {
		b, err := readByte(r)
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

func writeVarLong(w io.Writer, i VarLong) error {
	ui := uint64(i)
	for {
		if (ui & ^varPartLong) == 0 {
			err := writeByte(w, byte(ui))
			return err
		}
		err := writeByte(w, byte((ui&varPartLong)|0x80))
		if err != nil {
			return err
		}
		ui >>= 7
	}
}

func readVarLong(r io.Reader) (VarLong, error) {
	var size uint
	var val uint64
	for {
		b, err := readByte(r)
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

func writeString(w io.Writer, str string) error {
	b := []byte(str)
	err := writeVarInt(w, VarInt(len(b)))
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func readString(r io.Reader) (string, error) {
	l, err := readVarInt(r)
	if err != nil {
		return "", nil
	}
	buf := make([]byte, int(l))
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}

func writeBool(w io.Writer, b bool) error {
	if b {
		return writeByte(w, 1)
	}
	return writeByte(w, 0)
}

func readBool(r io.Reader) (bool, error) {
	b, err := readByte(r)
	if b == 0 {
		return false, err
	}
	return true, err
}

func writeByte(w io.Writer, b byte) error {
	if bw, ok := w.(io.ByteWriter); ok {
		return bw.WriteByte(b)
	}
	var buf [1]byte
	buf[0] = b
	_, err := w.Write(buf[:1])
	return err
}

func readByte(r io.Reader) (byte, error) {
	if br, ok := r.(io.ByteReader); ok {
		return br.ReadByte()
	}
	var buf [1]byte
	_, err := r.Read(buf[:1])
	return buf[0], err
}

func readNBT(r io.Reader) (*nbt.Compound, error) {
	b, err := readByte(r)
	if err != nil || b == 0 { // 0 == No tag
		return nil, err
	}
	n := nbt.NewCompound()
	err = n.Deserialize(r)
	return n, err
}

func writeNBT(w io.Writer, n *nbt.Compound) error {
	if n == nil {
		return writeByte(w, 0)
	}
	return n.Serialize(w)
}
