package gl

import (
	"github.com/go-gl/gl/v3.2-core/gl"
)

type Type uint32

const (
	UnsignedByte  Type = gl.UNSIGNED_BYTE
	UnsignedShort Type = gl.UNSIGNED_SHORT
	Short         Type = gl.SHORT
	Float         Type = gl.FLOAT
)
