package main

type blockSimple struct {
	baseBlock
}

type simpleConfig struct {
	NotCullAgainst bool
}

func initSimple(name string) *BlockSet {
	return initSimpleConfig(name, simpleConfig{})
}

func initSimpleConfig(name string, config simpleConfig) *BlockSet {
	s := &blockSimple{}
	s.init(name)
	set := alloc(s)

	s.cullAgainst = !config.NotCullAgainst

	return set
}

func (b *blockSimple) toData() int {
	if b == b.Parent.Base {
		return 0
	}
	return -1
}
