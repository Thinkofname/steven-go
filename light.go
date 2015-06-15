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

import (
	"math"

	"github.com/thinkofdeath/steven/render"
)

func getLightColor(block, sky float64) (r, g, b float64) {
	bl := math.Pow(float64(render.LightLevel), 15.0-block)*15.0 + 1.0
	sk := math.Pow(float64(render.LightLevel), 15.0-sky)*15.0 + 1.0

	br := (0.879552*bl*bl + 0.871148*bl + 32.9821)
	bg := (1.22181*bl*bl - 4.78113*bl + 36.7125)
	bb := (1.67612*bl*bl - 12.9764*bl + 48.8321)

	sr := (0.131653*sk*sk - 0.761625*sk + 35.0393)
	sg := (0.136555*sk*sk - 0.853782*sk + 29.6143)
	sb := (0.327311*sk*sk - 1.62017*sk + 28.0929)
	srl := (0.996148*sk*sk - 4.19629*sk + 51.4036)
	sgl := (1.03904*sk*sk - 4.81516*sk + 47.0911)
	sbl := (1.076164*sk*sk - 5.36376*sk + 43.9089)

	skyOffset := float64(render.SkyOffset)

	sr = srl*skyOffset + sr*(1.0-skyOffset)
	sg = sgl*skyOffset + sg*(1.0-skyOffset)
	sb = sbl*skyOffset + sb*(1.0-skyOffset)

	return math.Sqrt((br*br+sr*sr)/2) / 255.0,
		math.Sqrt((bg*bg+sg*sg)/2) / 255.0,
		math.Sqrt((bb*bb+sb*sb)/2) / 255.0
}
