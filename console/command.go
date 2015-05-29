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
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/thinkofdeath/steven/chat"
)

// Registry contains information required to store
// and execute commands
//
// ExtraParameters can be set to specify how many
// extra parameters a command function will have
// after the caller argument
type registry struct {
	ExtraParameters int

	root *commandNode

	typeHandlers map[reflect.Type]TypeHandler
}

var (
	// ErrCommandNotFound is returned when no matching command was found
	ErrCommandNotFound = errors.New("command not found")
)

var quotedStringRegex = regexp.MustCompile(`[^\s"]+|"([^"]*)"`)

// Register adds the passed function to the command registry
// using the description to decided on its name and location.
//
// f must be a function
//
// This is designed to panic instead of returning an error
// because its intended to be used in init methods and fail
// early on mistakes
func Register(desc string, f interface{}) {
	defaultRegistry.Register(desc, f)
}

func (r *registry) Register(desc string, f interface{}) {
	if r.typeHandlers == nil {
		r.initTypes()
	}

	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		panic("f must be a function")
	}
	args := strings.Split(desc, " ")

	if len(args[0]) < 1 {
		panic("Invalid command desc")
	}

	if r.root == nil {
		r.root = &commandNode{
			childNodes: map[string]*commandNode{},
		}
	}

	current := r.root
	pos := 0
main:
	for _, arg := range args {
		if arg[0] == '%' {
			// Make sure the function can actually handle this
			offset := r.ExtraParameters + pos
			if offset >= t.NumIn() {
				panic("not enough parameters for function")
			}
			pos++
			arg = arg[1:]
			// Get the handler
			aT := t.In(offset)
			handler, ok := r.typeHandlers[aT]
			if !ok {
				panic(fmt.Errorf("no handler for %s", aT))
			}
			data := handler.DefineType(arg)

			// Check for an existing entry
			for _, info := range current.types {
				if info.handler != handler {
					continue
				}
				if handler.Equals(data, info.data) {
					current = info.node
					continue main
				}
			}

			info := typeInfo{
				handler: handler,
				data:    data,
				node: &commandNode{
					childNodes: map[string]*commandNode{},
				},
			}
			current.types = append(current.types, info)
			current = info.node
			continue
		}
		current = current.subNode(arg)
	}
	if r.ExtraParameters+pos != t.NumIn() {
		panic("too many parameters for function")
	}

	if current.f != nil {
		panic("Double registered command")
	}

	current.f = f
}

// Execute tries to execute the specified command.
//
// Panics if the number of extra arguments doesn't match the
// amount specified in Registry's ExtraParameters
func Execute(cmd string, extra ...interface{}) (err error) {
	Component(chat.
		Build("> ").
		Color(chat.Yellow).
		Append(cmd).
		Create(),
	)
	return defaultRegistry.Execute(cmd, extra...)
}

func (r *registry) Execute(cmd string, extra ...interface{}) (err error) {
	if len(extra) != r.ExtraParameters {
		panic("Incorrect number of extra parameters")
	}

	if r.root == nil {
		return ErrCommandNotFound
	}
	// Catch and return any errors thrown to prevent crashing
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	// Unlike most command systems this supports quoting of
	// arguments using ". The regex used leaves the quotes
	// in the resulting string so we go through a strip them
	parts := quotedStringRegex.FindAllString(cmd, -1)
	for i, p := range parts {
		if strings.HasPrefix(p, `"`) && strings.HasSuffix(p, `"`) {
			parts[i] = p[1 : len(p)-1]
		}
	}

	current, vals, err := r.exec(r.root, parts)
	if err != nil {
		return err
	}

	// Its possible to reach a node which doesn't have
	// a command assigned (e.g. part of a sub-command)
	// so we have to check that too
	if current != nil && current.f != nil {
		// No checks are preformed on the function here
		// as they should have been checked in Register
		f := reflect.ValueOf(current.f)
		args := make([]reflect.Value, r.ExtraParameters+len(vals))
		for i, e := range extra {
			args[i] = reflect.ValueOf(e)
		}
		copy(args[r.ExtraParameters:], vals)
		f.Call(args)
		return nil
	}

	return ErrCommandNotFound
}

func (r *registry) exec(node *commandNode, args []string, vals ...reflect.Value) (*commandNode, []reflect.Value, error) {
	if len(args) == 0 {
		return node, vals, nil
	}
	part := args[0]

	// Try types first
	var err error

	for _, info := range node.types {
		var a interface{}
		a, err = info.handler.ParseType(part, info.data)
		if err == nil {
			var cn *commandNode
			var v []reflect.Value
			cn, v, err = r.exec(info.node, args[1:], append(vals, reflect.ValueOf(a))...)
			if err == nil {
				return cn, v, nil
			}
		}
	}

	if cn, ok := node.childNodes[strings.ToLower(part)]; ok {
		return r.exec(cn, args[1:], vals...)
	}
	if err == nil {
		err = ErrCommandNotFound
	}
	return nil, vals, err
}

type commandNode struct {
	childNodes map[string]*commandNode
	types      []typeInfo
	f          interface{}
}

type typeInfo struct {
	handler TypeHandler
	data    interface{}
	node    *commandNode
}

func (cn *commandNode) subNode(name string) *commandNode {
	name = strings.ToLower(name)
	node, ok := cn.childNodes[name]
	if !ok {
		node = &commandNode{
			childNodes: map[string]*commandNode{},
		}
		cn.childNodes[name] = node
	}
	return node
}
