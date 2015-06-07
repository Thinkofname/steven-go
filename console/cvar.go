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
	"bytes"
	"fmt"
	"strings"

	"github.com/thinkofdeath/steven/chat"
)

var cvars []cvarI

// Property is a flag for cvar to control how its handled
type Property int

const (
	// Mutable marks whether the var can be changed from the
	// console
	Mutable Property = 1 << iota
	// Serializable marks the var to be saved and loaded
	Serializable
)

func (p Property) is(po Property) bool { return p&po == po }

type cvar struct {
	properties Property

	documentation string
}

type cvarI interface {
	props() Property
}

func (c *cvar) props() Property { return c.properties }

// Doc sets the documentation for this cvar
func (c *cvar) Doc(d string) { c.documentation = strings.TrimSpace(d) }

func (c *cvar) String() string {
	var buf bytes.Buffer
	for _, line := range strings.Split(c.documentation, "\n") {
		buf.WriteString("# ")
		buf.WriteString(line)
		buf.WriteRune('\n')
	}
	return buf.String()
}

func (c *cvar) printDoc() {
	for _, line := range strings.Split(c.documentation, "\n") {
		Component(chat.Build("# ").
			Color(chat.DarkGray).
			Append(line).
			Color(chat.Gray).
			Create())
	}
}

// IntVar is a console var that contains an integer
type IntVar struct {
	cvar
	name  string
	value int
}

// NewIntVar creates and registers a integer console variable
func NewIntVar(name string, val int, props ...Property) *IntVar {
	i := &IntVar{
		name:  name,
		value: val,
	}
	cvars = append(cvars, i)
	for _, p := range props {
		i.properties |= p
	}
	Register(fmt.Sprintf("%s", name), func() {
		i.printDoc()
		i.print()
	})
	if i.properties.is(Mutable) {
		Register(fmt.Sprintf("%s %%", name), func(v int) {
			i.SetValue(v)
			i.print()
		})
	}
	return i
}

func (i *IntVar) String() string {
	var buf bytes.Buffer
	buf.WriteString(i.cvar.String())
	buf.WriteString(i.name)
	buf.WriteRune(' ')
	fmt.Fprint(&buf, i.value)
	return buf.String()
}

// Doc sets the documentation for this cvar
func (i *IntVar) Doc(d string) *IntVar { i.cvar.Doc(d); return i }

func (i *IntVar) Value() int { return i.value }
func (i *IntVar) SetValue(v int) {
	i.value = v
	if i.properties.is(Serializable) {
		saveConf()
	}
}

func (i *IntVar) print() {
	Component(chat.Build(i.name).
		Color(chat.Aqua).
		Append(" ").
		Append(fmt.Sprint(i.value)).
		Color(chat.Aqua).
		Create(),
	)
}

// StringVar is a console var that contains an string
type StringVar struct {
	cvar
	name  string
	value string
}

// NewStringVar creates and registers a string console variable
func NewStringVar(name, val string, props ...Property) *StringVar {
	s := &StringVar{
		name:  name,
		value: val,
	}
	cvars = append(cvars, s)
	for _, p := range props {
		s.properties |= p
	}
	Register(fmt.Sprintf("%s", name), func() {
		s.printDoc()
		s.print()
	})
	if s.properties.is(Mutable) {
		Register(fmt.Sprintf("%s %%", name), func(v string) {
			s.value = v
			s.print()
		})
	}
	return s
}

func (s *StringVar) String() string {
	var buf bytes.Buffer
	buf.WriteString(s.cvar.String())
	buf.WriteString(s.name)
	buf.WriteRune(' ')
	fmt.Fprintf(&buf, "\"%s\"", s.value)
	return buf.String()
}

// Doc sets the documentation for this cvar
func (s *StringVar) Doc(d string) *StringVar { s.cvar.Doc(d); return s }
func (s *StringVar) Value() string           { return s.value }
func (s *StringVar) SetValue(v string) {
	s.value = v
	if s.properties.is(Serializable) {
		saveConf()
	}
}

func (s *StringVar) print() {
	Component(chat.Build(s.name).
		Color(chat.Aqua).
		Append(" ").
		Append("\"").Color(chat.Yellow).
		Append(s.value).
		Color(chat.Aqua).
		Append("\"").Color(chat.Yellow).
		Create(),
	)
}
