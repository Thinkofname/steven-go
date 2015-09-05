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

import "testing"

func TestMap(t *testing.T) {
	m := NewMap(4096, 4)
	for i := 0; i < 4096; i++ {
		for j := 0; j < 16; j++ {
			m.Set(i, j)
			if m.Get(i) != j {
				t.Fatalf("Index(%d) wanted %d and got %d", i, j, m.Get(i))
			}
		}
	}
}

func TestMapOdd(t *testing.T) {
	for size := 1; size <= 16; size++ {
		m := NewMap(64*3, size)
		max := (1 << uint(size)) - 1
		for i := 0; i < 64*3; i++ {
			for j := 0; j < max; j++ {
				m.Set(i, j)
				if m.Get(i) != j {
					t.Fatalf("Index(%d) wanted %d and got %d", i, j, m.Get(i))
				}
			}
		}
	}
}
