package vmath

type Frustum struct {
	planes [6]fPlane
}

type fPlane struct {
	N Vector3
	D float32
}

func (f *Frustum) FromMatrix(m *Matrix4) {
	for i := range f.planes {
		off := i >> 1
		f.planes[i] = fPlane{
			N: Vector3{
				X: m.Get(0, 3) - m.Get(0, off),
				Y: m.Get(1, 3) - m.Get(1, off),
				Z: m.Get(2, 3) - m.Get(2, off),
			},
			D: m.Get(3, 3) - m.Get(3, off),
		}
	}

	for i := range f.planes {
		f.planes[i].N.Normalize()
	}
}

func (f *Frustum) IsSphereInside(x, y, z, radius float32) bool {
	center := Vector3{x, y, z}
	for i := 0; i < 6; i++ {
		if center.Dot(f.planes[i].N)+f.planes[i].D+radius <= 0 {
			return false
		}
	}
	return true
}
