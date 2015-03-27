package main

type stateCollection interface {
	key() string
	values() []interface{}
}

type intState struct {
	name     string
	min, max int
}

func newIntState(name string, min, max int) intState {
	return intState{
		name: name,
		min:  min, max: max,
	}
}

func (is intState) key() string {
	return is.name
}

func (is intState) values() []interface{} {
	vals := make([]interface{}, is.max-is.min+1)
	for i := is.min; i <= is.max; i++ {
		vals[i-is.min] = i
	}
	return vals
}
