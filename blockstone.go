package main

import "fmt"

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

func (l *blockstone) clone() Block {
	return &blockstone{
		baseBlock: *(l.baseBlock.clone().(*baseBlock)),
		Variant:   l.Variant,
	}
}

func (l *blockstone) ModelName() string {
	return l.Variant.String()
}

func (l *blockstone) toData() int {
	data := int(l.Variant)
	return data
}
