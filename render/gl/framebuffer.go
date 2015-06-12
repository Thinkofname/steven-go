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

package gl

import (
	"github.com/thinkofdeath/gl/v3.2-core/gl"
)

type Attachment uint32

const (
	ColorAttachment0 Attachment = gl.COLOR_ATTACHMENT0
	ColorAttachment1 Attachment = gl.COLOR_ATTACHMENT1
	ColorAttachment2 Attachment = gl.COLOR_ATTACHMENT2
	DepthAttachment  Attachment = gl.DEPTH_ATTACHMENT
)

type Framebuffer struct {
	internal uint32
}

func NewFramebuffer() Framebuffer {
	f := Framebuffer{}
	gl.GenFramebuffers(1, &f.internal)
	return f
}

func BindFragDataLocation(p Program, cn int, name string) {
	n := gl.Str(name + "\x00")
	gl.BindFragDataLocation(uint32(p), uint32(cn), n)
}

func (f Framebuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, f.internal)
}

func UnbindFramebuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func DrawBuffers(bufs []Attachment) {
	gl.DrawBuffers(int32(len(bufs)), (*uint32)(gl.Ptr(bufs)))
}

func (f Framebuffer) Texture2D(attachment Attachment, texTarget TextureTarget, tex Texture, level int) {
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, uint32(attachment), uint32(texTarget), tex.internal, int32(level))
}

func (f Framebuffer) Delete() {
	gl.DeleteFramebuffers(1, &f.internal)
}
