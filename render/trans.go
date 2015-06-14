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
	"github.com/thinkofdeath/steven/console"
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
	transDepth       gl.Texture
	transCreated     bool
	transState       = struct {
		program gl.Program
		shader  *transShader
		array   gl.VertexArray
		buffer  gl.Buffer
	}{}

	rSamples = console.NewIntVar("r_samples", 1, console.Serializable, console.Mutable).Doc(`
r_samples controls the number of samples taken whilst rendering. 
Increasing this will increase the amount of AA used but decrease
performance.
`).Callback(func() { lastWidth = -1; lastHeight = -1 })
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

	transDepth = gl.CreateTexture()
	transDepth.Bind(gl.Texture2D)
	transDepth.Image2DEx(0, lastWidth, lastHeight, gl.DepthComponent24, gl.DepthComponent, gl.UnsignedByte, nil)
	transDepth.Parameter(gl.TextureMinFilter, gl.Linear)
	transDepth.Parameter(gl.TextureMagFilter, gl.Linear)
	transFramebuffer.Texture2D(gl.DepthAttachment, gl.Texture2D, transDepth, 0)

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
	fbColor.Bind(gl.Texture2DMultisample)
	fbColor.Image2DSample(rSamples.Value(), lastWidth, lastHeight, gl.RGBA8, false)
	mainFramebuffer.Texture2D(gl.ColorAttachment0, gl.Texture2DMultisample, fbColor, 0)

	fbDepth = gl.CreateTexture()
	fbDepth.Bind(gl.Texture2DMultisample)
	fbDepth.Image2DSample(rSamples.Value(), lastWidth, lastHeight, gl.DepthComponent24, false)
	mainFramebuffer.Texture2D(gl.DepthAttachment, gl.Texture2DMultisample, fbDepth, 0)

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
	fbColor.Bind(gl.Texture2DMultisample)

	transState.program.Use()
	transState.shader.Accum.Int(0)
	transState.shader.Revealage.Int(1)
	transState.shader.Color.Int(2)
	transState.shader.Samples.Int(rSamples.Value())
	transState.array.Bind()
	gl.DrawArrays(gl.Triangles, 0, 6)
}

type transShader struct {
	Position  gl.Attribute `gl:"aPosition"`
	Accum     gl.Uniform   `gl:"taccum"`
	Revealage gl.Uniform   `gl:"trevealage"`
	Color     gl.Uniform   `gl:"tcolor"`
	Samples   gl.Uniform   `gl:"samples"`
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
uniform sampler2DMS tcolor;

uniform int samples;

out vec4 fragColor;

void main() {	
	ivec2 C = ivec2(gl_FragCoord.xy);
	vec4 accum = texelFetch(taccum, C, 0);
	float aa = texelFetch(trevealage, C, 0).r;
	vec4 col = texelFetch(tcolor, C, 0);

	for (int i = 1; i < samples; i++) {
		col += texelFetch(tcolor, C, i);
	}
	col /= float(samples);

	float r = accum.a;
	accum.a = aa;
	if (r >= 1.0) {
		fragColor = vec4(col.rgb, 0.0);
	} else {
		vec3 alp = clamp(accum.rgb / clamp(accum.a, 1e-4, 5e4), 0.0, 1.0);
		fragColor = vec4(col.rgb * r  + alp * (1.0 - r), 0.0);
	}
}
`
)
