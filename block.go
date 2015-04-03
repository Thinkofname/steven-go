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

package main

import (
	"bytes"
	"fmt"
	"image"
	"reflect"
	"strconv"
	"strings"
)

var (
	nextBlockID   int
	blocks        [0x10000]Block
	blockSetsByID [0x100]*BlockSet
	missingBlock  = &baseBlock{
		name:        "missing_block",
		plugin:      "steven",
		cullAgainst: true,
	}
)

// Block is a type of tile in the world. All blocks, excluding the special
// 'missing block', belong to a set.
type Block interface {
	// Is returns whether this block is a member of the passed Set
	Is(s *BlockSet) bool

	Plugin() string
	Name() string

	ModelName() string
	ModelVariant() string
	Models() blockVariants
	ForceShade() bool
	ShouldCullAgainst() bool
	TintImage() *image.NRGBA
	IsTranslucent() bool

	String() string

	clone() Block
	toData() int
}

// base of most (if not all) blocks
type baseBlock struct {
	self          Block
	plugin, name  string
	Parent        *BlockSet
	cullAgainst   bool
	BlockVariants blockVariants
	translucent   bool
}

// Is returns whether this block is a member of the passed Set
func (b *baseBlock) Is(s *BlockSet) bool {
	return b.Parent == s
}

func (b *baseBlock) init(name string) {
	// plugin:name format
	if strings.ContainsRune(name, ':') {
		pos := strings.IndexRune(name, ':')
		b.plugin = name[:pos]
		b.name = name[pos+1:]
		return
	}
	b.name = name
	b.plugin = "minecraft"
	b.cullAgainst = true
}

func (b *baseBlock) String() string {
	return fmt.Sprintf("%s:%s", b.plugin, b.name)
}

func (b *baseBlock) Plugin() string {
	return b.plugin
}

func (b *baseBlock) Name() string {
	return b.name
}

func (b *baseBlock) Models() blockVariants {
	return b.BlockVariants
}

func (b *baseBlock) ModelName() string {
	return b.name
}
func (b *baseBlock) ModelVariant() string {
	return "normal"
}

func (b *baseBlock) toData() int {
	panic("toData on baseBlock")
}

func (b *baseBlock) ShouldCullAgainst() bool {
	return b.cullAgainst
}

func (b *baseBlock) ForceShade() bool {
	return false
}

func (b *baseBlock) TintImage() *image.NRGBA {
	return nil
}

func (b *baseBlock) IsTranslucent() bool {
	return b.translucent
}

func (b *baseBlock) clone() Block {
	return &baseBlock{
		plugin:      b.plugin,
		name:        b.name,
		Parent:      b.Parent,
		cullAgainst: b.cullAgainst,
		translucent: b.translucent,
	}
}

// GetBlockByCombinedID returns the block with the matching combined id.
// The combined id is:
//     block id << 4 | data
func GetBlockByCombinedID(id uint16) Block {
	b := blocks[id]
	if b == nil {
		return missingBlock
	}
	return b
}

// BlockSet is a collection of Blocks.
type BlockSet struct {
	ID int

	Base   Block
	Blocks []Block
	states []state
}

type state struct {
	name  string
	field reflect.StructField
	count int
}

func alloc(initial Block) *BlockSet {
	id := nextBlockID
	nextBlockID++
	bs := &BlockSet{
		ID:     id,
		Blocks: []Block{initial},
		Base:   initial,
	}
	blockSetsByID[id] = bs

	t := reflect.TypeOf(initial).Elem()

	v := reflect.ValueOf(initial).Elem()
	v.FieldByName("Parent").Set(
		reflect.ValueOf(bs),
	)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		s := f.Tag.Get("state")
		if s == "" {
			continue
		}
		args := strings.Split(s, ",")
		name := args[0]
		args = args[1:]

		var vals []interface{}
		switch f.Type.Kind() {
		case reflect.Bool:
			vals = []interface{}{false, true}
		case reflect.Int:
			var min, max int
			if args[0][0] != '@' {
				rnge := strings.Split(args[0], "-")
				min, _ = strconv.Atoi(rnge[0])
				max, _ = strconv.Atoi(rnge[1])
			} else {
				ret := v.Addr().MethodByName(args[0][1:]).Call([]reflect.Value{})
				min = int(ret[0].Int())
				max = int(ret[1].Int())
			}
			vals = make([]interface{}, max-min+1)
			for j := min; j <= max; j++ {
				vals[j-min] = j
			}
		case reflect.Uint:
			var min, max uint
			if args[0][0] != '@' {
				rnge := strings.Split(args[0], "-")
				mint, _ := strconv.Atoi(rnge[0])
				maxt, _ := strconv.Atoi(rnge[1])
				min = uint(mint)
				max = uint(maxt)
			} else {
				ret := v.Addr().MethodByName(args[0][1:]).Call([]reflect.Value{})
				min = uint(ret[0].Uint())
				max = uint(ret[1].Uint())
			}
			vals = make([]interface{}, max-min+1)
			for j := min; j <= max; j++ {
				vals[j-min] = j
			}
		default:
			panic("invalid state kind " + f.Type.Kind().String())
		}

		old := bs.Blocks
		bs.Blocks = make([]Block, 0, len(old)*len(vals))
		bs.states = append(bs.states, state{
			name:  name,
			field: f,
			count: len(vals),
		})
		for _, val := range vals {
			rval := reflect.ValueOf(val)
			for _, o := range old {
				// allocate a new block
				nb := o.clone()
				// set the new state
				ff := reflect.ValueOf(nb).Elem().Field(i)
				ff.Set(rval.Convert(ff.Type()))
				// now add back to the set
				bs.Blocks = append(bs.Blocks, nb)
			}
		}
	}
	return bs
}

func (bs *BlockSet) stringify(block Block) string {
	v := reflect.ValueOf(block).Elem()
	buf := bytes.NewBufferString(block.Plugin())
	buf.WriteRune(':')
	buf.WriteString(block.Name())
	if len(bs.states) > 0 {
		buf.WriteRune('[')
		for i, state := range bs.states {
			fv := v.FieldByIndex(state.field.Index)
			buf.WriteString(fmt.Sprintf("%s=%v", state.name, fv.Interface()))
			if i != len(bs.states)-1 {
				buf.WriteRune(',')
			}
		}
		buf.WriteRune(']')
	}
	return buf.String()
}

func initBlocks() {
	var missingModel *blockStateModel
	if missingModel = findStateModel("steven", "missing_block"); missingModel != nil {
		reflect.ValueOf(missingBlock).Elem().FieldByName("BlockVariants").Set(
			reflect.ValueOf(missingModel.variant("normal")),
		)
	}
	// Flatten the ids
	for _, bs := range blockSetsByID {
		if bs == nil {
			continue
		}
		for _, b := range bs.Blocks {
			data := b.toData()
			if data != -1 {
				blocks[(bs.ID<<4)|data] = b
			}
			// Liquids have custom rendering and air is never
			// rendered
			if _, ok := b.(*blockLiquid); ok || b.Is(BlockAir) {
				continue
			}
			if model := findStateModel(b.Plugin(), b.ModelName()); model != nil {
				if variants := model.variant(b.ModelVariant()); variants != nil {
					reflect.ValueOf(b).Elem().FieldByName("BlockVariants").Set(
						reflect.ValueOf(variants),
					)
					continue
				}
				fmt.Printf("Missing block variant (%s) for %s\n", b.ModelVariant(), b)
			} else {
				fmt.Printf("Missing block model for %s\n", b)
			}
			reflect.ValueOf(b).Elem().FieldByName("BlockVariants").Set(
				reflect.ValueOf(missingModel.variant("normal")),
			)

		}
	}
}
