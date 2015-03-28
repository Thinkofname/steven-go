package main

import "fmt"

type blockLeaves struct {
	baseBlock
	Variant    string
	Decayable  bool
	CheckDecay bool
}

func initLeaves(name string) *BlockSet {
	l := &blockLeaves{}
	l.init(name)
	set := alloc(l)
	l.Parent = set

	l.cullAgainst = false

	set.supportsData = true
	set.state(newIntState("variant", 0, 3))
	set.state(newBoolState("decayable"))
	set.state(newBoolState("check_decay"))

	return set
}

func (l *blockLeaves) String() string {
	return fmt.Sprintf("%s[variant=%s,decayable=%t,check_decay]", l.baseBlock.String(), l.Variant, l.Decayable, l.CheckDecay)
}

func (l *blockLeaves) clone() Block {
	return &blockLeaves{
		baseBlock:  *(l.baseBlock.clone().(*baseBlock)),
		Variant:    l.Variant,
		Decayable:  l.Decayable,
		CheckDecay: l.CheckDecay,
	}
}

func (l *blockLeaves) setState(key string, val interface{}) {
	switch key {
	case "variant":
		l.Variant = logVariant(l.name, val.(int))
	case "decayable":
		l.Decayable = val.(bool)
	case "check_decay":
		l.CheckDecay = val.(bool)
	default:
		panic("invalid state " + key)
	}
}

func (l *blockLeaves) ModelName() string {
	return l.Variant + "_" + l.name
}

func (l *blockLeaves) toData() int {
	data := 0
	switch l.Variant {
	case "oak":
		data = 0
	case "spurce":
		data = 1
	case "birch":
		data = 2
	case "jungle":
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
