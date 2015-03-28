package main

import "fmt"

type Axis int

const (
	AxisY Axis = iota
	AxisZ
	AxisX
	AxisNone
)

type blockLog struct {
	baseBlock
	Variant string
	Axis    Axis
}

func initLog(name string) *BlockSet {
	l := &blockLog{}
	l.init(name)
	set := alloc(l)
	l.Parent = set
	set.supportsData = true
	set.state(newIntState("variant", 0, 3))
	set.state(newIntState("axis", 0, 3))

	return set
}

func (l *blockLog) String() string {
	var axis string
	switch l.Axis {
	case AxisNone:
		axis = "none"
	case AxisX:
		axis = "x"
	case AxisY:
		axis = "y"
	case AxisZ:
		axis = "z"
	}
	return fmt.Sprintf("%s[variant=%s,axis=%s]", l.baseBlock.String(), l.Variant, axis)
}

func (l *blockLog) clone() Block {
	return &blockLog{
		baseBlock: *(l.baseBlock.clone().(*baseBlock)),
		Variant:   l.Variant,
		Axis:      l.Axis,
	}
}

func (l *blockLog) setState(key string, val interface{}) {
	switch key {
	case "variant":
		l.Variant = logVariant(l.name, val.(int))
	case "axis":
		l.Axis = Axis(val.(int))
	default:
		panic("invalid state " + key)
	}
}

func (l *blockLog) ModelName() string {
	return l.Variant + "_" + l.name
}

func (l *blockLog) ModelVariant() string {
	var axis string
	switch l.Axis {
	case AxisNone:
		axis = "none"
	case AxisX:
		axis = "x"
	case AxisY:
		axis = "y"
	case AxisZ:
		axis = "z"
	}
	return fmt.Sprintf("axis=%s", axis)
}

func (l *blockLog) toData() int {
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
	data |= int(l.Axis) << 2
	return data
}

func logVariant(name string, val int) string {
	switch name {
	case "log", "leaves":
		switch val {
		case 0:
			return "oak"
		case 1:
			return "spruce"
		case 2:
			return "birch"
		case 3:
			return "jungle"
		}
	}
	panic("unsupported log " + name)
}
