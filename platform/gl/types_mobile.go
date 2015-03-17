// +build mobile

package gl

import (
	"golang.org/x/mobile/gl"
)

type Type uint32

const (
	UnsignedByte  Type = gl.UNSIGNED_BYTE
	UnsignedShort Type = gl.UNSIGNED_SHORT
	Short         Type = gl.SHORT
	Float         Type = gl.FLOAT
)
