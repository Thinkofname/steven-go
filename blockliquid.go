package main

import "fmt"

// BlockLiquid is a liquid based block e.g. lava and water
type BlockLiquid struct {
	baseBlock
	Level int
}

func initLiquid(name string, lava bool) *BlockSet {
	l := &BlockLiquid{}
	l.init(name)
	set := alloc(l)
	l.Parent = set
	set.supportsData = true
	set.state(newIntState("level", 0, 15))

	return set
}

func (l *BlockLiquid) String() string {
	return fmt.Sprintf("%s[level=%d]", l.baseBlock.String(), l.Level)
}

func (l *BlockLiquid) clone() Block {
	return &BlockLiquid{
		baseBlock: *(l.baseBlock.clone().(*baseBlock)),
		Level:     l.Level,
	}
}

func (l *BlockLiquid) setState(key string, val interface{}) {
	switch key {
	case "level":
		l.Level = val.(int)
	default:
		panic("invalid state " + key)
	}
}

func (l *BlockLiquid) toData() int {
	return l.Level
}
