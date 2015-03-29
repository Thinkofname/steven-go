package main

import "image"

type blockLeaves struct {
	baseBlock
	Variant    treeVariant
	Decayable  bool `state:"decayable"`
	CheckDecay bool `state:"check_decay"`
}

func initLeaves(name string) *BlockSet {
	l := &blockLeaves{}
	l.init(name)
	l.cullAgainst = false
	set := alloc(l)
	return set
}

func (l *blockLeaves) String() string {
	return l.Parent.stringify(l)
}

func (l *blockLeaves) clone() Block {
	return &blockLeaves{
		baseBlock:  *(l.baseBlock.clone().(*baseBlock)),
		Variant:    l.Variant,
		Decayable:  l.Decayable,
		CheckDecay: l.CheckDecay,
	}
}

func (l *blockLeaves) ModelName() string {
	return l.Variant.String() + "_" + l.name
}

func (l *blockLeaves) ForceShade() bool {
	return true
}

func (l *blockLeaves) TintImage() *image.NRGBA {
	return foliageBiomeColors
}

func (l *blockLeaves) toData() int {
	data := int(l.Variant)
	if l.Decayable {
		data |= 0x4
	}
	if l.CheckDecay {
		data |= 0x8
	}
	return data
}
