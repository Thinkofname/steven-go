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

// Package glsl handles glsl shaders and applies some preprocessing to them
package glsl

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

var shaders = map[string]string{}

// Register adds the named shader to the database.
// Once included they may be obtained from Get or
// used within other shaders via #include
func Register(name, source string) {
	_, ok := shaders[name]
	if ok {
		panic("shader " + name + " is already defined")
	}
	shaders[name] = strings.TrimSpace(source)
}

// Get returns the named shader optionally adding
// the passed defines. The output from this will be
// processed finding and expanding #include statements.
// This also adds debugging information to the shaders
// (line number and name).
func Get(name string, defines ...string) string {
	var out bytes.Buffer
	out.WriteString("#version 150\n")
	for _, def := range defines {
		fmt.Fprintf(&out, "#define %s\n", def)
	}
	get(&out, name)
	return out.String()
}

func get(out *bytes.Buffer, name string) {
	src, ok := shaders[name]
	if !ok {
		panic("unknown shader " + name)
	}
	s := bufio.NewScanner(strings.NewReader(src))
	lineNo := 0
	for s.Scan() {
		lineNo++
		line := s.Text()
		if strings.HasPrefix(line, "#include ") {
			inc := strings.TrimSpace(line[len("#include "):])
			get(out, inc)
			continue
		}
		out.WriteString(line)
		out.WriteRune('\n')
		if os.Getenv("STEVEN_DEBUG") == "true" {
			fmt.Fprintf(out, "#line %d \"%s\"\n", lineNo, name)
		}
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
}
