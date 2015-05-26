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
	"bytes"
	"fmt"
	"image"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/thinkofdeath/steven/type/vmath"
)

var (
	nextBlockID   int
	blocks        [0x10000]Block
	blockSetsByID [0x100]*BlockSet
	allBlocks     = make([]Block, 0, math.MaxUint16)
)

// Block is a type of tile in the world. All blocks, excluding the special
// 'missing block', belong to a set.
type Block interface {
	// Is returns whether this block is a member of the passed Set
	Is(s *BlockSet) bool
	BlockSet() *BlockSet

	Plugin() string
	Name() string
	SID() uint16
	Set(key string, val interface{}) Block
	UpdateState(x, y, z int) Block
	states() []blockState

	Collidable() bool
	CollisionBounds() []vmath.AABB

	Hardness() float64

	Renderable() bool
	ModelName() string
	ModelVariant() string
	Models() blockVariants
	ForceShade() bool
	ShouldCullAgainst() bool
	TintImage() *image.NRGBA
	IsTranslucent() bool

	LightReduction() int
	LightEmitted() int
	String() string
	NameLocaleKey() string

	CreateBlockEntity() BlockEntity

	init(name string)
	toData() int
}

type blockState struct {
	Key   string
	Value interface{}
}

// base of most (if not all) blocks
type baseBlock struct {
	plugin, name  string
	Parent        *BlockSet
	Index         int
	StevenID      uint16
	cullAgainst   bool
	BlockVariants blockVariants
	translucent   bool
	collidable    bool
	renderable    bool
	bounds        []vmath.AABB
	hardness      float64
}

// Is returns whether this block is a member of the passed Set
func (b *baseBlock) Is(s *BlockSet) bool {
	return b.Parent == s
}

func (b *baseBlock) BlockSet() *BlockSet {
	return b.Parent
}

func (b *baseBlock) init(name string) {
	// plugin:name format
	if strings.ContainsRune(name, ':') {
		pos := strings.IndexRune(name, ':')
		b.plugin = name[:pos]
		b.name = name[pos+1:]
	} else {
		b.name = name
		b.plugin = "minecraft"
	}
	b.cullAgainst = true
	b.collidable = true
	b.renderable = true
	b.hardness = 1.0
}

func (b *baseBlock) NameLocaleKey() string {
	return fmt.Sprintf("tile.%s.name", b.name)
}

func (b *baseBlock) Hardness() float64 {
	return b.hardness
}

func (b *baseBlock) String() string {
	return b.Parent.stringify(b.Parent.Blocks[b.Index])
}

func (b *baseBlock) Plugin() string {
	return b.plugin
}

func (b *baseBlock) Name() string {
	return b.name
}

func (b *baseBlock) SID() uint16 {
	return b.StevenID
}

func (b *baseBlock) Collidable() bool {
	return b.collidable
}

func (b *baseBlock) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(0, 0, 0, 1.0, 1.0, 1.0),
		}
	}
	return b.bounds
}

func (b *baseBlock) Renderable() bool {
	return b.renderable
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

func (b *baseBlock) LightReduction() int {
	if b.ShouldCullAgainst() {
		return 15
	}
	return 0
}

func (b *baseBlock) LightEmitted() int {
	return 0
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

func (b *baseBlock) CreateBlockEntity() BlockEntity {
	return nil
}

func (b *baseBlock) UpdateState(x, y, z int) Block {
	return b.Parent.Blocks[b.Index]
}

func (b *baseBlock) Set(key string, val interface{}) Block {
	index := 0
	cur := reflect.ValueOf(b.Parent.Blocks[b.Index]).Elem()
	for i := range b.Parent.states {
		state := b.Parent.states[len(b.Parent.states)-1-i]
		index *= state.count
		var sval reflect.Value
		// Need to lookup the current value if this isn't the
		// state we are changing
		if state.name != key {
			sval = reflect.ValueOf(cur.FieldByIndex(state.field.Index).Interface())
		} else {
			sval = reflect.ValueOf(val)
		}
		args := strings.Split(state.field.Tag.Get("state"), ",")
		args = args[1:]
		switch state.field.Type.Kind() {
		case reflect.Bool:
			if sval.Bool() {
				index += 1
			}
		case reflect.Int:
			var min int
			if args[0][0] != '@' {
				rnge := strings.Split(args[0], "-")
				min, _ = strconv.Atoi(rnge[0])
			} else {
				ret := cur.Addr().MethodByName(args[0][1:]).Call([]reflect.Value{})
				min = int(ret[0].Int())
			}
			v := int(sval.Int())
			index += v - min
		case reflect.Uint:
			var min uint
			if args[0][0] != '@' {
				rnge := strings.Split(args[0], "-")
				mint, _ := strconv.Atoi(rnge[0])
				min = uint(mint)
			} else {
				ret := cur.Addr().MethodByName(args[0][1:]).Call([]reflect.Value{})
				min = uint(ret[0].Uint())
			}
			v := uint(sval.Uint())
			index += int(v - min)
		default:
			panic("invalid state kind " + state.field.Type.Kind().String())
		}

	}
	return b.Parent.Blocks[index]
}

func (b *baseBlock) states() (out []blockState) {
	self := reflect.ValueOf(b.Parent.Blocks[b.Index]).Elem()
	for _, state := range b.Parent.states {
		out = append(out, blockState{
			Key:   state.name,
			Value: self.FieldByIndex(state.field.Index).Interface(),
		})
	}
	return
}

// GetBlockByCombinedID returns the block with the matching combined id.
// The combined id is:
//     block id << 4 | data
func GetBlockByCombinedID(id uint16) Block {
	b := blocks[id]
	if b == nil {
		return Blocks.MissingBlock.Base
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
				nb := cloneBlock(o)
				// set the new state
				ff := reflect.ValueOf(nb).Elem().Field(i)
				ff.Set(rval.Convert(ff.Type()))
				// now add back to the set
				bs.Blocks = append(bs.Blocks, nb)
			}
		}
	}
	bs.Base = bs.Blocks[0]
	return bs
}

func cloneBlock(b Block) Block {
	v := reflect.ValueOf(b).Elem()
	nv := reflect.New(v.Type()).Elem()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		nv.Field(i).Set(f)
	}
	return nv.Addr().Interface().(Block)
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
	missingModel := findStateModel("steven", "missing_block")
	// Flatten the ids
	for _, bs := range blockSetsByID {
		if bs == nil {
			continue
		}
		for i, b := range bs.Blocks {
			br := reflect.ValueOf(b).Elem()
			br.FieldByName("Index").SetInt(int64(i))
			br.FieldByName("StevenID").SetUint(uint64(len(allBlocks)))
			allBlocks = append(allBlocks, b)
			if len(allBlocks) > math.MaxUint16 {
				panic("ran out of ids, time to do this correctly :(")
			}
			data := b.toData()
			if data != -1 {
				blocks[(bs.ID<<4)|data] = b
			}
			// Liquids have custom rendering
			if _, ok := b.(*blockLiquid); ok || !b.Renderable() {
				continue
			}

			if model := findStateModel(b.Plugin(), b.ModelName()); model != nil {
				if variants := model.variant(b.ModelVariant()); variants != nil {
					br.FieldByName("BlockVariants").Set(
						reflect.ValueOf(variants),
					)
					continue
				}
				log.Printf("Missing block variant (%s) for %s\n", b.ModelVariant(), b)
			} else {
				log.Printf("Missing block model for %s\n", b)
			}
			br.FieldByName("BlockVariants").Set(
				reflect.ValueOf(missingModel.variant("normal")),
			)

		}
	}
}

func reinitBlocks() {
	blockStateModels = map[pluginKey]*blockStateModel{}
	missingModel := findStateModel("steven", "missing_block")
	// Flatten the ids
	for _, bs := range blockSetsByID {
		if bs == nil {
			continue
		}
		for _, b := range bs.Blocks {
			br := reflect.ValueOf(b).Elem()
			// Liquids have custom rendering
			if _, ok := b.(*blockLiquid); ok || !b.Renderable() {
				continue
			}

			if model := findStateModel(b.Plugin(), b.ModelName()); model != nil {
				if variants := model.variant(b.ModelVariant()); variants != nil {
					br.FieldByName("BlockVariants").Set(
						reflect.ValueOf(variants),
					)
					continue
				}
				log.Printf("Missing block variant (%s) for %s\n", b.ModelVariant(), b)
			} else {
				log.Printf("Missing block model for %s\n", b)
			}
			br.FieldByName("BlockVariants").Set(
				reflect.ValueOf(missingModel.variant("normal")),
			)

		}
	}
}
