package protocol

import (
	"encoding/binary"
	"errors"
	"io"
)

// Metadata is a simple index -> value map used in the Minecraft protocol.
// A limited number of types are supported:
//     int8
//     int16
//     int32
//     float32
//     string
//     ItemStack
//     []int32
//     []float32
type Metadata map[int]interface{}

func readMetadata(r io.Reader) (Metadata, error) {
	m := make(Metadata)
	for {
		b, err := readByte(r)
		if err != nil || b == 0x7F {
			return m, err
		}
		index := int(b & 0x1F)
		t := b >> 5

		switch t {
		case 0:
			var val int8
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		case 1:
			var val int16
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		case 2:
			var val int32
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		case 3:
			var val float32
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		case 4:
			m[index], err = readString(r)
		case 5:
			i := ItemStack{}
			err = i.Deserialize(r)
			m[index] = i
		case 6:
			var val [3]int32
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		case 7:
			var val [3]float32
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		default:
			err = errors.New("invalid metadata type")
		}
		if err != nil {
			return m, err
		}
	}
}

func writeMetadata(w io.Writer, m Metadata) error {
	for index, v := range m {
		t := 0

		var val interface{} = v
		switch v.(type) {
		case int8:
			t = 0
		case int16:
			t = 1
		case int32:
			t = 2
		case float32:
			t = 3
		case string:
			t = 4
			val = nil
		case ItemStack:
			t = 5
			val = nil
		case [3]int32:
			t = 6
		case [3]float32:
			t = 7
		default:
			return errors.New("invalid metadata type")
		}
		if err := writeByte(w, byte(index)|(byte(t)<<5)); err != nil {
			return err
		}
		var err error
		if val != nil {
			err = binary.Write(w, binary.BigEndian, val)
		} else {
			switch v := v.(type) {
			case string:
				err = writeString(w, v)
			case ItemStack:
				err = v.Serialize(w)
			}
		}
		if err != nil {
			return err
		}
	}
	return writeByte(w, 0x7F)
}
