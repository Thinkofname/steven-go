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

import "github.com/thinkofdeath/phteven/render/gl"

type lineShader struct {
	Position          gl.Attribute `gl:"aPosition"`
	Color             gl.Attribute `gl:"aColor"`
	PerspectiveMatrix gl.Uniform   `gl:"perspectiveMatrix"`
	CameraMatrix      gl.Uniform   `gl:"cameraMatrix"`
}

const (
	vertexLine = `
#version 150
in vec3 aPosition;
in vec3 aColor;

uniform mat4 perspectiveMatrix;
uniform mat4 cameraMatrix;

out vec3 vColor;

void main() {
	vec3 pos = vec3(aPosition.x, -aPosition.y, aPosition.z);
	gl_Position = perspectiveMatrix * cameraMatrix * vec4(pos, 1.0);
	vColor = aColor;
}
`
	fragmentLine = `
#version 150

const float atlasSize = ` + atlasSizeStr + `;

in vec3 vColor;

out vec4 fragColor;

void main() {
	fragColor = vec4(vColor, 1.0);
}
`
)
