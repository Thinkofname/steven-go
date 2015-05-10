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

package steven

import "github.com/thinkofdeath/steven/type/vmath"

func newCow() Entity {
	type cow struct {
		networkComponent
		positionComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	c := &cow{
		debugComponent: debugComponent{255, 0, 0},
	}
	c.networkID = 92
	c.bounds = vmath.NewAABB(-0.45, 0, -0.45, 0.9, 1.3, 0.9)
	return c
}

func newPlayer() Entity {
	type player struct {
		positionComponent
		targetPositionComponent
		sizeComponent

		debugComponent
	}
	p := &player{
		debugComponent: debugComponent{255, 0, 255},
	}
	p.bounds = vmath.NewAABB(-0.3, 0, -0.3, 0.6, 1.8, 0.6)
	return p
}
