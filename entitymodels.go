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
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/thinkofdeath/steven/entitysys"
	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/type/direction"
)

func init() {
	addSystem(entitysys.Add, esPlayerModelAdd)
	addSystem(entitysys.Tick, esPlayerModelTick)
	addSystem(entitysys.Remove, esPlayerModelRemove)

	// Generic removal
	addSystem(entitysys.Remove, esModelRemove)
}

func appendBox(verts []*render.StaticVertex, x, y, z, w, h, d float32, textures [6]render.TextureInfo) []*render.StaticVertex {
	return appendBoxExtra(verts, x, y, z, w, h, d, textures, [6][2]float64{
		{1.0, 1.0},
		{1.0, 1.0},
		{1.0, 1.0},
		{1.0, 1.0},
		{1.0, 1.0},
		{1.0, 1.0},
	})
}
func appendBoxExtra(verts []*render.StaticVertex, x, y, z, w, h, d float32, textures [6]render.TextureInfo, extra [6][2]float64) []*render.StaticVertex {
	for i, face := range faceVertices {
		tex := textures[i]
		if tex == nil {
			continue
		}
		for _, v := range face.verts {
			var rr, gg, bb byte = 255, 255, 255
			if direction.Type(i) == direction.West || direction.Type(i) == direction.East {
				rr = byte(255 * 0.8)
				gg = byte(255 * 0.8)
				bb = byte(255 * 0.8)
			}
			vert := &render.StaticVertex{
				X:        float32(v.X)*w + x,
				Y:        float32(v.Y)*h + y,
				Z:        float32(v.Z)*d + z,
				Texture:  tex,
				TextureX: float64(v.TOffsetX) * extra[i][0],
				TextureY: float64(v.TOffsetY) * extra[i][1],
				R:        rr,
				G:        gg,
				B:        bb,
				A:        255,
			}
			verts = append(verts, vert)
		}
	}
	return verts
}

// Player

type playerModelComponent struct {
	model         *render.StaticModel
	skin          string
	hasHead       bool
	hasNameTag    bool
	isFirstPerson bool

	dir        float64
	time       float64
	idleTime   float64
	manualMove bool
	walking    bool

	armTime float64

	heldModel *render.StaticModel
	heldMat   mgl32.Mat4
}

func (p *playerModelComponent) Model() *render.StaticModel { return p.model }

func (p *playerModelComponent) SwingArm() {
	p.armTime = 15
}

func (p *playerModelComponent) SetCurrentItem(item *ItemStack) {
	if p.heldModel != nil {
		p.heldModel.Free()
	}
	if item == nil {
		return
	}
	mdl := getModel(item.Type.Name())
	if mdl == nil {
		mdl = getModel("stone")
	}

	var blk Block
	if bt, ok := item.Type.(*blockItem); ok {
		blk = bt.block
	}
	mode := "thirdperson"

	var out []*render.StaticVertex
	if mdl.builtIn == builtInGenerated {
		out, p.heldMat = genStaticModelFromItem(mdl, blk, mode)
	} else if mdl.builtIn == builtInFalse {
		out, p.heldMat = staticModelFromItem(mdl, blk, mode)
	}

	p.heldModel = render.NewStaticModel([][]*render.StaticVertex{
		out,
	})
	p.heldModel.Radius = 3
}

type PlayerModelComponent interface {
	SwingArm()
	SetCurrentItem(item *ItemStack)
}

const (
	playerModelHead = iota
	playerModelBody
	playerModelLegLeft
	playerModelLegRight
	playerModelArmLeft
	playerModelArmRight
	playerModelNameTag
)

func esPlayerModelAdd(p *playerModelComponent, pl PlayerComponent) {
	uuid := pl.UUID()
	info := Client.playerList.info[uuid]
	if info == nil {
		panic("missing player info")
	}
	skin := info.skin
	p.skin = info.skinHash
	if p.skin != "" {
		render.RefSkin(p.skin)
	}

	var hverts []*render.StaticVertex
	if p.hasHead {
		hverts = appendBox(hverts, -4/16.0, 0, -4/16.0, 8/16.0, 8/16.0, 8/16.0, [6]render.TextureInfo{
			direction.North: skin.Sub(8, 8, 8, 8),
			direction.South: skin.Sub(24, 8, 8, 8),
			direction.West:  skin.Sub(0, 8, 8, 8),
			direction.East:  skin.Sub(16, 8, 8, 8),
			direction.Up:    skin.Sub(8, 0, 8, 8),
			direction.Down:  skin.Sub(16, 0, 8, 8),
		})
		hverts = appendBox(hverts, -4.2/16.0, -.2/16.0, -4.2/16.0, 8.4/16.0, 8.4/16.0, 8.4/16.0, [6]render.TextureInfo{
			direction.North: skin.Sub(8+32, 8, 8, 8),
			direction.South: skin.Sub(24+32, 8, 8, 8),
			direction.West:  skin.Sub(0+32, 8, 8, 8),
			direction.East:  skin.Sub(16+32, 8, 8, 8),
			direction.Up:    skin.Sub(8+32, 0, 8, 8),
			direction.Down:  skin.Sub(16+32, 0, 8, 8),
		})
	}

	bverts := appendBox(nil, -4/16.0, -6/16.0, -2/16.0, 8/16.0, 12/16.0, 4/16.0, [6]render.TextureInfo{
		direction.North: skin.Sub(20, 20, 8, 12),
		direction.South: skin.Sub(32, 20, 8, 12),
		direction.West:  skin.Sub(16, 20, 4, 12),
		direction.East:  skin.Sub(28, 20, 4, 12),
		direction.Up:    skin.Sub(20, 16, 8, 4),
		direction.Down:  skin.Sub(28, 16, 8, 4),
	})
	bverts = appendBox(bverts, -4.2/16.0, -6.2/16.0, -2.2/16.0, 8.4/16.0, 12.4/16.0, 4.4/16.0, [6]render.TextureInfo{
		direction.North: skin.Sub(20, 20+16, 8, 12),
		direction.South: skin.Sub(32, 20+16, 8, 12),
		direction.West:  skin.Sub(16, 20+16, 4, 12),
		direction.East:  skin.Sub(28, 20+16, 4, 12),
		direction.Up:    skin.Sub(20, 16+16, 8, 4),
		direction.Down:  skin.Sub(28, 16+16, 8, 4),
	})

	var lverts [4][]*render.StaticVertex

	for i, off := range [][4]int{
		{0, 16, 0, 32},
		{16, 48, 0, 48},
		{32, 48, 48, 48},
		{40, 16, 40, 32},
	} {
		ox, oy := off[0], off[1]
		lverts[i] = appendBox(nil, -2/16.0, -12/16.0, -2/16.0, 4/16.0, 12/16.0, 4/16.0, [6]render.TextureInfo{
			direction.North: skin.Sub(ox+4, oy+4, 4, 12),
			direction.South: skin.Sub(ox+12, oy+4, 4, 12),
			direction.West:  skin.Sub(ox+0, oy+4, 4, 12),
			direction.East:  skin.Sub(ox+8, oy+4, 4, 12),
			direction.Up:    skin.Sub(ox+4, oy, 4, 4),
			direction.Down:  skin.Sub(ox+8, oy, 4, 4),
		})
		ox, oy = off[2], off[3]
		lverts[i] = appendBox(lverts[i], -2.2/16.0, -12.2/16.0, -2.2/16.0, 4.4/16.0, 12.4/16.0, 4.4/16.0, [6]render.TextureInfo{
			direction.North: skin.Sub(ox+4, oy+4, 4, 12),
			direction.South: skin.Sub(ox+12, oy+4, 4, 12),
			direction.West:  skin.Sub(ox+0, oy+4, 4, 12),
			direction.East:  skin.Sub(ox+8, oy+4, 4, 12),
			direction.Up:    skin.Sub(ox+4, oy, 4, 4),
			direction.Down:  skin.Sub(ox+8, oy, 4, 4),
		})
	}

	var nverts []*render.StaticVertex
	if p.hasNameTag {
		nverts = createNameTag(info.name)
	}

	model := render.NewStaticModel([][]*render.StaticVertex{
		playerModelHead:     hverts,
		playerModelBody:     bverts,
		playerModelLegRight: lverts[0],
		playerModelLegLeft:  lverts[1],
		playerModelArmRight: lverts[2],
		playerModelArmLeft:  lverts[3],
		playerModelNameTag:  nverts,
	})
	p.model = model
	model.Radius = 3
}

func createNameTag(name string) (verts []*render.StaticVertex) {
	width := render.SizeOfString(name) + 4
	tex := render.GetTexture("solid")
	for _, v := range faceVertices[direction.North].verts {
		vert := &render.StaticVertex{
			X:        float32(v.X)*float32(width*0.01) - float32((width/2)*0.01),
			Y:        float32(v.Y)*0.2 - 0.1,
			Texture:  tex,
			TextureX: float64(v.TOffsetX),
			TextureY: float64(v.TOffsetY),
			R:        0,
			G:        0,
			B:        0,
			A:        100,
		}
		verts = append(verts, vert)
	}
	offset := -(width/2)*0.01 + (2 * 0.01)
	for _, r := range name {
		tex := render.CharacterTexture(r)
		if tex == nil {
			continue
		}
		s := render.SizeOfCharacter(r)
		for _, v := range faceVertices[direction.North].verts {
			vert := &render.StaticVertex{
				X:        float32(v.X)*float32(s*0.01) - float32(offset+s*0.01),
				Y:        float32(v.Y)*0.16 - 0.08,
				Z:        -0.01,
				Texture:  tex,
				TextureX: float64(v.TOffsetX),
				TextureY: float64(v.TOffsetY),
				R:        255,
				G:        255,
				B:        255,
				A:        255,
			}
			verts = append(verts, vert)
		}
		offset += (s + 2) * 0.01
	}
	return verts
}

func esPlayerModelRemove(p *playerModelComponent) {
	if p.skin != "" {
		render.FreeSkin(p.skin)
	}
	if p.heldModel != nil {
		p.heldModel.Free()
	}
}

func esModelRemove(p interface {
	Model() *render.StaticModel
}) {
	if p.Model() != nil {
		p.Model().Free()
	}
}

var moveLimit = 1e-5

func esPlayerModelTick(p *playerModelComponent,
	pos PositionComponent, t *targetPositionComponent, r RotationComponent) {
	x, y, z := pos.Position()
	model := p.model

	model.X, model.Y, model.Z = -float32(x), -float32(y), float32(z)
	if p.heldModel != nil {
		p.heldModel.X, p.heldModel.Y, p.heldModel.Z = -float32(x), -float32(y), float32(z)
	}

	offMat := mgl32.Translate3D(float32(x), -float32(y), float32(z)).
		Mul4(mgl32.Rotate3DY(math.Pi - float32(r.Yaw())).Mat4())

	// TODO This isn't the most optimal way of doing this
	if p.hasNameTag {
		val := math.Atan2(x-render.Camera.X, z-render.Camera.Z)
		model.Matrix[playerModelNameTag] = mgl32.Translate3D(float32(x), -float32(y), float32(z)).
			Mul4(mgl32.Translate3D(0, -12/16.0-12/16.0-0.6, 0)).
			Mul4(mgl32.Rotate3DY(float32(val)).Mat4())
	}

	model.Matrix[playerModelHead] = offMat.Mul4(mgl32.Translate3D(0, -12/16.0-12/16.0, 0)).
		Mul4(mgl32.Rotate3DX(float32(r.Pitch())).Mat4())
	model.Matrix[playerModelBody] = offMat.Mul4(mgl32.Translate3D(0, -12/16.0-6/16.0, 0))

	time := p.time
	dir := p.dir
	if dir == 0 {
		dir = 1
		time = 15
	}
	ang := ((time / 15) - 1) * (math.Pi / 4)

	model.Matrix[playerModelLegRight] = offMat.Mul4(mgl32.Translate3D(2/16.0, -12/16.0, 0)).
		Mul4(mgl32.Rotate3DX(float32(ang)).Mat4())
	model.Matrix[playerModelLegLeft] = offMat.Mul4(mgl32.Translate3D(-2/16.0, -12/16.0, 0)).
		Mul4(mgl32.Rotate3DX(-float32(ang)).Mat4())

	iTime := p.idleTime
	iTime += Client.delta * 0.02
	if iTime > math.Pi*2 {
		iTime -= math.Pi * 2
	}
	p.idleTime = iTime

	if p.armTime <= 0 {
		p.armTime = 0
	} else {
		p.armTime -= Client.delta
	}

	model.Matrix[playerModelArmRight] = offMat.Mul4(mgl32.Translate3D(6/16.0, -12/16.0-12/16.0, 0))
	model.Matrix[playerModelArmRight] = model.Matrix[playerModelArmRight].
		Mul4(mgl32.Rotate3DX(-float32(ang * 0.75)).Mat4()).
		Mul4(mgl32.Rotate3DZ(float32(math.Cos(iTime)*0.06) - 0.06).Mat4()).
		Mul4(mgl32.Rotate3DX(float32(math.Sin(iTime)*0.06) - float32((7.5-math.Abs(p.armTime-7.5))/7.5)).Mat4())

	if p.heldModel != nil {
		p.heldModel.Matrix[0] = offMat.Mul4(mgl32.Translate3D(6/16.0, -12/16.0-12/16.0, 0.0)).
			Mul4(mgl32.Rotate3DX(-float32(ang * 0.75)).Mat4()).
			Mul4(mgl32.Rotate3DZ(float32(math.Cos(iTime)*0.06) - 0.06).Mat4()).
			Mul4(mgl32.Rotate3DX(float32(math.Sin(iTime)*0.06) - float32((7.5-math.Abs(p.armTime-7.5))/7.5)).Mat4()).
			Mul4(mgl32.Translate3D(0, 11/16.0, -5/16.0)).
			Mul4(mgl32.Rotate3DX(math.Pi).Mat4()).
			Mul4(p.heldMat)
	}

	model.Matrix[playerModelArmLeft] = offMat.Mul4(mgl32.Translate3D(-6/16.0, -12/16.0-12/16.0, 0)).
		Mul4(mgl32.Rotate3DX(float32(ang * 0.75)).Mat4()).
		Mul4(mgl32.Rotate3DZ(-float32(math.Cos(iTime)*0.06) + 0.06).Mat4()).
		Mul4(mgl32.Rotate3DX(-float32(math.Sin(iTime) * 0.06)).Mat4())

	update := true
	if (!p.manualMove && t.X == t.sX && t.Y == t.sY && t.Z == t.sZ) || (p.manualMove && !p.walking) {
		if t.stillTime > 5.0 {
			if math.Abs(time-15) <= 1.5*Client.delta {
				time = 15
				update = false
			}
			dir = math.Copysign(1, 15-time)
		} else {
			t.stillTime += Client.delta
		}
	} else {
		t.stillTime = 0
	}

	if update {
		time += Client.delta * 1.5 * dir
		if time > 30 {
			time = 30
			dir = -1
		} else if time < 0 {
			time = 0
			dir = 1
		}
	}
	p.dir = dir
	p.time = time
}
