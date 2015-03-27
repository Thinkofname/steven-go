package main

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
	if lava {
		l.color = 0xFF0000
	} else {
		l.color = 0x0000FF
	}
	set.supportsData = true
	set.state(newIntState("level", 0, 15))

	return set
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
		l.color = 0xFF - uint32(l.Level*(0xFF/15))
	default:
		panic("invalid state " + key)
	}
}

func (l *BlockLiquid) toData() int {
	return l.Level
}
