package main

// BlockSimple is a basic block with no special attributes
type BlockSimple struct {
	baseBlock
}

type simpleConfig struct {
	Color uint32
}

func initSimple(name string, config simpleConfig) *BlockSet {
	s := &BlockSimple{}
	s.init(name)
	set := alloc(s)
	s.Parent = set

	s.color = config.Color
	return set
}
