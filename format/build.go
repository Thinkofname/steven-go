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

// Builder allows for easy creation of simple messages.
type Builder struct {
	main *TextComponent

	cur *TextComponent
}

// Build creates a builder starting with the initial message.
func Build(msg string) *Builder {
	return &Builder{
		main: &TextComponent{},
		cur:  &TextComponent{Text: msg},
	}
}

// Append adds the message to the builder and sets it as the current
// target.
func (b *Builder) Append(msg string) *Builder {
	b.main.Extra = append(b.main.Extra, Wrap(b.cur))
	old := b.cur
	b.cur = &TextComponent{Text: msg}
	b.cur.Color = old.Color
	return b
}

// Color changes the color of current message
func (b *Builder) Color(c Color) *Builder {
	b.cur.Color = c
	return b
}

// Create returns the component created by this builder
func (b *Builder) Create() AnyComponent {
	b.main.Extra = append(b.main.Extra, Wrap(b.cur))
	b.cur = nil
	return Wrap(b.main)
}
