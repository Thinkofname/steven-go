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

package steven

import (
	"github.com/thinkofdeath/steven/encoding/nbt"
	"github.com/thinkofdeath/steven/protocol"
)

type ItemStack struct {
	Type  ItemType
	Count int
}

func ItemStackFromProtocol(p protocol.ItemStack) *ItemStack {
	it := ItemById(int(p.ID))
	if it == nil {
		return nil
	}
	i := &ItemStack{
		Type:  it,
		Count: int(p.Count),
	}
	i.Type.ParseDamage(p.Damage)
	if p.NBT != nil {
		i.Type.ParseTag(p.NBT)
	}
	return i
}

type ItemType interface {
	Name() string
	NameLocaleKey() string

	Stackable() bool

	ParseDamage(d int16)
	ParseTag(tag *nbt.Compound)
}

func ItemById(id int) (ty ItemType) {
	if id == -1 {
		return nil
	}
	if id < 256 {
		ty = ItemOfBlock(blockSetsByID[id].Base)
	} else {
		if f, ok := itemsByID[id]; ok {
			ty = f()
		}
	}
	if ty == nil {
		ty = ItemOfBlock(Blocks.Stone.Base)
	}
	return ty
}

type displayTag struct {
	name string
	lore []string
}

func (d *displayTag) ParseTag(tag *nbt.Compound) {
	display, ok := tag.Items["display"].(*nbt.Compound)
	if !ok {
		return
	}
	d.name, _ = display.Items["Name"].(string)
	lore, ok := display.Items["Lore"].([]interface{})
	if !ok {
		return
	}
	d.lore = make([]string, len(lore))
	for i := range lore {
		d.lore[i], _ = lore[i].(string)
	}
}
func (d *displayTag) DisplayName() string { return d.name }
func (d *displayTag) Lore() []string      { return d.lore }

type DisplayTag interface {
	DisplayName() string
	Lore() []string
}

type blockItem struct {
	itemNamed
	block Block
	displayTag
}

func ItemOfBlock(b Block) ItemType {
	return &blockItem{
		block: b,
		itemNamed: itemNamed{
			name: b.ModelName(),
		},
	}
}

func (b *blockItem) NameLocaleKey() string {
	return b.block.NameLocaleKey()
}

func (b *blockItem) ParseDamage(d int16) {
	d &= 0xF
	nb := GetBlockByCombinedID(uint16(b.block.BlockSet().ID<<4) | uint16(d))
	if nb.Is(b.block.BlockSet()) {
		b.block = nb
		b.itemNamed.name = nb.ModelName()
	}
}
func (b *blockItem) ParseTag(tag *nbt.Compound) {
	b.displayTag.ParseTag(tag)
}

func (b *blockItem) Stackable() bool {
	return true
}

type itemSimpleLocale struct {
	locale string
}

func (i *itemSimpleLocale) NameLocaleKey() string {
	return i.locale
}

type itemDamagable struct {
	damage, maxDamage int16
}

func (i *itemDamagable) ParseDamage(d int16) {
	i.damage = d
}
func (i *itemDamagable) Damage() int16    { return i.damage }
func (i *itemDamagable) MaxDamage() int16 { return i.maxDamage }

type ItemDamagable interface {
	Damage() int16
	MaxDamage() int16
}

type itemNamed struct {
	name string
}

func (i *itemNamed) Name() string {
	return i.name
}
