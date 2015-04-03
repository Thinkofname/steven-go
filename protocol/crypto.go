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

package protocol

import (
	"crypto/cipher"
)

type cfb8Stream struct {
	block    cipher.Block
	iv, orig []byte
	offset   int
	decrypt  bool
}

func newCFB8(block cipher.Block, key []byte, decrypt bool) *cfb8Stream {
	cp := make([]byte, 16)
	copy(cp, key)
	orig := make([]byte, 32)
	copy(orig, key)
	copy(orig[16:], key)
	return &cfb8Stream{
		block:   block,
		iv:      cp,
		orig:    orig,
		decrypt: decrypt,
	}
}

func (c *cfb8Stream) XORKeyStream(dst, src []byte) {
	for i, b := range src {
		c.block.Encrypt(c.iv, c.iv)
		b ^= c.iv[0]

		if c.decrypt {
			c.orig[c.offset] = src[i]
			c.orig[c.offset+16] = src[i]
		} else {
			c.orig[c.offset] = b
			c.orig[c.offset+16] = b
		}
		c.offset = (c.offset + 1) & 0xF
		copy(c.iv, c.orig[c.offset:])
		dst[i] = b
	}
}
