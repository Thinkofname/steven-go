package vmath

type Vector3 struct {
	X, Y, Z float32
}

func (v Vector3) Dot(other Vector3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vector3) Equals(other Vector3) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}
