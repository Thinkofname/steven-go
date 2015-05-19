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

package vmath

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Frustum struct {
	planes                  [6]fPlane
	fovy, aspect, near, far float32
	tang, nh, nw, fh, fw    float32
}

type fPlane struct {
	N, P mgl32.Vec3
	D    float32
}

func (f *fPlane) setPoints(v1, v2, v3 mgl32.Vec3) {
	aux1 := v1.Sub(v2)
	aux2 := v3.Sub(v2)

	f.N = aux2.Cross(aux1)
	f.N = safeNormalize(f.N)
	f.P = v2
	f.D = -(f.N.Dot(f.P))
}

func NewFrustum() *Frustum {
	return &Frustum{}
}

func safeNormalize(v mgl32.Vec3) mgl32.Vec3 {
	v = v.Normalize()
	if math.IsInf(float64(v[0]), 0) || math.IsNaN(float64(v[0])) {
		return mgl32.Vec3{}
	}
	return v
}

func (f *Frustum) SetPerspective(fovy, aspect, near, far float32) {
	f.fovy = fovy
	f.aspect = aspect
	f.near = near
	f.far = far

	f.tang = float32(math.Tan(float64(fovy * 0.5)))
	f.nh = near * f.tang
	f.nw = f.nh * aspect
	f.fh = far * f.tang
	f.fw = f.fh * aspect
}

func (f *Frustum) SetCamera(p, l, u mgl32.Vec3) {
	Z := p.Sub(l)
	Z = safeNormalize(Z)

	X := u.Cross(Z)
	X = safeNormalize(X)

	Y := Z.Cross(X)

	nc := p.Sub(Z.Mul(f.near))
	fc := p.Sub(Z.Mul(f.far))

	ntl := nc.Add(Y.Mul(f.nh)).Sub(X.Mul(f.nw))
	ntr := nc.Add(Y.Mul(f.nh)).Add(X.Mul(f.nw))
	nbl := nc.Sub(Y.Mul(f.nh)).Sub(X.Mul(f.nw))
	nbr := nc.Sub(Y.Mul(f.nh)).Add(X.Mul(f.nw))

	ftl := fc.Add(Y.Mul(f.fh)).Sub(X.Mul(f.fw))
	ftr := fc.Add(Y.Mul(f.fh)).Add(X.Mul(f.fw))
	fbl := fc.Sub(Y.Mul(f.fh)).Sub(X.Mul(f.fw))
	fbr := fc.Sub(Y.Mul(f.fh)).Add(X.Mul(f.fw))

	const (
		top = iota
		bottom
		left
		right
		nearP
		farP
	)
	f.planes[top].setPoints(ntr, ntl, ftl)
	f.planes[bottom].setPoints(nbl, nbr, fbr)
	f.planes[left].setPoints(ntl, nbl, fbl)
	f.planes[right].setPoints(nbr, ntr, fbr)
	f.planes[nearP].setPoints(ntl, ntr, nbr)
	f.planes[farP].setPoints(ftr, ftl, fbl)
}

func (f *Frustum) IsSphereInside(x, y, z, radius float32) bool {
	p := mgl32.Vec3{x, y, z}
	for i := range f.planes {
		dist := f.planes[i].D + f.planes[i].N.Dot(p)
		if dist < -radius {
			return false
		}
	}
	return true
}

func (f *Frustum) IsAABBInside(aabb AABB) bool {
	for i := range f.planes {
		v := aabb.Min
		for j := 0; j < 3; j++ {
			if f.planes[i].N[j] >= 0 {
				v[j] = aabb.Max[j]
			}
		}
		if f.planes[i].N.Dot(v)+f.planes[i].D < 0 {
			return false
		}
	}
	return true
}
