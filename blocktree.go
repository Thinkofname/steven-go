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
)

// Logs have 4 possible 'rotations' (more like shapes),
// one for each possible axis and a special 'none' one.
// The 'none' state causes the bark texture to be displayed
// on all sides.
type blockAxis int

const (
	axisY blockAxis = iota
	axisZ
	axisX
	axisNone
)

func (b blockAxis) String() string {
	switch b {
	case axisNone:
		return "none"
	case axisX:
		return "x"
	case axisY:
		return "y"
	case axisZ:
		return "z"
	}
	return fmt.Sprintf("blockAxis(%d)", b)
}

// The tree variants are split across two log/leaves types
// due to the size limits of the old data field.
type treeVariant int

const (
	treeOak treeVariant = iota
	treeSpruce
	treeBirch
	treeJungle
	treeAcacia
	treeDarkOak
)

func (t treeVariant) String() string {
	switch t {
	case treeOak:
		return "oak"
	case treeSpruce:
		return "spruce"
	case treeBirch:
		return "birch"
	case treeJungle:
		return "jungle"
	case treeAcacia:
		return "acacia"
	case treeDarkOak:
		return "dark_oak"
	}
	return fmt.Sprintf("treeVariant(%d)", t)
}

type blockLog struct {
	baseBlock
	Variant treeVariant `state:"variant,@VariantRange"`
	second  bool
	Axis    blockAxis `state:"axis,0-3"`
}

func initLog(name string, second bool) *BlockSet {
	l := &blockLog{}
	l.init(name)
	l.second = second
	set := alloc(l)
	return set
}

func (l *blockLog) VariantRange() (int, int) {
	if l.second {
		return 4, 5
	}
	return 0, 3
}

func (l *blockLog) String() string {
	return l.Parent.stringify(l)
}

func (l *blockLog) clone() Block {
	return &blockLog{
		baseBlock: *(l.baseBlock.clone().(*baseBlock)),
		Variant:   l.Variant,
		second:    l.second,
		Axis:      l.Axis,
	}
}

func (l *blockLog) ModelName() string {
	return l.Variant.String() + "_log"
}

func (l *blockLog) ModelVariant() string {
	return fmt.Sprintf("axis=%s", l.Axis)
}

func (l *blockLog) toData() int {
	data := int(l.Variant)
	if l.second {
		data -= 4
	}
	data |= int(l.Axis) << 2
	return data
}

type blockLeaves struct {
	baseBlock
	Variant    treeVariant `state:"variant,@VariantRange"`
	second     bool
	Decayable  bool `state:"decayable"`
	CheckDecay bool `state:"check_decay"`
}

func initLeaves(name string, second bool) *BlockSet {
	l := &blockLeaves{}
	l.init(name)
	l.second = second
	l.cullAgainst = false
	set := alloc(l)
	return set
}

func (l *blockLeaves) VariantRange() (int, int) {
	if l.second {
		return 4, 5
	}
	return 0, 3
}

func (l *blockLeaves) String() string {
	return l.Parent.stringify(l)
}

func (l *blockLeaves) clone() Block {
	return &blockLeaves{
		baseBlock:  *(l.baseBlock.clone().(*baseBlock)),
		Variant:    l.Variant,
		second:     l.second,
		Decayable:  l.Decayable,
		CheckDecay: l.CheckDecay,
	}
}

func (l *blockLeaves) ModelName() string {
	return l.Variant.String() + "_leaves"
}

func (l *blockLeaves) ForceShade() bool {
	return true
}

func (l *blockLeaves) TintImage() *image.NRGBA {
	return foliageBiomeColors
}

func (l *blockLeaves) toData() int {
	data := int(l.Variant)
	if l.second {
		data -= 4
	}
	if l.Decayable {
		data |= 0x4
	}
	if l.CheckDecay {
		data |= 0x8
	}
	return data
}
