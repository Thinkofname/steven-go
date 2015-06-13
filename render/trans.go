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

import (
	"github.com/thinkofdeath/steven/render/builder"
	"github.com/thinkofdeath/steven/render/gl"
)

var (
	mainFramebuffer  gl.Framebuffer
	fbColor          gl.Texture
	fbDepth          gl.Texture
	transFramebuffer gl.Framebuffer
	accum            gl.Texture
	revealage        gl.Texture
	transCreated     bool
	transState       = struct {
		program gl.Program
		shader  *transShader
		array   gl.VertexArray
		buffer  gl.Buffer
	}{}
)

func initTrans() {
	if transCreated {
		accum.Delete()
		revealage.Delete()
		transFramebuffer.Delete()
		fbColor.Delete()
		fbDepth.Delete()
		mainFramebuffer.Delete()
	}
	transCreated = true
	transFramebuffer = gl.NewFramebuffer()
	transFramebuffer.Bind()

	accum = gl.CreateTexture()
	accum.Bind(gl.Texture2D)
	accum.Image2DEx(0, lastWidth, lastHeight, gl.RGBA16F, gl.RGBA, gl.Float, nil)
	accum.Parameter(gl.TextureMinFilter, gl.Linear)
	accum.Parameter(gl.TextureMagFilter, gl.Linear)
	transFramebuffer.Texture2D(gl.ColorAttachment0, gl.Texture2D, accum, 0)

	revealage = gl.CreateTexture()
	revealage.Bind(gl.Texture2D)
	revealage.Image2DEx(0, lastWidth, lastHeight, gl.R16F, gl.Red, gl.Float, nil)
	revealage.Parameter(gl.TextureMinFilter, gl.Linear)
	revealage.Parameter(gl.TextureMagFilter, gl.Linear)
	transFramebuffer.Texture2D(gl.ColorAttachment1, gl.Texture2D, revealage, 0)

	fbDepth = gl.CreateTexture()
	fbDepth.Bind(gl.Texture2D)
	fbDepth.Image2DEx(0, lastWidth, lastHeight, gl.DepthComponent24, gl.DepthComponent, gl.UnsignedByte, nil)
	fbDepth.Parameter(gl.TextureMinFilter, gl.Linear)
	fbDepth.Parameter(gl.TextureMagFilter, gl.Linear)
	transFramebuffer.Texture2D(gl.DepthAttachment, gl.Texture2D, fbDepth, 0)

	chunkProgramT.Use()
	gl.BindFragDataLocation(chunkProgramT, 0, "accum")
	gl.BindFragDataLocation(chunkProgramT, 1, "revealage")

	gl.DrawBuffers([]gl.Attachment{
		gl.ColorAttachment0,
		gl.ColorAttachment1,
	})

	mainFramebuffer = gl.NewFramebuffer()
	mainFramebuffer.Bind()

	fbColor = gl.CreateTexture()
	fbColor.Bind(gl.Texture2D)
	fbColor.Image2DEx(0, lastWidth, lastHeight, gl.RGBA8, gl.RGBA, gl.UnsignedByte, nil)
	fbColor.Parameter(gl.TextureMinFilter, gl.Linear)
	fbColor.Parameter(gl.TextureMagFilter, gl.Linear)
	mainFramebuffer.Texture2D(gl.ColorAttachment0, gl.Texture2D, fbColor, 0)

	fbDepth.Bind(gl.Texture2D)
	mainFramebuffer.Texture2D(gl.DepthAttachment, gl.Texture2D, fbDepth, 0)

	gl.UnbindFramebuffer()

	transState.program = CreateProgram(vertexTrans, fragmentTrans)
	transState.shader = &transShader{}
	InitStruct(transState.shader, transState.program)

	transState.array = gl.CreateVertexArray()
	transState.array.Bind()
	transState.buffer = gl.CreateBuffer()
	transState.buffer.Bind(gl.ArrayBuffer)

	data := builder.New()
	for _, f := range []float32{-1, 1, 1, -1, -1, -1, 1, 1, 1, -1, -1, 1} {
		data.Float(f)
	}

	transState.buffer.Data(data.Data(), gl.StaticDraw)
	transState.shader.Position.Enable()
	transState.shader.Position.Pointer(2, gl.Float, false, 8, 0)
}

func transDraw() {
	gl.ActiveTexture(0)
	accum.Bind(gl.Texture2D)
	gl.ActiveTexture(1)
	revealage.Bind(gl.Texture2D)
	gl.ActiveTexture(2)
	fbColor.Bind(gl.Texture2D)

	transState.program.Use()
	transState.shader.Accum.Int(0)
	transState.shader.Revealage.Int(1)
	transState.shader.Color.Int(2)
	transState.array.Bind()
	gl.DrawArrays(gl.Triangles, 0, 6)
}

type transShader struct {
	Position  gl.Attribute `gl:"aPosition"`
	Accum     gl.Uniform   `gl:"taccum"`
	Revealage gl.Uniform   `gl:"trevealage"`
	Color     gl.Uniform   `gl:"tcolor"`
}

const (
	vertexTrans = `
#version 150
in vec2 aPosition;

void main() {
    gl_Position = vec4(aPosition,0,1);
}
`
	fragmentTrans = `
#version 150

uniform sampler2D taccum;
uniform sampler2D trevealage;
uniform sampler2D tcolor;

out vec4 fragColor;

void main() {	
	ivec2 C = ivec2(gl_FragCoord.xy);
    vec4 accum = texelFetch(taccum, C, 0);
    float r = accum.a;
    accum.a = texelFetch(trevealage, C, 0).r;
	vec4 col = texelFetch(tcolor, C, 0);
    if (r >= 1.0) {
		fragColor = vec4(col.rgb, 0.0);
	} else {
		vec3 alp = clamp(accum.rgb / clamp(accum.a, 1e-4, 5e4), 0.0, 1.0);
		fragColor = vec4(col.rgb * r  + alp * (1.0 - r), 0.0);
	}
}
`
)
