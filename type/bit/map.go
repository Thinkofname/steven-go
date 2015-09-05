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

import "fmt"

type Map struct {
	bits    []uint64
	BitSize int
	Length  int
}

func NewMap(l, size int) *Map {
	return &Map{
		BitSize: size,
		bits:    make([]uint64, (l*size)/64),
		Length:  l,
	}
}

func NewMapFromRaw(bits []uint64, size int) *Map {
	return &Map{
		BitSize: size,
		bits:    bits,
		Length:  (len(bits)*64 + 63) / size,
	}
}

func (m *Map) ResizeBits(size int) *Map {
	n := NewMap(m.Length, size)
	for i := 0; i < m.Length; i++ {
		n.Set(i, m.Get(i))
	}
	return n
}

func (m *Map) Set(i, val int) {
	if val < 0 || val >= (int(1)<<uint(m.BitSize)) {
		panic(fmt.Sprintf("invalid value %d %d", val, m.BitSize))
	}
	i *= m.BitSize
	pos := i / 64
	mask := (uint64(1) << uint(m.BitSize)) - 1
	ii := i % 64
	m.bits[pos] = (m.bits[pos] & ^(mask << uint64(ii))) | (uint64(val) << uint64(ii))
	pos2 := (i + m.BitSize - 1) / 64
	if pos2 != pos {
		used := 64 - ii
		rem := m.BitSize - used
		m.bits[pos2] = (m.bits[pos2] >> uint(rem) << uint(rem)) | (uint64(val) >> uint64(used))
	}
}

func (m *Map) Get(i int) int {
	i *= m.BitSize
	pos := i / 64
	mask := (uint64(1) << uint(m.BitSize)) - 1
	ii := i % 64
	pos2 := (i + m.BitSize - 1) / 64
	if pos2 != pos {
		used := 64 - ii
		return int(((m.bits[pos] >> uint64(ii)) | (m.bits[pos2])<<uint(used)) & mask)
	} else {
		return int((m.bits[pos] >> uint64(ii)) & mask)
	}
}
