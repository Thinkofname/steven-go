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
	"strings"
	"testing"
)

func TestTypeString(t *testing.T) {
	r := &registry{}

	called := false
	r.Register("testing %", func(arg string) {
		if arg != "hello" {
			t.FailNow()
		}
		called = true
	})

	checkError(t, r.Execute("testing hello"))
	if !called {
		t.FailNow()
	}
}

func TestTypeLimit(t *testing.T) {
	r := &registry{}

	r.Register("testing %3", func(arg string) {
		t.FailNow()
	})

	err := r.Execute("testing hello")
	if err == nil || !strings.Contains(err.Error(), "too long") {
		t.FailNow()
	}
}

func TestTypeSub(t *testing.T) {
	r := &registry{}

	called := 0
	r.Register("testing % a", func(arg string) {
		if arg != "hello" {
			t.FailNow()
		}
		if called != 0 {
			t.FailNow()
		}
		called++
	})
	r.Register("testing % b", func(arg string) {
		if arg != "hello" {
			t.FailNow()
		}
		if called != 1 {
			t.FailNow()
		}
		called++
	})

	checkError(t, r.Execute("testing hello a"))
	checkError(t, r.Execute("testing hello b"))
	if called != 2 {
		t.FailNow()
	}
}
