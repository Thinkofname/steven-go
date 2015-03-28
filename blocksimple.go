package main

// BlockSimple is a basic block with no special attributes
type BlockSimple struct {
	baseBlock
}

type simpleConfig struct {
	NotCullAgainst bool
}

func initSimple(name string, config simpleConfig) *BlockSet {
	s := &BlockSimple{}
	s.init(name)
	set := alloc(s)
	s.Parent = set

	s.cullAgainst = !config.NotCullAgainst

	return set
}
