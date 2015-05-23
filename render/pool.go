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

var (
	freeTextures []*reusableTexture
)

type reusableTexture struct {
	info *textureInfo
	W, H int
}

func getFreeTexture(w, h int) *reusableTexture {
	for i, f := range freeTextures {
		if f.H == h && f.W == w {
			freeTextures = append(freeTextures[:i], freeTextures[i+1:]...)
			return f
		}
	}
	return &reusableTexture{
		info: addTexture(make([]byte, w*h*4), w, h),
		W:    w, H: h,
	}
}

func freeTexture(t *reusableTexture) {
	freeTextures = append(freeTextures, t)
}
