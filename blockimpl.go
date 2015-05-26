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
	"fmt"
	"image"
	"math"
	"reflect"

	"github.com/thinkofdeath/steven/type/direction"
	"github.com/thinkofdeath/steven/type/vmath"
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

type blockStone struct {
	baseBlock
	Variant stoneVariant `state:"variant,0-6"`
}

func (b *blockStone) ModelName() string {
	return b.Variant.String()
}

func (b *blockStone) NameLocaleKey() string {
	switch b.Variant {
	case stoneNormal:
		return "tile.stone.stone.name"
	case stoneGranite:
		return "tile.stone.granite.name"
	case stoneSmoothGranite:
		return "tile.stone.graniteSmooth.name"
	case stoneDiorite:
		return "tile.stone.diorite.name"
	case stoneSmoothDiorite:
		return "tile.stone.dioriteSmooth.name"
	case stoneAndesite:
		return "tile.stone.andesite.name"
	case stoneSmoothAndesite:
		return "tile.stone.andesiteSmooth.name"

	}
	return "unknown"
}

func (b *blockStone) toData() int {
	data := int(b.Variant)
	return data
}

// Grass

type blockGrass struct {
	baseBlock
	Snowy bool `state:"snowy"`
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

func (b *blockTallGrass) load(tag reflect.StructTag) {
	b.cullAgainst = false
	b.collidable = false
}

func (b *blockTallGrass) ModelName() string {
	return b.Type.String()
}

func (b *blockTallGrass) TintImage() *image.NRGBA {
	return grassBiomeColors
}

func (b *blockTallGrass) toData() int {
	return int(b.Type)
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

func (b *blockBed) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockBed) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(0, 0, 0, 1, 9.0/16.0, 1),
		}
	}
	return b.bounds
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

func (b *blockDoor) load(tag reflect.StructTag) {
	b.cullAgainst = false
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

func (b *blockPoweredRail) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockPoweredRail) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(0, 0, 0, 1.0, 1.0/16.0, 1.0),
		}
	}
	return b.bounds
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

func (b *blockRail) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockRail) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(0, 0, 0, 1.0, 1.0/16.0, 1.0),
		}
	}
	return b.bounds
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

func (b *blockDeadBush) load(tag reflect.StructTag) {
	b.cullAgainst = false
	b.collidable = false
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

func (b *blockFence) load(tag reflect.StructTag) {
	getBool := wrapTagBool(tag)
	b.cullAgainst = false
	b.Wood = getBool("wood", true)
}

func (b *blockFence) UpdateState(x, y, z int) Block {
	var block Block = b
	for _, d := range direction.Values {
		if d < 2 {
			continue
		}
		ox, oy, oz := d.Offset()
		bl := chunkMap.Block(x+ox, y+oy, z+oz)
		_, ok2 := bl.(*blockFenceGate)
		if fence, ok := bl.(*blockFence); bl.ShouldCullAgainst() || (ok && fence.Wood == b.Wood) || ok2 {
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

// Fence Gate

type blockFenceGate struct {
	baseBlock
	Facing  direction.Type `state:"facing,2-5"`
	InWall  bool           `state:"in_wall"`
	Open    bool           `state:"open"`
	Powered bool           `state:"powered"`
}

func (b *blockFenceGate) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockFenceGate) UpdateState(x, y, z int) Block {
	var block Block = b
	ox, oy, oz := b.Facing.Clockwise().Offset()
	if _, ok := chunkMap.Block(x+ox, y+oy, z+oz).(*blockWall); ok {
		return block.Set("in_wall", true)
	}
	ox, oy, oz = b.Facing.CounterClockwise().Offset()
	if _, ok := chunkMap.Block(x+ox, y+oy, z+oz).(*blockWall); ok {
		return block.Set("in_wall", true)
	}
	return block.Set("in_wall", false)
}

func (b *blockFenceGate) ModelVariant() string {
	return fmt.Sprintf("facing=%s,in_wall=%t,open=%t", b.Facing, b.InWall, b.Open)
}

func (b *blockFenceGate) toData() int {
	if b.Powered || b.InWall {
		return -1
	}
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
	if b.Open {
		data |= 0x4
	}
	return data
}

// Wall

type wallVariant int

const (
	wvCobblestone wallVariant = iota
	wvMossyCobblestone
)

func (w wallVariant) String() string {
	switch w {
	case wvCobblestone:
		return "cobblestone"
	case wvMossyCobblestone:
		return "mossy_cobblestone"
	}
	return fmt.Sprintf("wallVariant(%d)", w)
}

type blockWall struct {
	baseBlock
	Variant wallVariant `state:"variant,0-1"`
	Up      bool        `state:"up"`
	North   bool        `state:"north"`
	South   bool        `state:"south"`
	East    bool        `state:"east"`
	West    bool        `state:"west"`
}

func (b *blockWall) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockWall) UpdateState(x, y, z int) Block {
	var block Block = b
	for _, d := range direction.Values {
		if d == direction.Down {
			continue
		}
		ox, oy, oz := d.Offset()
		bl := chunkMap.Block(x+ox, y+oy, z+oz)
		_, ok := bl.(*blockWall)
		_, ok2 := bl.(*blockFenceGate)
		if bl.ShouldCullAgainst() || ok || ok2 {
			block = block.Set(d.String(), true)
		} else {
			block = block.Set(d.String(), false)
		}
	}
	return block
}

func (b *blockWall) ModelName() string {
	return b.Variant.String() + "_wall"
}

func (b *blockWall) ModelVariant() string {
	return fmt.Sprintf("east=%t,north=%t,south=%t,up=%t,west=%t", b.East, b.North, b.South, b.Up, b.West)
}

func (b *blockWall) toData() int {
	if !b.North && !b.South && !b.East && !b.West && !b.Up {
		return int(b.Variant)
	}
	return -1
}

// Stained glass

type color int

const (
	cWhite color = iota
	cOrange
	cMagenta
	cLightBlue
	cYellow
	cLime
	cPink
	cGray
	cSilver
	cCyan
	cPurple
	cBlue
	cBrown
	cGreen
	cRed
	cBlack
)

func (c color) String() string {
	switch c {
	case cWhite:
		return "white"
	case cOrange:
		return "orange"
	case cMagenta:
		return "magenta"
	case cLightBlue:
		return "light_blue"
	case cYellow:
		return "yellow"
	case cLime:
		return "lime"
	case cPink:
		return "pink"
	case cGray:
		return "gray"
	case cSilver:
		return "silver"
	case cCyan:
		return "cyan"
	case cPurple:
		return "purple"
	case cBlue:
		return "blue"
	case cBrown:
		return "brown"
	case cGreen:
		return "green"
	case cRed:
		return "red"
	case cBlack:
		return "black"
	}
	return fmt.Sprintf("color(%d)", c)
}

type blockStainedGlass struct {
	baseBlock
	Color color `state:"color,0-15"`
}

func (b *blockStainedGlass) load(tag reflect.StructTag) {
	b.translucent = true
	b.cullAgainst = false
}

func (b *blockStainedGlass) ModelName() string {
	return b.Color.String() + "_stained_glass"
}

func (b *blockStainedGlass) toData() int {
	return int(b.Color)
}

// Connectable

type blockConnectable struct {
	baseBlock
	North bool `state:"north"`
	South bool `state:"south"`
	East  bool `state:"east"`
	West  bool `state:"west"`
}

func (b *blockConnectable) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (blockConnectable) connectable() {}

func (b *blockConnectable) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		all := !b.North && !b.South && !b.West && !b.East
		aa := vmath.NewAABB(0, 0, 7.0/16.0, 1.0, 1.0, 9.0/16.0)
		bb := vmath.NewAABB(7.0/16.0, 0, 0, 9.0/16.0, 1.0, 1.0)
		if !b.North && !all {
			bb.Min[2] = 7.0 / 16.0
		}
		if !b.South && !all {
			bb.Max[2] = 9.0 / 16.0
		}
		if !b.West && !all {
			aa.Min[0] = 7.0 / 16.0
		}
		if !b.East && !all {
			aa.Max[0] = 9.0 / 16.0
		}
		b.bounds = []vmath.AABB{aa, bb}
	}
	return b.bounds
}

func (b *blockConnectable) UpdateState(x, y, z int) Block {
	type connectable interface {
		connectable()
	}
	var block Block = b
	for _, d := range direction.Values {
		if d < 2 {
			continue
		}
		ox, oy, oz := d.Offset()
		bl := chunkMap.Block(x+ox, y+oy, z+oz)
		if _, ok := bl.(connectable); bl.ShouldCullAgainst() || ok {
			block = block.Set(d.String(), true)
		} else {
			block = block.Set(d.String(), false)
		}
	}
	return block
}

func (b *blockConnectable) ModelVariant() string {
	return fmt.Sprintf("east=%t,north=%t,south=%t,west=%t", b.East, b.North, b.South, b.West)
}

func (b *blockConnectable) toData() int {
	if !b.North && !b.South && !b.East && !b.West {
		return 0
	}
	return -1
}

// Stained Glass Pane

type blockStainedGlassPane struct {
	baseBlock
	Color color `state:"color,0-15"`
	North bool  `state:"north"`
	South bool  `state:"south"`
	East  bool  `state:"east"`
	West  bool  `state:"west"`
}

func (b *blockStainedGlassPane) load(tag reflect.StructTag) {
	b.translucent = true
	b.cullAgainst = false
}

func (b *blockStainedGlassPane) ModelName() string {
	return b.Color.String() + "_" + b.name
}

func (b *blockStainedGlassPane) ModelVariant() string {
	return fmt.Sprintf("east=%t,north=%t,south=%t,west=%t", b.East, b.North, b.South, b.West)
}

func (blockStainedGlassPane) connectable() {}

func (b *blockStainedGlassPane) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		all := !b.North && !b.South && !b.West && !b.East
		aa := vmath.NewAABB(0, 0, 7.0/16.0, 1.0, 1.0, 9.0/16.0)
		bb := vmath.NewAABB(7.0/16.0, 0, 0, 9.0/16.0, 1.0, 1.0)
		if !b.North && !all {
			bb.Min[2] = 7.0 / 16.0
		}
		if !b.South && !all {
			bb.Max[2] = 9.0 / 16.0
		}
		if !b.West && !all {
			aa.Min[0] = 7.0 / 16.0
		}
		if !b.East && !all {
			aa.Max[0] = 9.0 / 16.0
		}
		b.bounds = []vmath.AABB{aa, bb}
	}
	return b.bounds
}

func (b *blockStainedGlassPane) UpdateState(x, y, z int) Block {
	type connectable interface {
		connectable()
	}
	var block Block = b
	for _, d := range direction.Values {
		if d < 2 {
			continue
		}
		ox, oy, oz := d.Offset()
		bl := chunkMap.Block(x+ox, y+oy, z+oz)
		if _, ok := bl.(connectable); bl.ShouldCullAgainst() || ok {
			block = block.Set(d.String(), true)
		} else {
			block = block.Set(d.String(), false)
		}
	}
	return block
}

func (b *blockStainedGlassPane) toData() int {
	if !b.North && !b.South && !b.East && !b.West {
		return int(b.Color)
	}
	return -1
}

// Stairs

type stairHalf int

const (
	shTop stairHalf = iota
	shBottom
)

func (sh stairHalf) String() string {
	switch sh {
	case shTop:
		return "top"
	case shBottom:
		return "bottom"
	}
	return fmt.Sprintf("stairHalf(%d)", sh)
}

type stairShape int

const (
	ssStraight stairShape = iota
	ssInnerLeft
	ssInnerRight
	ssOuterLeft
	ssOuterRight
)

func (sh stairShape) String() string {
	switch sh {
	case ssStraight:
		return "straight"
	case ssInnerLeft:
		return "inner_left"
	case ssInnerRight:
		return "inner_right"
	case ssOuterLeft:
		return "outer_left"
	case ssOuterRight:
		return "outer_right"
	}
	return fmt.Sprintf("stairShape(%d)", sh)
}

type blockStairs struct {
	baseBlock
	Facing direction.Type `state:"facing,2-5"`
	Half   stairHalf      `state:"half,0-1"`
	Shape  stairShape     `state:"shape,0-4"`
}

func (b *blockStairs) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockStairs) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		switch b.Shape {
		case ssStraight:
			b.bounds = []vmath.AABB{
				vmath.NewAABB(0, 0, 0, 1, 0.5, 1),
				vmath.NewAABB(0, 0.5, 0, 1, 1, 0.5),
			}
		case ssInnerLeft:
			b.bounds = []vmath.AABB{
				vmath.NewAABB(0, 0, 0, 1, 0.5, 1),
				vmath.NewAABB(0, 0.5, 0, 1, 1, 0.5),
				vmath.NewAABB(0, 0.5, 0.5, 0.5, 1, 1.0),
			}
		case ssInnerRight:
			b.bounds = []vmath.AABB{
				vmath.NewAABB(0, 0, 0, 1, 0.5, 1),
				vmath.NewAABB(0, 0.5, 0, 1, 1, 0.5),
				vmath.NewAABB(0.5, 0.5, 0.5, 1.0, 1, 1.0),
			}
		case ssOuterLeft:
			b.bounds = []vmath.AABB{
				vmath.NewAABB(0, 0, 0, 1, 0.5, 1),
				vmath.NewAABB(0, 0.5, 0, 0.5, 1, 0.5),
			}
		case ssOuterRight:
			b.bounds = []vmath.AABB{
				vmath.NewAABB(0, 0, 0, 1, 0.5, 1),
				vmath.NewAABB(0.5, 0.5, 0, 1.0, 1, 0.5),
			}
		default:
			b.bounds = []vmath.AABB{
				vmath.NewAABB(0, 0, 0, 1, 1, 1),
			}
		}
		for i := range b.bounds {
			if b.Half == shTop {
				b.bounds[i] = b.bounds[i].RotateX(-math.Pi, 0.5, 0.5, 0.5)
				b.bounds[i] = b.bounds[i].RotateY(-math.Pi, 0.5, 0.5, 0.5)
			}
			switch b.Facing {
			case direction.North:
			case direction.South:
				b.bounds[i] = b.bounds[i].RotateY(-math.Pi, 0.5, 0.5, 0.5)
			case direction.East:
				b.bounds[i] = b.bounds[i].RotateY(-math.Pi*0.5, 0.5, 0.5, 0.5)
			case direction.West:
				b.bounds[i] = b.bounds[i].RotateY(-math.Pi*1.5, 0.5, 0.5, 0.5)
			}
		}
	}
	return b.bounds
}

func (b *blockStairs) ModelVariant() string {
	return fmt.Sprintf("facing=%s,half=%s,shape=%s", b.Facing, b.Half, b.Shape)
}

func (b *blockStairs) UpdateState(x, y, z int) Block {
	// Facing is the side of the back of the stairs
	// If the stair in front of the back doesn't have the
	// same facing as this one or the opposite facing then
	// it will join in the 'outer' shape.
	// If it didn't join with the backface then the front
	// is tested in the same way but forming an 'inner' shape

	ox, oy, oz := b.Facing.Offset()
	if s, ok := chunkMap.Block(x+ox, y+oy, z+oz).(*blockStairs); ok &&
		s.Facing != b.Facing && s.Facing != b.Facing.Opposite() {
		r := false
		if s.Facing == b.Facing.Clockwise() {
			r = true
		}
		if r == (b.Half == shBottom) {
			return b.Set("shape", ssOuterRight)
		}
		return b.Set("shape", ssOuterLeft)
	}

	ox, oy, oz = b.Facing.Opposite().Offset()
	if s, ok := chunkMap.Block(x+ox, y+oy, z+oz).(*blockStairs); ok &&
		s.Facing != b.Facing && s.Facing != b.Facing.Opposite() {
		r := false
		if s.Facing == b.Facing.Clockwise() {
			r = true
		}
		if r == (b.Half == shBottom) {
			return b.Set("shape", ssInnerRight)
		}
		return b.Set("shape", ssInnerLeft)
	}
	return b
}

func (b *blockStairs) toData() int {
	if b.Shape != ssStraight {
		return -1
	}
	data := 0
	switch b.Facing {
	case direction.East:
		data = 0
	case direction.West:
		data = 1
	case direction.South:
		data = 2
	case direction.North:
		data = 3
	}
	if b.Half == shTop {
		data |= 0x4
	}
	return data
}

// Vines

type blockVines struct {
	baseBlock
	Up    bool `state:"up"`
	North bool `state:"north"`
	South bool `state:"south"`
	East  bool `state:"east"`
	West  bool `state:"west"`
}

func (b *blockVines) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockVines) ModelVariant() string {
	return fmt.Sprintf("east=%t,north=%t,south=%t,up=%t,west=%t", b.East, b.North, b.South, b.Up, b.West)
}

func (b *blockVines) UpdateState(x, y, z int) Block {
	if b := chunkMap.Block(x, y+1, z); b.ShouldCullAgainst() {
		return b.Set("up", true)
	}
	return b.Set("up", false)
}

func (b *blockVines) TintImage() *image.NRGBA {
	return foliageBiomeColors
}

func (b *blockVines) toData() int {
	data := 0
	if b.South {
		data |= 0x1
	}
	if b.West {
		data |= 0x2
	}
	if b.North {
		data |= 0x4
	}
	if b.East {
		data |= 0x8
	}
	return data
}

// Stained clay

type blockStainedClay struct {
	baseBlock
	Color color `state:"color,0-15"`
}

func (b *blockStainedClay) ModelName() string {
	return b.Color.String() + "_stained_hardened_clay"
}

func (b *blockStainedClay) toData() int {
	return int(b.Color)
}

// Wool

type blockWool struct {
	baseBlock
	Color color `state:"color,0-15"`
}

func (b *blockWool) ModelName() string {
	return b.Color.String() + "_wool"
}

func (b *blockWool) toData() int {
	return int(b.Color)
}

// Piston

type blockPiston struct {
	baseBlock
	Facing   direction.Type `state:"facing,0-5"`
	Extended bool           `state:"extended"`
}

func (b *blockPiston) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockPiston) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		bo := vmath.NewAABB(0, 0, 0, 1.0, 1.0, 1.0)
		if b.Extended {
			bo.Min[2] = 4.0 / 16.0
		}
		switch b.Facing {
		case direction.North:
		case direction.South:
			bo = bo.RotateY(-math.Pi, 0.5, 0.5, 0.5)
		case direction.West:
			bo = bo.RotateY(-math.Pi*1.5, 0.5, 0.5, 0.5)
		case direction.East:
			bo = bo.RotateY(-math.Pi*0.5, 0.5, 0.5, 0.5)
		case direction.Up:
			bo = bo.RotateX(-math.Pi*1.5, 0.5, 0.5, 0.5)
		case direction.Down:
			bo = bo.RotateX(-math.Pi*0.5, 0.5, 0.5, 0.5)
		}
		b.bounds = []vmath.AABB{bo}
	}
	return b.bounds
}

func (b *blockPiston) LightReduction() int {
	return 6
}

func (b *blockPiston) ModelVariant() string {
	return fmt.Sprintf("extended=%t,facing=%s", b.Extended, b.Facing)
}

func (b *blockPiston) toData() int {
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
	if b.Extended {
		data |= 0x8
	}
	return data
}

type pistonType int

const (
	ptNormal pistonType = iota
	ptSticky
)

func (p pistonType) String() string {
	switch p {
	case ptNormal:
		return "normal"
	case ptSticky:
		return "sticky"
	}
	return fmt.Sprintf("pistonType(%d)", p)
}

type blockPistonHead struct {
	baseBlock
	Facing direction.Type `state:"facing,0-5"`
	Short  bool           `state:"short"`
	Type   pistonType     `state:"type,0-1"`
}

func (b *blockPistonHead) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockPistonHead) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(0, 0, 0, 1.0, 1.0, 4.0/16.0),
			vmath.NewAABB(6.0/16.0, 6.0/16.0, 4.0/16.0, 10.0/16.0, 10.0/16.0, 1.0),
		}
		if !b.Short {
			b.bounds[1].Max[2] += 4.0 / 16.0
		}
		for i := range b.bounds {
			switch b.Facing {
			case direction.North:
			case direction.South:
				b.bounds[i] = b.bounds[i].RotateY(-math.Pi, 0.5, 0.5, 0.5)
			case direction.West:
				b.bounds[i] = b.bounds[i].RotateY(-math.Pi*1.5, 0.5, 0.5, 0.5)
			case direction.East:
				b.bounds[i] = b.bounds[i].RotateY(-math.Pi*0.5, 0.5, 0.5, 0.5)
			case direction.Up:
				b.bounds[i] = b.bounds[i].RotateX(-math.Pi*1.5, 0.5, 0.5, 0.5)
			case direction.Down:
				b.bounds[i] = b.bounds[i].RotateX(-math.Pi*0.5, 0.5, 0.5, 0.5)
			}
		}
	}
	return b.bounds
}

func (b *blockPistonHead) LightReduction() int {
	return 0
}

func (b *blockPistonHead) ModelVariant() string {
	return fmt.Sprintf("facing=%s,short=%t,type=%s", b.Facing, b.Short, b.Type)
}

func (b *blockPistonHead) toData() int {
	if b.Short {
		return -1
	}
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
	if b.Type == ptSticky {
		data |= 0x8
	}
	return data
}

// Slabs

type slabHalf int

const (
	slabTop slabHalf = iota
	slabBottom
)

func (s slabHalf) String() string {
	switch s {
	case slabTop:
		return "top"
	case slabBottom:
		return "bottom"
	}
	return fmt.Sprintf("slabHalf(%d)", s)
}

type slabVariant int

const (
	slabStone slabVariant = iota
	slabSandstone
	slabWooden
	slabCobblestone
	slabBricks
	slabStoneBrick
	slabNetherBrick
	slabQuartz
	slabRedSandstone
	slabOak
	slabSpruce
	slabBirch
	slabJungle
	slabAcacia
	slabDarkOak
)

func (s slabVariant) String() string {
	switch s {
	case slabStone:
		return "stone"
	case slabSandstone:
		return "sandstone"
	case slabWooden:
		return "wood_old"
	case slabCobblestone:
		return "cobblestone"
	case slabBricks:
		return "brick"
	case slabStoneBrick:
		return "stone_brick"
	case slabNetherBrick:
		return "nether_brick"
	case slabQuartz:
		return "quartz"
	case slabRedSandstone:
		return "red_sandstone"
	case slabOak:
		return "oak"
	case slabSpruce:
		return "spruce"
	case slabBirch:
		return "birch"
	case slabJungle:
		return "jungle"
	case slabAcacia:
		return "acacia"
	case slabDarkOak:
		return "dark_oak"
	}
	return fmt.Sprintf("slabVariant(%d)", s)
}

type blockSlab struct {
	baseBlock
	Half    slabHalf    `state:"half,0-1"`
	Variant slabVariant `state:"variant,@TypeRange"`
	Type    string
}

func (b *blockSlab) load(tag reflect.StructTag) {
	b.Type = tag.Get("variant")
	b.cullAgainst = false
}

func (b *blockSlab) TypeRange() (int, int) {
	switch b.Type {
	case "stone":
		return 0, 7
	case "stone2":
		return 8, 8
	case "wood":
		return 9, 14
	}
	panic("invalid type " + b.Type)
}

func (b *blockSlab) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(0, 0, 0, 1.0, 0.5, 1.0),
		}
		if b.Half == slabTop {
			b.bounds[0] = b.bounds[0].Shift(0, 0.5, 0.0)
		}
	}
	return b.bounds
}

func (b *blockSlab) ModelVariant() string {
	return fmt.Sprintf("half=%s", b.Half)
}

func (b *blockSlab) ModelName() string {
	return fmt.Sprintf("%s_slab", b.Variant)
}

func (b *blockSlab) toData() int {
	data := 0
	switch b.Type {
	case "stone":
		data = int(b.Variant)
	case "stone2":
		data = int(b.Variant - 8)
	case "wood":
		data = int(b.Variant - 9)
	}
	if b.Half == slabTop {
		data |= 0x8
	}
	return data
}

type blockSlabDouble struct {
	baseBlock
	Variant slabVariant `state:"variant,@TypeRange"`
	Type    string
}

func (b *blockSlabDouble) load(tag reflect.StructTag) {
	b.Type = tag.Get("variant")
}

func (b *blockSlabDouble) TypeRange() (int, int) {
	switch b.Type {
	case "stone":
		return 0, 7
	case "stone2":
		return 8, 8
	case "wood":
		return 9, 14
	}
	panic("invalid type " + b.Type)
}

func (b *blockSlabDouble) ModelName() string {
	return fmt.Sprintf("%s_double_slab", b.Variant)
}

func (b *blockSlabDouble) toData() int {
	data := 0
	switch b.Type {
	case "stone":
		data = int(b.Variant)
	case "stone2":
		data = int(b.Variant - 8)
	case "wood":
		data = int(b.Variant - 9)
	}
	return data
}

type blockSlabDoubleSeamless struct {
	baseBlock
	Seamless bool        `state:"seamless"`
	Variant  slabVariant `state:"variant,@TypeRange"`
	Type     string
}

func (b *blockSlabDoubleSeamless) load(tag reflect.StructTag) {
	b.Type = tag.Get("variant")
}

func (b *blockSlabDoubleSeamless) TypeRange() (int, int) {
	switch b.Type {
	case "stone":
		return 0, 7
	case "stone2":
		return 8, 8
	case "wood":
		return 9, 14
	}
	panic("invalid type " + b.Type)
}

func (b *blockSlabDoubleSeamless) ModelVariant() string {
	if b.Seamless {
		return "all"
	}
	return "normal"
}

func (b *blockSlabDoubleSeamless) ModelName() string {
	return fmt.Sprintf("%s_double_slab", b.Variant)
}

func (b *blockSlabDoubleSeamless) toData() int {
	data := 0
	switch b.Type {
	case "stone":
		data = int(b.Variant)
	case "stone2":
		data = int(b.Variant - 8)
	case "wood":
		data = int(b.Variant - 9)
	}
	if b.Seamless {
		data |= 0x8
	}
	return data
}

// Carpet

type blockCarpet struct {
	baseBlock
	Color color `state:"color,0-15"`
}

func (b *blockCarpet) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockCarpet) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(0, 0, 0, 1.0, 1.0/16.0, 1.0),
		}
	}
	return b.bounds
}

func (b *blockCarpet) ModelName() string {
	return b.Color.String() + "_carpet"
}

func (b *blockCarpet) toData() int {
	return int(b.Color)
}

// Torch

type blockTorch struct {
	baseBlock
	Facing int `state:"facing,0-4"`
	Model  string
}

func (b *blockTorch) load(tag reflect.StructTag) {
	b.Model = tag.Get("model")
	b.cullAgainst = false
	b.collidable = false
}

func (b *blockTorch) LightEmitted() int {
	return 13
}

func (b *blockTorch) ModelName() string {
	return b.Model
}

func (b *blockTorch) ModelVariant() string {
	facing := b.facing()
	return fmt.Sprintf("facing=%s", facing)
}

func (b *blockTorch) facing() direction.Type {
	switch b.Facing {
	case 0:
		return direction.East
	case 1:
		return direction.West
	case 2:
		return direction.South
	case 3:
		return direction.North
	case 4:
		return direction.Up
	}
	return direction.Invalid
}

func (b *blockTorch) toData() int {
	switch b.facing() {
	case direction.East:
		return 1
	case direction.West:
		return 2
	case direction.South:
		return 3
	case direction.North:
		return 4
	case direction.Up:
		return 5
	}
	return -1
}

// Wall Sign

type blockWallSign struct {
	baseBlock
	Facing direction.Type `state:"facing,2-5"`
}

func (b *blockWallSign) load(tag reflect.StructTag) {
	b.cullAgainst = false
	b.collidable = false
	b.renderable = false
}

func (b *blockWallSign) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(-0.5, -4/16.0, -0.5/16.0, 0.5, 4/16.0, 0.5/16.0),
		}
		f := b.Facing
		ang := float32(0)
		switch f {
		case direction.South:
			ang = math.Pi
		case direction.West:
			ang = math.Pi / 2
		case direction.East:
			ang = -math.Pi / 2
		}
		b.bounds[0] = b.bounds[0].Shift(0.5, 0.5, 0.5-7.5/16.0)
		b.bounds[0] = b.bounds[0].RotateY(ang+math.Pi, 0.5, 0.5, 0.5)
	}
	return b.bounds
}

func (b *blockWallSign) CreateBlockEntity() BlockEntity {
	type wallSign struct {
		blockComponent
		signComponent
	}
	w := &wallSign{}
	w.oz = 7.5 / 16.0
	switch b.Facing {
	case direction.North:
	case direction.South:
		w.rotation = math.Pi
	case direction.West:
		w.rotation = math.Pi / 2
	case direction.East:
		w.rotation = -math.Pi / 2
	}
	return w
}

func (b *blockWallSign) toData() int {
	return int(b.Facing)
}

// Floor Sign

type blockFloorSign struct {
	baseBlock
	Rotation int `state:"rotation,0-15"`
}

func (b *blockFloorSign) load(tag reflect.StructTag) {
	b.cullAgainst = false
	b.collidable = false
	b.renderable = false
}

func (b *blockFloorSign) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(-0.5, -4/16.0, -0.5/16.0, 0.5, 4/16.0, 0.5/16.0),
			vmath.NewAABB(7.5/16.0, 0, 7.5/16.0, 8.5/16.0, 9/16.0, 8.5/16.0),
		}
		b.bounds[0] = b.bounds[0].Shift(0.5, 0.5+5/16.0, 0.5)
		b.bounds[0] = b.bounds[0].RotateY((-float32(b.Rotation)/16)*math.Pi*2+math.Pi, 0.5, 0.5, 0.5)
	}
	return b.bounds
}

func (b *blockFloorSign) CreateBlockEntity() BlockEntity {
	type floorSign struct {
		blockComponent
		signComponent
	}
	w := &floorSign{}
	w.rotation = (-float64(b.Rotation)/16)*math.Pi*2 + math.Pi
	w.oy = 5 / 16.0
	w.hasStand = true
	return w
}

func (b *blockFloorSign) toData() int {
	return b.Rotation
}

// Skull

type blockSkull struct {
	baseBlock
	Facing direction.Type `state:"facing,0-5"`
	NoDrop bool           `state:"nodrop"`
}

func (b *blockSkull) load(tag reflect.StructTag) {
	b.cullAgainst = false
	b.renderable = false
}

func (b *blockSkull) CreateBlockEntity() BlockEntity {
	type skull struct {
		blockComponent
		skullComponent
	}
	w := &skull{}
	w.Facing = b.Facing
	return w
}

func (b *blockSkull) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(0.5-(4/16.0), 0, 0.5-(4/16.0), 0.5+(4/16.0), 8/16.0, 0.5+(4/16.0)),
		}
		f := b.Facing
		if f != direction.Up {
			ang := float32(0)
			switch f {
			case direction.South:
				ang = math.Pi
			case direction.East:
				ang = math.Pi / 2
			case direction.West:
				ang = -math.Pi / 2
			}
			b.bounds[0] = b.bounds[0].Shift(0, 4/16.0, 4/16.0)
			b.bounds[0] = b.bounds[0].RotateY(ang, 0.5, 0.5, 0.5)
		}
	}
	return b.bounds
}

func (b *blockSkull) toData() int {
	data := 0
	switch b.Facing {
	case direction.Up:
		data = 1
	case direction.North:
		data = 2
	case direction.South:
		data = 3
	case direction.East:
		data = 4
	case direction.West:
		data = 5
	}
	if b.NoDrop {
		data |= 0x8
	}
	return data
}

// Portal

type blockPortal struct {
	baseBlock
	Axis blockAxis `state:"axis,1-2"`
}

func (b *blockPortal) load(tag reflect.StructTag) {
	b.cullAgainst = false
	b.collidable = false
	b.translucent = true
}

func (b *blockPortal) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(6/16.0, 0, 0, 10/16.0, 1.0, 1.0),
		}
		if b.Axis == axisX {
			b.bounds[0] = b.bounds[0].RotateY(math.Pi/2, 0.5, 0.5, 0.5)
		}
	}
	return b.bounds
}

func (b *blockPortal) ModelVariant() string {
	return fmt.Sprintf("axis=%s", b.Axis)
}

func (b *blockPortal) toData() int {
	switch b.Axis {
	case axisX:
		return 1
	case axisZ:
		return 2
	}
	return 0
}

// Lilypad

type blockLilypad struct {
	baseBlock
}

func (b *blockLilypad) load(tag reflect.StructTag) {
	b.cullAgainst = false
}

func (b *blockLilypad) CollisionBounds() []vmath.AABB {
	if b.bounds == nil {
		b.bounds = []vmath.AABB{
			vmath.NewAABB(0, 0, 0, 1.0, 1/64.0, 1.0),
		}
	}
	return b.bounds
}

func (b *blockLilypad) TintImage() *image.NRGBA {
	return foliageBiomeColors
}

func (b *blockLilypad) toData() int {
	return 0
}

// Stone brick

type stoneBrickVariant int

const (
	stoneBrickNormal stoneBrickVariant = iota
	stoneBrickMossy
	stoneBrickCracked
	stoneBrickChiseled
)

func (s stoneBrickVariant) String() string {
	switch s {
	case stoneBrickNormal:
		return "stonebrick"
	case stoneBrickMossy:
		return "mossy_stonebrick"
	case stoneBrickCracked:
		return "cracked_stonebrick"
	case stoneBrickChiseled:
		return "chiseled_stonebrick"
	}
	return fmt.Sprintf("stoneBrickVariant(%d)", s)
}

type blockStoneBrick struct {
	baseBlock
	Variant stoneBrickVariant `state:"variant,0-3"`
}

func (b *blockStoneBrick) ModelName() string {
	return b.Variant.String()
}

func (b *blockStoneBrick) toData() int {
	data := int(b.Variant)
	return data
}
