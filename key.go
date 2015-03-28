package main

type pluginKey struct {
	Plugin, Name string
}

func (pk pluginKey) String() string {
	return pk.Plugin + ":" + pk.Name
}
