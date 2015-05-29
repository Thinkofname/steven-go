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

import "testing"

func TestBasic(t *testing.T) {
	r := registry{}

	called := false
	r.Register("test", func() {
		called = true
	})
	checkError(t, r.Execute("test"))
	if !called {
		t.FailNow()
	}
}

func TestSubCommands(t *testing.T) {
	r := registry{}

	called := 0
	r.Register("test a", func() {
		called++
	})
	r.Register("test b", func() {
		called++
	})
	r.Register("test c a", func() {
		called++
	})
	r.Register("test c b", func() {
		called++
	})
	r.Register("test d e f g", func() {
		called++
	})

	checkError(t, r.Execute("test a"))
	checkError(t, r.Execute("test b"))
	checkError(t, r.Execute("test c a"))
	checkError(t, r.Execute("test c b"))
	checkError(t, r.Execute("test d e f g"))

	if called != 5 {
		t.FailNow()
	}
}

func TestNonFunction(t *testing.T) {
	shouldPanic(t, func() {
		r := registry{}
		r.Register("test a", "")
	})
}

func TestInvalidDesc(t *testing.T) {
	shouldPanic(t, func() {
		r := registry{}
		r.Register("", func() {

		})
	})
}

func TestDoubleAdd(t *testing.T) {
	shouldPanic(t, func() {
		r := registry{}
		r.Register("test", func() {

		})
		r.Register("test", func() {

		})
	})
}

func TestExtraParams(t *testing.T) {
	r := registry{
		ExtraParameters: 2,
	}
	r.Register("extra", func(a, b string) {
		if a != "a" || b != "b" {
			t.FailNow()
		}
	})
	checkError(t, r.Execute("extra", "a", "b"))
}

func TestExtraParamsFail(t *testing.T) {
	shouldPanic(t, func() {
		r := registry{
			ExtraParameters: 2,
		}
		r.Register("extra", func(a, b string) {
			t.FailNow()
		})
		r.Execute("extra", "a", "b", "c")
	})
}

func TestEmpty(t *testing.T) {
	r := registry{}
	err := r.Execute("test")
	if err != ErrCommandNotFound {
		t.FailNow()
	}
}

func TestMissing(t *testing.T) {
	r := registry{}
	r.Register("hello world", func() {
		t.FailNow()
	})
	err := r.Execute("test")
	if err != ErrCommandNotFound {
		t.FailNow()
	}
}

func TestMissing2(t *testing.T) {
	r := registry{}
	r.Register("hello world", func() {
		t.FailNow()
	})
	err := r.Execute("hello")
	if err != ErrCommandNotFound {
		t.FailNow()
	}
}

func TestQuoted(t *testing.T) {
	r := registry{}
	called := false
	r.Register("hello world", func() {
		called = true
	})
	checkError(t, r.Execute("hello \"world\""))
	if !called {
		t.FailNow()
	}
}

func TestCommandPanic(t *testing.T) {
	r := registry{}
	r.Register("hello world", func() {
		panic("test panic")
	})
	err := r.Execute("hello \"world\"")
	if err == nil || err.Error() != "test panic" {
		t.FailNow()
	}
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func shouldPanic(t *testing.T, f func()) {
	defer func() {
		if err := recover(); err == nil {
			t.FailNow()
		}
	}()
	f()
}
