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
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/thinkofdeath/steven/format"
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
//     VarInt
//     UUID
//     format.AnyComponent
type Metadata map[int]interface{}

func readMetadata(r io.Reader) (Metadata, error) {
	m := make(Metadata)
	for {
		id, err := ReadByte(r)
		if err != nil || id == 0xFF {
			return m, err
		}
		index := int(id)
		t, err := ReadByte(r)
		if err != nil {
			return m, err
		}
		switch t {
		case 0:
			var val int8
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		case 1:
			var val VarInt
			val, err = ReadVarInt(r)
			m[index] = int(val)
		case 2:
			var val float32
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		case 3:
			m[index], err = ReadString(r)
		case 4:
			str, err := ReadString(r)
			if err != nil {
				return m, err
			}
			var msg format.AnyComponent
			err = json.Unmarshal([]byte(str), &msg)
			m[index] = msg
		case 5:
			i := ItemStack{}
			err = i.Deserialize(r)
			m[index] = i
		case 6:
			m[index], err = ReadBool(r)
		case 7:
			var val [3]float32
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		case 8:
			var val Position
			err = binary.Read(r, binary.BigEndian, &val)
			m[index] = val
		case 9:
			var ok bool
			ok, err = ReadBool(r)
			if err != nil {
				return m, err
			}
			pos := new(Position)
			if ok {
				err = binary.Read(r, binary.BigEndian, pos)
			}
			m[index] = pos
		case 10: // Direction
			m[index], err = ReadVarInt(r)
		case 11:
			var ok bool
			ok, err = ReadBool(r)
			if err != nil {
				return m, err
			}
			uuid := new(UUID)
			if ok {
				err = uuid.Deserialize(r)
			}
			m[index] = uuid
		default:
			err = errors.New("invalid metadata type " + strconv.Itoa(int(t)))
		}
		if err != nil {
			return m, err
		}
	}
}

func writeMetadata(w io.Writer, m Metadata) error {
	for index, v := range m {
		err := WriteByte(w, byte(index))
		if err != nil {
			return err
		}
		switch v := v.(type) {
		case int8:
			WriteByte(w, 0)
			err = binary.Write(w, binary.BigEndian, v)
		case int:
			WriteByte(w, 1)
			err = WriteVarInt(w, VarInt(v))
		case float32:
			WriteByte(w, 2)
			err = binary.Write(w, binary.BigEndian, v)
		case string:
			WriteByte(w, 3)
			err = WriteString(w, v)
		case format.AnyComponent:
			WriteByte(w, 4)
			val, _ := json.Marshal(v)
			err = WriteString(w, string(val))
		case ItemStack:
			WriteByte(w, 5)
			v.Serialize(w)
		case bool:
			WriteByte(w, 6)
			err = WriteBool(w, v)
		case [3]float32:
			WriteByte(w, 7)
			err = binary.Write(w, binary.BigEndian, v)
		case Position:
			WriteByte(w, 8)
			err = binary.Write(w, binary.BigEndian, v)
		case *Position:
			WriteByte(w, 9)
			WriteBool(w, v != nil)
			if v != nil {
				err = binary.Write(w, binary.BigEndian, v)
			}
		case VarInt: // Direction
			WriteByte(w, 10)
			err = WriteVarInt(w, v)
		case *UUID:
			WriteByte(w, 11)
			WriteBool(w, v != nil)
			if v != nil {
				v.Serialize(w)
			}
		}
		if err != nil {
			return err
		}
	}
	return WriteByte(w, 0xFF)
}
