package chat

import (
	"bytes"
	"fmt"
)

// String provides a string version of the component without
// formatting.
func (a AnyComponent) String() string {
	return fmt.Sprint(a.Value)
}

// String provides a string version of the component without
// formatting.
func (c *Component) String() string {
	var buf bytes.Buffer
	for _, e := range c.Extra {
		fmt.Fprint(&buf, e)
	}
	return buf.String()
}

// String provides a string version of the component without
// formatting.
func (t *TextComponent) String() string {
	return t.Text + t.Component.String()
}
