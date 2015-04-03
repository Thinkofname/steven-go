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

import "math"

// Camera is the main camera for the renderer
var Camera = &ClientCamera{}

func init() {
	Camera.Pitch = math.Pi
}

// ClientCamera is a camera within the world
type ClientCamera struct {
	X, Y, Z    float64
	Yaw, Pitch float64
}
