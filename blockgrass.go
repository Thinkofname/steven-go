package main

import "fmt"

type blockGrass struct {
	baseBlock
	Snowy bool
}

func initGrass() *BlockSet {
	g := &blockGrass{}
	g.init("grass")
	set := alloc(g)
	g.Parent = set
	set.supportsData = true
	set.state(newBoolState("snowy"))

	return set
}

func (g *blockGrass) String() string {
	return fmt.Sprintf("%s[snowy=%t]", g.baseBlock.String(), g.Snowy)
}

func (g *blockGrass) clone() Block {
	return &blockGrass{
		baseBlock: *(g.baseBlock.clone().(*baseBlock)),
		Snowy:     g.Snowy,
	}
}

func (g *blockGrass) setState(key string, val interface{}) {
	switch key {
	case "snowy":
		g.Snowy = val.(bool)
	default:
		panic("invalid state " + key)
	}
}

func (g *blockGrass) ModelVariant() string {
	return fmt.Sprintf("snowy=%t", g.Snowy)
}

func (g *blockGrass) toData() int {
	if g.Snowy {
		return -1
	}
	return 0
}
