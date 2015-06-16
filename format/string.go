// Copyright 2015 Matthew Collins
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package format

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
