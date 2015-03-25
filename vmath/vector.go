package vmath

// Vector3 is a 3 component vector
type Vector3 struct {
	X, Y, Z float32
}

// Dot returns the result of preforming the dot operation on this
// vector and the passed vector.
func (v Vector3) Dot(other Vector3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}
