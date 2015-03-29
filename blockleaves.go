package main

import "fmt"

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
	return fmt.Sprintf("%s[variant=%s,decayable=%t,check_decay=%t]", l.baseBlock.String(), l.Variant, l.Decayable, l.CheckDecay)
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

func (l *blockLeaves) toData() int {
	data := 0
	switch l.Variant {
	case treeOak:
		data = 0
	case treeSpruce:
		data = 1
	case treeBirch:
		data = 2
	case treeJungle:
		data = 3
	}
	if l.Decayable {
		data |= 0x4
	}
	if l.CheckDecay {
		data |= 0x8
	}
	return data
}
