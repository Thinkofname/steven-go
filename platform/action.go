package platform

//go:generate stringer -type=Action

// Action is something that a platform implementation
// can trigger.
type Action int

const (
	InvalidAction Action = iota
	Debug
)
