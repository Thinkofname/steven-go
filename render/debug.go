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
	"github.com/thinkofdeath/steven/render/gl"
)

type debugBuffer struct {
	fb *gl.Framebuffer
	a  gl.Attachment
}

var debugBuffers = []debugBuffer{
	{&transFramebuffer, gl.ColorAttachment0},
	{&transFramebuffer, gl.ColorAttachment1},
}

func blitBuffers() {
	// Screen
	gl.Framebuffer{}.BindDraw()

	const width, height = 300, 150

	ox, oy := 5, 5
	for _, d := range debugBuffers {
		d.fb.BindRead()
		d.fb.ReadBuffer(d.a)
		gl.BlitFramebuffer(
			0, 0, lastWidth, lastHeight,
			ox, lastHeight-height-oy, ox+width, lastHeight-oy,
			gl.ColorBufferBit, gl.Nearest,
		)
		ox += width
		if ox+width >= lastWidth {
			ox = 5
			oy += height
		}
	}

	gl.UnbindFramebufferDraw()
	gl.UnbindFramebufferRead()
}
