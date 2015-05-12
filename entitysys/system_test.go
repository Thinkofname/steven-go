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

import "testing"

type nameable struct {
	name string
}

func (n *nameable) Name() string        { return n.name }
func (n *nameable) SetName(name string) { n.name = name }

type Nameable interface {
	Name() string
	SetName(string)
}

type countable struct {
	counter int
}

func (c *countable) IncCounter()  { c.counter++ }
func (c *countable) Counter() int { return c.counter }

type Countable interface {
	IncCounter()
	Counter() int
}

func TestBasic(t *testing.T) {
	c := NewContainer()
	tSys := func(n Nameable) {
		n.SetName("test_" + n.Name())
	}
	c.AddSystem(Tick, tSys)

	type testEntity struct {
		nameable
	}

	t1 := &testEntity{}
	t1.name = "bob"
	t2 := &testEntity{}
	t2.name = "steven"

	c.AddEntity(t1)
	c.AddEntity(t2)

	c.Tick()

	if t1.name != "test_bob" {
		t.Log(t1.name)
		t.FailNow()
	}
	if t2.name != "test_steven" {
		t.Log(t2.name)
		t.FailNow()
	}
}
func TestDifferent(t *testing.T) {
	c := NewContainer()
	tSys := func(n Nameable) {
		n.SetName("test_" + n.Name())
	}
	c.AddSystem(Tick, tSys)

	type testEntity struct {
		nameable
	}
	type countEntity struct {
		nameable
		countable
	}

	t1 := &testEntity{}
	t1.name = "bob"
	t2 := &countEntity{}
	t2.name = "county"
	t3 := &testEntity{}
	t3.name = "steven"

	c.AddEntity(t1)
	c.AddEntity(t2)
	c.AddEntity(t3)

	c.Tick()

	if t1.name != "test_bob" {
		t.Log(t1.name)
		t.FailNow()
	}
	if t2.name != "test_county" {
		t.Log(t2.name)
		t.FailNow()
	}
	if t3.name != "test_steven" {
		t.Log(t3.name)
		t.FailNow()
	}
}

func TestNot(t *testing.T) {
	c := NewContainer()
	tSys := func(n Nameable) {
		n.SetName("test_" + n.Name())
	}
	c.AddSystem(Tick, tSys, Not(Type((*Countable)(nil))))

	type testEntity struct {
		nameable
	}
	type countEntity struct {
		nameable
		countable
	}

	t1 := &testEntity{}
	t1.name = "bob"
	t2 := &countEntity{}
	t2.name = "county"
	t3 := &testEntity{}
	t3.name = "steven"

	c.AddEntity(t1)
	c.AddEntity(t2)
	c.AddEntity(t3)

	c.Tick()

	if t1.name != "test_bob" {
		t.Log(t1.name)
		t.FailNow()
	}
	if t2.name != "county" {
		t.Log(t2.name)
		t.FailNow()
	}
	if t3.name != "test_steven" {
		t.Log(t3.name)
		t.FailNow()
	}
}

func TestMultiple(t *testing.T) {
	c := NewContainer()
	tSys := func(n Nameable) {
		n.SetName("test_" + n.Name())
	}
	c.AddSystem(Tick, tSys)
	tSys2 := func(c Countable) {
		c.IncCounter()
	}
	c.AddSystem(Tick, tSys2)

	type testEntity struct {
		nameable
	}
	type countEntity struct {
		nameable
		countable
	}

	t1 := &testEntity{}
	t1.name = "bob"
	t2 := &countEntity{}
	t2.name = "county"
	t3 := &testEntity{}
	t3.name = "steven"

	c.AddEntity(t1)
	c.AddEntity(t2)
	c.AddEntity(t3)

	c.Tick()

	if t1.name != "test_bob" {
		t.Log(t1.name)
		t.FailNow()
	}
	if t2.name != "test_county" {
		t.Log(t2.name)
		t.FailNow()
	}
	if t2.counter != 1 {
		t.Log(t2.counter)
		t.FailNow()
	}
	if t3.name != "test_steven" {
		t.Log(t3.name)
		t.FailNow()
	}
}

func TestDuel(t *testing.T) {
	c := NewContainer()
	tSys := func(n Nameable, c Countable) {
		n.SetName("test_" + n.Name())
		c.IncCounter()
	}
	c.AddSystem(Tick, tSys)

	type testEntity struct {
		nameable
	}
	type countEntity struct {
		nameable
		countable
	}

	t1 := &testEntity{}
	t1.name = "bob"
	t2 := &countEntity{}
	t2.name = "county"
	t3 := &testEntity{}
	t3.name = "steven"

	c.AddEntity(t1)
	c.AddEntity(t2)
	c.AddEntity(t3)

	c.Tick()

	if t1.name != "bob" {
		t.Log(t1.name)
		t.FailNow()
	}
	if t2.name != "test_county" {
		t.Log(t2.name)
		t.FailNow()
	}
	if t2.counter != 1 {
		t.Log(t2.counter)
		t.FailNow()
	}
	if t3.name != "steven" {
		t.Log(t3.name)
		t.FailNow()
	}
}

func TestAdd(t *testing.T) {
	c := NewContainer()
	tSys := func(n Nameable) {
		n.SetName("add_" + n.Name())
	}
	c.AddSystem(Add, tSys)

	type testEntity struct {
		nameable
	}
	t1 := &testEntity{}
	t1.name = "bob"
	t2 := &testEntity{}
	t2.name = "steven"

	c.AddEntity(t1)
	c.AddEntity(t2)

	if t1.name != "add_bob" {
		t.Log(t1.name)
		t.FailNow()
	}
	if t2.name != "add_steven" {
		t.Log(t2.name)
		t.FailNow()
	}
}

func TestRemove(t *testing.T) {
	c := NewContainer()
	tSys := func(n Nameable) {
		n.SetName("remove_" + n.Name())
	}
	c.AddSystem(Remove, tSys)

	type testEntity struct {
		nameable
	}
	t1 := &testEntity{}
	t1.name = "bob"
	t2 := &testEntity{}
	t2.name = "steven"

	c.AddEntity(t1)
	c.AddEntity(t2)

	if t1.name != "bob" {
		t.Log(t1.name)
		t.FailNow()
	}
	if t2.name != "steven" {
		t.Log(t2.name)
		t.FailNow()
	}

	c.RemoveEntity(t1)

	if t1.name != "remove_bob" {
		t.Log(t1.name)
		t.FailNow()
	}

	c.RemoveEntity(t2)

	if t2.name != "remove_steven" {
		t.Log(t2.name)
		t.FailNow()
	}
}

func TestRaw(t *testing.T) {
	c := NewContainer()
	tSys := func(n *nameable) {
		n.SetName("remove_" + n.Name())
	}
	c.AddSystem(Remove, tSys)

	type testEntity struct {
		nameable
	}
	t1 := &testEntity{}
	t1.name = "bob"
	t2 := &testEntity{}
	t2.name = "steven"

	c.AddEntity(t1)
	c.AddEntity(t2)

	if t1.name != "bob" {
		t.Log(t1.name)
		t.FailNow()
	}
	if t2.name != "steven" {
		t.Log(t2.name)
		t.FailNow()
	}

	c.RemoveEntity(t1)

	if t1.name != "remove_bob" {
		t.Log(t1.name)
		t.FailNow()
	}

	c.RemoveEntity(t2)

	if t2.name != "remove_steven" {
		t.Log(t2.name)
		t.FailNow()
	}
}
