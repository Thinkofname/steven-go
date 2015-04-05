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
	"fmt"
	"image"

	"github.com/thinkofdeath/steven/type/direction"
)

// Stone

type stoneVariant int

const (
	stoneNormal stoneVariant = iota
	stoneGranite
	stoneSmoothGranite
	stoneDiorite
	stoneSmoothDiorite
	stoneAndesite
	stoneSmoothAndesite
)

func (s stoneVariant) String() string {
	switch s {
	case stoneNormal:
		return "stone"
	case stoneGranite:
		return "granite"
	case stoneSmoothGranite:
		return "smooth_granite"
	case stoneDiorite:
		return "diorite"
	case stoneSmoothDiorite:
		return "smooth_diorite"
	case stoneAndesite:
		return "andesite"
	case stoneSmoothAndesite:
		return "smooth_andesite"
	}
	return fmt.Sprintf("stoneVariant(%d)", s)
}

type blockstone struct {
	baseBlock
	Variant stoneVariant `state:"variant,0-6"`
}

func initStone(name string) *BlockSet {
	l := &blockstone{}
	l.init(name)
	set := alloc(l)
	return set
}

func (b *blockstone) String() string {
	return b.Parent.stringify(b)
}

func (b *blockstone) clone() Block {
	return &blockstone{
		baseBlock: *(b.baseBlock.clone().(*baseBlock)),
		Variant:   b.Variant,
	}
}

func (b *blockstone) ModelName() string {
	return b.Variant.String()
}

func (b *blockstone) toData() int {
	data := int(b.Variant)
	return data
}

// Grass

type blockGrass struct {
	baseBlock
	Snowy bool `state:"snowy"`
}

func initGrass() *BlockSet {
	g := &blockGrass{}
	g.init("grass")
	set := alloc(g)
	return set
}

func (g *blockGrass) String() string {
	return g.Parent.stringify(g)
}

func (g *blockGrass) clone() Block {
	return &blockGrass{
		baseBlock: *(g.baseBlock.clone().(*baseBlock)),
		Snowy:     g.Snowy,
	}
}

func (g *blockGrass) ModelVariant() string {
	return fmt.Sprintf("snowy=%t", g.Snowy)
}

func (g *blockGrass) TintImage() *image.NRGBA {
	return grassBiomeColors
}

func (g *blockGrass) toData() int {
	if g.Snowy {
		return -1
	}
	return 0
}

// Tall grass

type tallGrassType int

const (
	tallGrassDeadBush = iota
	tallGrass
	tallGrassFern
)

func (t tallGrassType) String() string {
	switch t {
	case tallGrassDeadBush:
		return "dead_bush"
	case tallGrass:
		return "tall_grass"
	case tallGrassFern:
		return "fern"
	}
	return fmt.Sprintf("tallGrassType(%d)", t)
}

type blockTallGrass struct {
	baseBlock
	Type tallGrassType `state:"type,0-2"`
}

func initTallGrass() *BlockSet {
	t := &blockTallGrass{}
	t.init("tallgrass")
	t.cullAgainst = false
	set := alloc(t)
	return set
}

func (t *blockTallGrass) String() string {
	return t.Parent.stringify(t)
}

func (t *blockTallGrass) clone() Block {
	return &blockTallGrass{
		baseBlock: *(t.baseBlock.clone().(*baseBlock)),
		Type:      t.Type,
	}
}

func (t *blockTallGrass) ModelName() string {
	return t.Type.String()
}

func (t *blockTallGrass) TintImage() *image.NRGBA {
	return grassBiomeColors
}

func (t *blockTallGrass) toData() int {
	return int(t.Type)
}

// Bed

type bedPart int

const (
	bedHead bedPart = iota
	bedFoot
)

func (b bedPart) String() string {
	switch b {
	case bedHead:
		return "head"
	case bedFoot:
		return "foot"
	}
	return fmt.Sprintf("bedPart(%d)", b)
}

type blockBed struct {
	baseBlock
	Facing   direction.Type `state:"facing,2-5"`
	Occupied bool           `state:"occupied"`
	Part     bedPart        `state:"part,0-1"`
}

func initBed(name string) *BlockSet {
	b := &blockBed{}
	b.init(name)
	b.cullAgainst = false
	set := alloc(b)
	return set
}

func (b *blockBed) String() string {
	return b.Parent.stringify(b)
}

func (b *blockBed) clone() Block {
	return &blockBed{
		baseBlock: *(b.baseBlock.clone().(*baseBlock)),
		Facing:    b.Facing,
		Occupied:  b.Occupied,
		Part:      b.Part,
	}
}

func (b *blockBed) ModelVariant() string {
	return fmt.Sprintf("facing=%s,part=%s", b.Facing, b.Part)
}

func (b *blockBed) toData() int {
	data := 0
	switch b.Facing {
	case direction.South:
		data = 0
	case direction.West:
		data = 1
	case direction.North:
		data = 2
	case direction.East:
		data = 3
	}
	if b.Occupied {
		data |= 0x4
	}
	if b.Part == bedHead {
		data |= 0x8
	}
	return data
}

// Sponge

type blockSponge struct {
	baseBlock
	Wet bool `state:"wet"`
}

func initSponge(name string) *BlockSet {
	b := &blockSponge{}
	b.init(name)
	set := alloc(b)
	return set
}

func (b *blockSponge) String() string {
	return b.Parent.stringify(b)
}

func (b *blockSponge) clone() Block {
	return &blockSponge{
		baseBlock: *(b.baseBlock.clone().(*baseBlock)),
		Wet:       b.Wet,
	}
}

func (b *blockSponge) ModelVariant() string {
	return fmt.Sprintf("wet=%t", b.Wet)
}

func (b *blockSponge) toData() int {
	data := 0
	if b.Wet {
		data = 1
	}
	return data
}

// Door

type doorHalf int

const (
	doorUpper doorHalf = iota
	doorLower
)

func (d doorHalf) String() string {
	switch d {
	case doorUpper:
		return "upper"
	case doorLower:
		return "lower"
	}
	return fmt.Sprintf("doorLower(%d)", d)
}

type doorHinge int

const (
	doorLeft doorHinge = iota
	doorRight
)

func (d doorHinge) String() string {
	switch d {
	case doorLeft:
		return "left"
	case doorRight:
		return "right"
	}
	return fmt.Sprintf("doorRight(%d)", d)
}

type blockDoor struct {
	baseBlock
	Facing  direction.Type `state:"facing,2-5"`
	Half    doorHalf       `state:"half,0-1"`
	Hinge   doorHinge      `state:"hinge,0-1"`
	Open    bool           `state:"open"`
	Powered bool           `state:"powered"`
}

func initDoor(name string) *BlockSet {
	b := &blockDoor{}
	b.init(name)
	b.cullAgainst = false
	set := alloc(b)
	return set
}

func (b *blockDoor) String() string {
	return b.Parent.stringify(b)
}

func (b *blockDoor) clone() Block {
	return &blockDoor{
		baseBlock: *(b.baseBlock.clone().(*baseBlock)),
		Facing:    b.Facing,
		Half:      b.Half,
		Hinge:     b.Hinge,
		Open:      b.Open,
		Powered:   b.Powered,
	}
}

func (b *blockDoor) ModelVariant() string {
	return fmt.Sprintf("facing=%s,half=%s,hinge=%s,open=%t", b.Facing, b.Half, b.Hinge, b.Open)
}

func (b *blockDoor) UpdateState(x, y, z int) Block {
	if b.Half == doorUpper {
		o := chunkMap.Block(x, y-1, z)
		if d, ok := o.(*blockDoor); ok {
			return b.
				Set("facing", d.Facing).
				Set("open", d.Open)
		}
		return b
	}
	o := chunkMap.Block(x, y+1, z)
	if d, ok := o.(*blockDoor); ok {
		return b.Set("hinge", d.Hinge)
	}
	return b
}

func (b *blockDoor) toData() int {
	data := 0
	if b.Half == doorUpper {
		data |= 0x8
		if b.Hinge == doorRight {
			data |= 0x1
		}
		if b.Powered {
			data |= 0x2
		}
	} else {
		switch b.Facing {
		case direction.East:
			data = 0
		case direction.South:
			data = 1
		case direction.West:
			data = 2
		case direction.North:
			data = 3
		}
		if b.Open {
			data |= 0x4
		}
	}
	return data
}

// Dispenser

type blockDispenser struct {
	baseBlock
	Facing    direction.Type `state:"facing,0-5"`
	Triggered bool           `state:"triggered"`
}

func initDispenser(name string) *BlockSet {
	b := &blockDispenser{}
	b.init(name)
	set := alloc(b)
	return set
}

func (b *blockDispenser) String() string {
	return b.Parent.stringify(b)
}

func (b *blockDispenser) clone() Block {
	return &blockDispenser{
		baseBlock: *(b.baseBlock.clone().(*baseBlock)),
		Facing:    b.Facing,
		Triggered: b.Triggered,
	}
}

func (b *blockDispenser) ModelVariant() string {
	return fmt.Sprintf("facing=%s", b.Facing)
}

func (b *blockDispenser) toData() int {
	data := 0
	switch b.Facing {
	case direction.Down:
		data = 0
	case direction.Up:
		data = 1
	case direction.North:
		data = 2
	case direction.South:
		data = 3
	case direction.West:
		data = 4
	case direction.East:
		data = 5
	}
	if b.Triggered {
		data |= 0x8
	}
	return data
}

// Powered rail

type railShape int

const (
	rsNorthSouth railShape = iota
	rsEastWest
	rsAscendingEast
	rsAscendingWest
	rsAscendingNorth
	rsAscendingSouth
	rsSouthEast
	rsSouthWest
	rsNorthWest
	rsNorthEast
)

func (r railShape) String() string {
	switch r {
	case rsNorthSouth:
		return "north_south"
	case rsEastWest:
		return "east_west"
	case rsAscendingNorth:
		return "ascending_north"
	case rsAscendingSouth:
		return "ascending_south"
	case rsAscendingEast:
		return "ascending_east"
	case rsAscendingWest:
		return "ascending_west"
	case rsSouthEast:
		return "south_east"
	case rsSouthWest:
		return "south_west"
	case rsNorthWest:
		return "north_west"
	case rsNorthEast:
		return "north_east"
	}
	return fmt.Sprintf("railShape(%d)", r)
}

type blockPoweredRail struct {
	baseBlock
	Shape   railShape `state:"shape,0-5"`
	Powered bool      `state:"powered"`
}

func initPoweredRail(name string) *BlockSet {
	b := &blockPoweredRail{}
	b.init(name)
	b.cullAgainst = false
	set := alloc(b)
	return set
}

func (b *blockPoweredRail) String() string {
	return b.Parent.stringify(b)
}

func (b *blockPoweredRail) clone() Block {
	return &blockPoweredRail{
		baseBlock: *(b.baseBlock.clone().(*baseBlock)),
		Shape:     b.Shape,
		Powered:   b.Powered,
	}
}

func (b *blockPoweredRail) ModelVariant() string {
	return fmt.Sprintf("powered=%t,shape=%s", b.Powered, b.Shape)
}

func (b *blockPoweredRail) toData() int {
	data := int(b.Shape)
	if b.Powered {
		data |= 0x8
	}
	return data
}

// Rail

type blockRail struct {
	baseBlock
	Shape railShape `state:"shape,0-9"`
}

func initRail(name string) *BlockSet {
	b := &blockRail{}
	b.init(name)
	b.cullAgainst = false
	set := alloc(b)
	return set
}

func (b *blockRail) String() string {
	return b.Parent.stringify(b)
}

func (b *blockRail) clone() Block {
	return &blockRail{
		baseBlock: *(b.baseBlock.clone().(*baseBlock)),
		Shape:     b.Shape,
	}
}

func (b *blockRail) ModelVariant() string {
	return fmt.Sprintf("shape=%s", b.Shape)
}

func (b *blockRail) toData() int {
	return int(b.Shape)
}

// Dead bush

type blockDeadBush struct {
	baseBlock
}

func initDeadBush(name string) *BlockSet {
	b := &blockDeadBush{}
	b.init(name)
	b.cullAgainst = false
	set := alloc(b)
	return set
}

func (b *blockDeadBush) String() string {
	return b.Parent.stringify(b)
}

func (b *blockDeadBush) clone() Block {
	return &blockDeadBush{
		baseBlock: *(b.baseBlock.clone().(*baseBlock)),
	}
}

func (b *blockDeadBush) ModelName() string {
	return "dead_bush"
}

func (b *blockDeadBush) toData() int {
	return 0
}

// Fence

type blockFence struct {
	baseBlock
	Wood  bool
	North bool `state:"north"`
	South bool `state:"south"`
	East  bool `state:"east"`
	West  bool `state:"west"`
}

func initFence(name string, wood bool) *BlockSet {
	b := &blockFence{}
	b.init(name)
	b.Wood = wood
	b.cullAgainst = false
	set := alloc(b)
	return set
}

func (b *blockFence) String() string {
	return b.Parent.stringify(b)
}

func (b *blockFence) clone() Block {
	return &blockFence{
		baseBlock: *(b.baseBlock.clone().(*baseBlock)),
		North:     b.North,
		South:     b.South,
		East:      b.East,
		West:      b.West,
		Wood:      b.Wood,
	}
}

func (b *blockFence) UpdateState(x, y, z int) Block {
	var block Block = b
	for _, d := range direction.Values {
		if d < 2 {
			continue
		}
		ox, oy, oz := d.Offset()
		bl := chunkMap.Block(x+ox, y+oy, z+oz)
		if fence, ok := bl.(*blockFence); bl.ShouldCullAgainst() || (ok && fence.Wood == b.Wood) {
			block = block.Set(d.String(), true)
		} else {
			block = block.Set(d.String(), false)
		}
	}
	return block
}

func (b *blockFence) ModelVariant() string {
	return fmt.Sprintf("east=%t,north=%t,south=%t,west=%t", b.East, b.North, b.South, b.West)
}

func (b *blockFence) toData() int {
	if !b.North && !b.South && !b.East && !b.West {
		return 0
	}
	return -1
}
