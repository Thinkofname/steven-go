package vmath

import (
	"math"
)

// Matrix4 represents a 4x4 matrix
type Matrix4 [4 * 4]float32

func NewMatrix4() *Matrix4 {
	m := new(Matrix4)
	m.Identity()
	return m
}

// Identity sets the matrix to the identity matrix
func (m *Matrix4) Identity() {
	for i := range m {
		m[i] = 0
	}
	m[k(0, 0)] = 1
	m[k(1, 1)] = 1
	m[k(2, 2)] = 1
	m[k(3, 3)] = 1
}

// Perspective applies a perspective to the matrix
func (m *Matrix4) Perspective(fovy, aspect, near, far float32) {
	invDepth := 1 / (near - far)

	m[k(1, 1)] = float32(1 / math.Tan(float64(0.5*fovy)))
	m[k(0, 0)] = m[k(1, 1)] / aspect
	m[k(2, 2)] = (far + near) * invDepth
	m[k(3, 2)] = 2 * far * near * invDepth
	m[k(2, 3)] = -1
	m[k(3, 3)] = 0
}

// Scale scales the matrix by the passed values
func (m *Matrix4) Scale(x, y, z float32) {
	for i := 0; i < 4; i++ {
		m[i] *= x
		m[i+4] *= y
		m[i+8] *= z
	}
}

// Translate translates the matrix by the passed values
func (m *Matrix4) Translate(x, y, z float32) {
	for i := 0; i < 4; i++ {
		o := i * 4
		m[o] = m[o] + m[o+3]*x
		m[o+1] = m[o+1] + m[o+3]*y
		m[o+2] = m[o+2] + m[o+3]*z
	}
}

// RotateX rotates the matrix by the passed value around
// the X axis
func (m *Matrix4) RotateX(ang float32) {
	c := float32(math.Cos(float64(ang)))
	s := float32(math.Sin(float64(ang)))

	for i := 0; i < 4; i++ {
		o := i * 4
		t := m[o+1]
		m[o+1] = t*c + m[o+2]*s
		m[o+2] = t*-s + m[o+2]*c
	}
}

// RotateY rotates the matrix by the passed value around
// the Y axis
func (m *Matrix4) RotateY(ang float32) {
	c := float32(math.Cos(float64(ang)))
	s := float32(math.Sin(float64(ang)))

	for i := 0; i < 4; i++ {
		o := i * 4
		t := m[o]
		m[o] = t*c + m[o+2]*s
		m[o+2] = t*-s + m[o+2]*c
	}
}

// RotateZ rotates the matrix by the passed value around
// the Z axis
func (m *Matrix4) RotateZ(ang float32) {
	c := float32(math.Cos(float64(ang)))
	s := float32(math.Sin(float64(ang)))

	for i := 0; i < 4; i++ {
		o := i * 4
		t := m[o]
		m[o] = t*c + m[o+1]*s
		m[o+1] = t*-s + m[o+1]*c
	}
}

// Multiply multiplies this matrix and the passed matrix
// together, storing the result in this matrix
func (m *Matrix4) Multiply(o *Matrix4) {
	for i := 0; i < 16; i += 4 {
		a, b, c, d := m[i], m[i+1], m[i+2], m[i+3]

		m[i] = a*o[0] + b*o[4] + c*o[8] + d*o[12]
		m[i+1] = a*o[1] + b*o[5] + c*o[9] + d*o[13]
		m[i+2] = a*o[2] + b*o[6] + c*o[10] + d*o[14]
		m[i+3] = a*o[3] + b*o[7] + c*o[11] + d*o[15]
	}
}

// Helper method for indexing the matrix
func k(x, y int) int {
	return y + x*4
}
