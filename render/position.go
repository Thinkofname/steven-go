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

package render

type position struct {
	X, Y, Z int
}

type positionC struct {
	X, Z int
}

type transList []position

func (t transList) Swap(a, b int) {
	t[a], t[b] = t[b], t[a]
}

func (t transList) Less(aa, bb int) bool {
	a := t[aa]
	b := t[bb]
	dx := float64(a.X<<4) + 8 - Camera.X
	dy := float64(a.Y<<4) + 8 - Camera.Y
	dz := float64(a.Z<<4) + 8 - Camera.Z
	adist := dx*dx + dy*dy + dz*dz

	dx = float64(b.X<<4) + 8 - Camera.X
	dy = float64(b.Y<<4) + 8 - Camera.Y
	dz = float64(b.Z<<4) + 8 - Camera.Z
	bdist := dx*dx + dy*dy + dz*dz
	return adist > bdist
}

func (t transList) Len() int {
	return len(t)
}
