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

package console

import (
	"fmt"

	"github.com/thinkofdeath/steven/chat"
)

// IntVar is a console var that contains an integer
type IntVar struct {
	name  string
	Value int
}

// NewIntVar creates and registers a integer console variable
func NewIntVar(name string, val int) *IntVar {
	i := &IntVar{
		name:  name,
		Value: val,
	}
	Register(fmt.Sprintf("%s", name), func() {
		i.print()
	})
	Register(fmt.Sprintf("%s = %%", name), func(v int) {
		i.Value = v
		i.print()
	})
	return i
}

func (i *IntVar) print() {
	Component(chat.Build(i.name).
		Color(chat.Aqua).
		Append(" = ").
		Append(fmt.Sprint(i.Value)).
		Color(chat.Aqua).
		Create(),
	)
}

// StringVar is a console var that contains an string
type StringVar struct {
	name  string
	Value string
}

// NewStringVar creates and registers a string console variable
func NewStringVar(name string, val string) *StringVar {
	s := &StringVar{
		name:  name,
		Value: val,
	}
	Register(fmt.Sprintf("%s", name), func() {
		s.print()
	})
	Register(fmt.Sprintf("%s = %%", name), func(v string) {
		s.Value = v
		s.print()
	})
	return s
}

func (s *StringVar) print() {
	Component(chat.Build(s.name).
		Color(chat.Aqua).
		Append(" = ").
		Append("\"").Color(chat.Yellow).
		Append(s.Value).
		Color(chat.Aqua).
		Append("\"").Color(chat.Yellow).
		Create(),
	)
}
