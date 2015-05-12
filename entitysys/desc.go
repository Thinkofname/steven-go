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

import "reflect"

// Matcher is used to select when an entity is used for a
// system.
type Matcher interface {
	Match(e interface{}) bool
}

type typeMatcher struct {
	Type reflect.Type
}

func (t typeMatcher) Match(e interface{}) bool {
	if t.Type.Kind() == reflect.Interface {
		return reflect.TypeOf(e).Implements(t.Type)
	}
	_, ok := reflect.TypeOf(e).Elem().FieldByName(t.Type.Elem().Name())
	return ok
}

// Type returns a Matcher that matches when the type of element
// of the passed to this is implemented by the entity.
//
// This method should be used as followed for interfaces
//     Type((*MyInterface)(nil))
func Type(t interface{}) Matcher {
	return typeMatcher{
		Type: reflect.TypeOf(t).Elem(),
	}
}

type notMatcher struct {
	child Matcher
}

func (n notMatcher) Match(e interface{}) bool { return !n.child.Match(e) }

// Not inverts the passed Matcher.
func Not(d Matcher) Matcher {
	return notMatcher{child: d}
}
