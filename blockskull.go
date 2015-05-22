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

package steven

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/encoding/nbt"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/type/direction"
)

type skullType int

const (
	skullSkeleton skullType = iota
	skullWitherSkeleton
	skullZombie
	skullPlayer
	skullCreeper
)

type skullComponent struct {
	SkullType skullType
	Rotation  int
	Facing    direction.Type
	Owner     string
	OwnerSkin render.TextureInfo
	model     *render.StaticModel
	position  Position
}

func (s *skullComponent) Model() *render.StaticModel {
	return s.model
}

func (s *skullComponent) CanHandleAction(action int) bool {
	return action == 4
}

func (s *skullComponent) Deserilize(tag *nbt.Compound) {
	t, ok := tag.Items["SkullType"].(int8)
	if !ok {
		return
	}
	s.SkullType = skullType(t)
	rot, ok := tag.Items["Rot"].(int8)
	if !ok {
		return
	}
	s.Rotation = int(rot)

	if s.SkullType != skullPlayer {
		s.free()
		s.create()
		return
	}

	owner, ok := tag.Items["Owner"].(*nbt.Compound)
	if !ok {
		return
	}
	props, ok := owner.Items["Properties"].(*nbt.Compound)
	if !ok {
		return
	}
	tex, ok := props.Items["textures"].(*nbt.List)
	if !ok || tex.Type != nbt.TagCompound || len(tex.Elements) < 1 {
		return
	}
	texP := tex.Elements[0].(*nbt.Compound)
	value := texP.Items["Value"].(string)
	sigV, hasSig := texP.Items["Signature"].(string)

	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return
	}

	if hasSig {
		sig, err := base64.StdEncoding.DecodeString(sigV)
		if err != nil {
			return
		}

		if err := verifySkinSignature([]byte(value), sig); err != nil {
			return
		}
	}

	var blob skinBlob
	err = json.Unmarshal(data, &blob)
	if err != nil {
		return
	}
	url := blob.Textures.Skin.Url
	// We can only handle textures from textures.minecraft.net currently,
	// luckily these are the only ones we really see in practice and
	// mojang seemed to have blocked other urls
	if strings.HasPrefix(url, "http://textures.minecraft.net/texture/") {
		s.free()
		s.Owner = url[len("http://textures.minecraft.net/texture/"):]
		render.RefSkin(s.Owner)
		s.OwnerSkin = render.Skin(s.Owner)
		s.create()
	}
}

func (s *skullComponent) free() {
	if s.model != nil {
		s.model.Free()
		s.model = nil
	}
	if s.Owner != "" && s.OwnerSkin != nil {
		s.OwnerSkin = nil
		render.FreeSkin(s.Owner)
	}
}

func (s *skullComponent) create() {
	var skin render.TextureInfo
	if s.SkullType == skullPlayer && s.OwnerSkin != nil {
		skin = s.OwnerSkin
	} else {
		switch s.SkullType {
		case skullPlayer:
			skin = render.GetTexture("entity/steve")
		case skullZombie:
			skin = render.GetTexture("entity/zombie/zombie")
		case skullSkeleton:
			skin = render.GetTexture("entity/skeleton/skeleton")
		case skullWitherSkeleton:
			skin = render.GetTexture("entity/skeleton/wither_skeleton")
		case skullCreeper:
			skin = render.GetTexture("entity/creeper/creeper")
		}
	}

	var hverts []*render.StaticVertex
	// Base layer
	hverts = appendBox(hverts, -4/16.0, 0, -4/16.0, 8/16.0, 8/16.0, 8/16.0, [6]render.TextureInfo{
		direction.North: skin.Sub(8, 8, 8, 8),
		direction.South: skin.Sub(24, 8, 8, 8),
		direction.West:  skin.Sub(0, 8, 8, 8),
		direction.East:  skin.Sub(16, 8, 8, 8),
		direction.Up:    skin.Sub(8, 0, 8, 8),
		direction.Down:  skin.Sub(16, 0, 8, 8),
	})
	// Hat layer
	hverts = appendBox(hverts, -4.2/16.0, -.2/16.0, -4.2/16.0, 8.4/16.0, 8.4/16.0, 8.4/16.0, [6]render.TextureInfo{
		direction.North: skin.Sub(8+32, 8, 8, 8),
		direction.South: skin.Sub(24+32, 8, 8, 8),
		direction.West:  skin.Sub(0+32, 8, 8, 8),
		direction.East:  skin.Sub(16+32, 8, 8, 8),
		direction.Up:    skin.Sub(8+32, 0, 8, 8),
		direction.Down:  skin.Sub(16+32, 0, 8, 8),
	})

	s.model = render.NewStaticModel([][]*render.StaticVertex{
		hverts,
	})
	model := s.model
	model.Radius = 2

	x, y, z := s.position.X, s.position.Y, s.position.Z

	model.X, model.Y, model.Z = -float32(x)-0.5, -float32(y), float32(z)+0.5

	mat := mgl32.Translate3D(float32(x)+0.5, -float32(y), float32(z)+0.5)
	if s.Facing == direction.Up {
		mat = mat.Mul4(mgl32.Rotate3DY(-(math.Pi / 8) * float32(s.Rotation)).Mat4())
	} else {
		ang := float32(0)
		switch s.Facing {
		case direction.South:
			ang = math.Pi
		case direction.East:
			ang = math.Pi / 2
		case direction.West:
			ang = -math.Pi / 2
		}
		mat = mat.Mul4(mgl32.Rotate3DY(ang).Mat4())
		mat = mat.Mul4(mgl32.Translate3D(0, -4/16.0, 4/16.0))
	}
	model.Matrix[0] = mat
	lightBlockModel(model, s.position)
}
func esSkullAdd(s *skullComponent, p BlockComponent) {
	s.position = p.Position()
	s.create()
}

func esSkullRemove(s *skullComponent) {
	s.free()
}
