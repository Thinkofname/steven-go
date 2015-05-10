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

package entitysys

import (
	"reflect"
)

type Container struct {
	systems []*system
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) AddEntity(entity interface{}) {
	re := reflect.ValueOf(entity)
	for _, sys := range c.systems {
		if !sys.Matches(entity) {
			continue
		}
		se := &systemEntity{
			v:      re,
			params: make([]reflect.Value, len(sys.params)),
		}
		for i := range sys.params {
			se.params[i] = re
		}

		sys.entities = append(sys.entities, se)
	}
}

func (c *Container) AddSystem(f interface{}, desc ...Desc) {
	s := &system{
		f:    reflect.ValueOf(f),
		desc: desc,
	}
	t := s.f.Type()
	for i := 0; i < t.NumIn(); i++ {
		s.params = append(s.params, t.In(i))
		s.desc = append(s.desc, typeDesc{Type: t.In(i)})
	}
	c.systems = append(c.systems, s)
}

func (c *Container) RemoveEntity(e interface{}) {
	re := reflect.ValueOf(e)
	for _, sys := range c.systems {
	seLoop:
		for i, se := range sys.entities {
			if se.v == re {
				sys.entities = append(sys.entities[:i], sys.entities[i+1:]...)
				break seLoop
			}
		}
	}
}

func (c *Container) Tick() {
	for _, sys := range c.systems {
		for _, e := range sys.entities {
			sys.f.Call(e.params)
		}
	}
}

type system struct {
	f      reflect.Value
	params []reflect.Type
	desc   []Desc

	entities []*systemEntity
}

func (s *system) Matches(e interface{}) bool {
	for _, desc := range s.desc {
		if !desc.Match(e) {
			return false
		}
	}
	return true
}

type systemEntity struct {
	v reflect.Value

	params []reflect.Value
}
