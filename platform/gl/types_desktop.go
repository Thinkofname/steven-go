// +build !mobile

package gl

import (
	"github.com/go-gl/gl/v2.1/gl"
)

type Type uint32

const (
	UnsignedByte  Type = gl.UNSIGNED_BYTE
	UnsignedShort Type = gl.UNSIGNED_SHORT
	Short         Type = gl.SHORT
	Float         Type = gl.FLOAT
)
