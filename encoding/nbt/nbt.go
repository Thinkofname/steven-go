package nbt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// TODO(Think) Cleanup

var (
	ErrInvalidCompound = errors.New("invalid compound")
)

type Compound struct {
	Name  string
	items map[string]interface{}
}

func NewCompound() *Compound {
	return &Compound{
		items: make(map[string]interface{}),
	}
}

func (c *Compound) Serialize(w io.Writer) error {
	panic("TODO NBT Serialize")
	return nil
}

func (c *Compound) Deserialize(r io.Reader) error {
	id, err := readByte(r)
	if err != nil {
		return err
	}
	if id != 10 {
		return ErrInvalidCompound
	}
	c.Name, err = readString(r)
	if err != nil {
		return err
	}
	for {
		id, err := readByte(r)
		if err != nil {
			return err
		}
		// End of compound
		if id == 0 {
			break
		}
		name, err := readString(r)
		if err != nil {
			return err
		}
		c.items[name], err = readType(r, int(id))
		if err != nil {
			return err
		}
	}
	return nil
}

type List struct {
	Type     int
	Elements []interface{}
}

func (l *List) Deserialize(r io.Reader) error {
	t, err := readByte(r)
	if err != nil {
		return err
	}
	l.Type = int(t)
	var le int32
	err = binary.Read(r, binary.BigEndian, &le)
	if err != nil {
		return err
	}
	l.Elements = make([]interface{}, le)
	for i := 0; i < int(le); i++ {
		l.Elements[i], err = readType(r, l.Type)
		if err != nil {
			return err
		}
	}
	return nil
}

func readType(r io.Reader, id int) (interface{}, error) {
	switch id {
	case 1:
		return readByte(r)
	case 2:
		var v int16
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 3:
		var v int32
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 4:
		var v int64
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 5:
		var v float32
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 6:
		var v float64
		err := binary.Read(r, binary.BigEndian, &v)
		return v, err
	case 7:
		var l int32
		err := binary.Read(r, binary.BigEndian, &l)
		if err != nil {
			return nil, err
		}
		v := make([]byte, l)
		_, err = io.ReadFull(r, v)
		return v, err

	case 8:
		return readString(r)
	case 9:
		l := &List{}
		err := l.Deserialize(r)
		return l, err
	case 10:
		c := NewCompound()
		err := c.Deserialize(r)
		return c, err
	case 11:
		var l int32
		err := binary.Read(r, binary.BigEndian, &l)
		if err != nil {
			return nil, err
		}
		v := make([]int32, l)
		err = binary.Read(r, binary.BigEndian, v)
		return v, err
	}
	return nil, fmt.Errorf("invalid type %d", id)
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

func writeString(w io.Writer, str string) error {
	b := []byte(str)
	err := binary.Write(w, binary.BigEndian, int16(len(b)))
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func readString(r io.Reader) (string, error) {
	var l int16
	err := binary.Read(r, binary.BigEndian, &l)
	if err != nil {
		return "", nil
	}
	buf := make([]byte, int(l))
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}
