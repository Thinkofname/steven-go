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
	"io"

	"github.com/thinkofdeath/phteven/encoding/nbt"
)

// ItemStack is a stack of items of a single type that can be serilized
// in the protocol.
type ItemStack struct {
	ID     int16
	Count  byte
	Damage int16
	NBT    *nbt.Compound
}

// Serialize writes the item stack into the writer.
func (i *ItemStack) Serialize(w io.Writer) error {
	if err := binary.Write(w, binary.BigEndian, i.ID); err != nil {
		return err
	}
	if i.ID == -1 {
		return nil
	}
	if err := binary.Write(w, binary.BigEndian, i.Count); err != nil {
		return err
	}
	if err := binary.Write(w, binary.BigEndian, i.Damage); err != nil {
		return err
	}
	return WriteNBT(w, i.NBT)
}

// Deserialize reads an item stack from the reader into this item stack.
func (i *ItemStack) Deserialize(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &i.ID); err != nil {
		return err
	}
	if i.ID == -1 {
		return nil
	}
	if err := binary.Read(r, binary.BigEndian, &i.Count); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &i.Damage); err != nil {
		return err
	}
	var err error
	i.NBT, err = ReadNBT(r)
	return err
}
