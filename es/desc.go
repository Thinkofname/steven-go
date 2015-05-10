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

package es

import (
	"reflect"
)

type Desc interface {
	Match(e interface{}) bool
}

type typeDesc struct {
	Type reflect.Type
}

func (t typeDesc) Match(e interface{}) bool { return reflect.TypeOf(e).Implements(t.Type) }

func Type(t interface{}) Desc {
	return typeDesc{
		Type: reflect.TypeOf(t).Elem(),
	}
}

type notDesc struct {
	child Desc
}

func (n notDesc) Match(e interface{}) bool { return !n.child.Match(e) }

func Not(d Desc) Desc {
	return notDesc{child: d}
}
