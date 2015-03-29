package main

type blockLiquid struct {
	baseBlock
	Level int `state:"level,0-15"`
}

func initLiquid(name string, lava bool) *BlockSet {
	l := &blockLiquid{}
	l.init(name)
	set := alloc(l)
	return set
}

func (l *blockLiquid) String() string {
	return l.Parent.stringify(l)
}

func (l *blockLiquid) clone() Block {
	return &blockLiquid{
		baseBlock: *(l.baseBlock.clone().(*baseBlock)),
		Level:     l.Level,
	}
}

func (l *blockLiquid) toData() int {
	return l.Level
}
