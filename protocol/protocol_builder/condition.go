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

package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

type conditions []condition

type condition struct {
	l    conditionVar
	cond string
	r    conditionVar
}

type conditionVar struct {
	name        string
	structField int
}

func parseCondition(con string) conditions {
	var conds conditions
	for len(con) > 0 {
		v := conditionVar{}
		for con[0] == '.' {
			con = con[1:]
			v.structField++
		}
		pos := strings.IndexFunc(con, func(r rune) bool { return !unicode.In(r, unicode.Letter, unicode.Digit) })
		if pos == -1 {
			panic("invalid condition")
		}
		v.name = con[:pos]
		con = con[pos:]
		con = strings.TrimSpace(con)

		c := condition{l: v}

		pos = strings.IndexFunc(con, func(r rune) bool { return !unicode.IsSymbol(r) && r != '!' })
		c.cond = con[:pos]
		con = con[pos:]
		con = strings.TrimSpace(con)

		r := conditionVar{}
		for con[0] == '.' {
			r.structField++
			con = con[1:]
		}

		pos = strings.IndexFunc(con, func(r rune) bool { return !unicode.In(r, unicode.Letter, unicode.Digit) && r != '"' })
		if pos == -1 {
			pos = len(con)
		}
		r.name = con[:pos]
		c.r = r
		conds = append(conds, c)
		con = con[pos:]
		con = strings.TrimSpace(con)
	}
	return conds
}

func (c *condition) print(base string, buf *bytes.Buffer) {
	l := c.l.name
	if c.l.structField > 0 {
		l = base
		i := c.l.structField - 1
		for i > 0 {
			p := strings.LastIndex(l, ".")
			l = l[:p]
			i--
		}
		l += "." + c.l.name
	}
	r := c.r.name
	if c.r.structField > 0 {
		r = base
		i := c.r.structField - 1
		for i > 0 {
			p := strings.LastIndex(r, ".")
			r = r[:p]
			i--
		}
		r += "." + c.r.name
	}
	fmt.Fprintf(buf, "%s %s %s", l, c.cond, r)
}

func (c *condition) equals(o *condition) bool {
	if c == o {
		return true
	}
	if c == nil || o == nil {
		return false
	}
	return *c == *o
}

func (c conditions) equals(o conditions) bool {
	if c == nil || o == nil {
		return false
	}
	if len(c) != len(o) {
		return false
	}
	for i, cond := range c {
		if !cond.equals(&o[i]) {
			return false
		}
	}
	return true
}

func (c conditions) print(base string, buf *bytes.Buffer) {
	buf.WriteString("if ")
	for i, cond := range c {
		cond.print(base, buf)
		if i != len(c)-1 {
			buf.WriteString(" || ")
		}
	}
	buf.WriteString(" {\n")
}
