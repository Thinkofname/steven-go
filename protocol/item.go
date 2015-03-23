package protocol

import (
	"encoding/binary"
	"github.com/thinkofdeath/steven/encoding/nbt"
	"io"
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
	return writeNBT(w, i.NBT)
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
	i.NBT, err = readNBT(r)
	return err
}
