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

package bit

type Map struct {
	bits []uint64
	size int
	mask int
}

func NewMap(l, size int) *Map {
	return &Map{
		size: size,
		bits: make([]uint64, (l*size)/64),
	}
}

func NewMapFromRaw(bits []uint64, size int) *Map {
	mask := 1
	for i := 1; i < size; i++ {
		mask <<= 1
	}
	return &Map{
		size: size,
		bits: bits,
	}
}

func (m *Map) Set(i, val int) {
	i *= m.size
	pos := i / 64
	mask := (uint64(1) << uint(m.size)) - 1
	i %= 64
	m.bits[pos] = (m.bits[pos] & ^(mask << uint64(i))) | (uint64(val) << uint64(i))
}

func (m *Map) Get(i int) int {
	i *= m.size
	pos := i / 64
	mask := (uint64(1) << uint(m.size)) - 1
	i %= 64
	return int((m.bits[pos] >> uint64(i)) & mask)
}
