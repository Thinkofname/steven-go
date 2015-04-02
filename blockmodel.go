package main

import (
	"encoding/json"
	"fmt"
	"image"
	"math"
	"strings"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
	"github.com/thinkofdeath/steven/type/direction"
)

var (
	blockStateModels = map[pluginKey]*blockStateModel{}
)

type blockStateModel struct {
	variants map[string]blockVariants
}

type blockVariants []*blockModel

func findStateModel(plugin, name string) *blockStateModel {
	key := pluginKey{plugin, name}
	if bs, ok := blockStateModels[key]; ok {
		return bs
	}

	// Hack to add our 'missing block' into the game without a
	// model for it. We hijack the clay model and then replace
	// the textures for it with ours.
	if plugin == "steven" && name == "missing_block" {
		key.Plugin = "minecraft"
		key.Name = "clay"
	}
	bs := loadStateModel(key)

	// See above comment
	if plugin == "steven" && name == "missing_block" {
		v := bs.variants["normal"]
		for _, m := range v {
			for _, e := range m.elements {
				for i := range e.faces {
					e.faces[i].texture = "missing_texture"
				}
			}
		}
		key.Plugin = "steven"
		key.Name = "missing_block"
	}

	blockStateModels[key] = bs
	return bs
}

func loadStateModel(key pluginKey) *blockStateModel {
	data := map[string]interface{}{}
	err := loadJSON(key.Plugin, fmt.Sprintf("blockstates/%s.json", key.Name), &data)
	if err != nil {
		fmt.Printf("Error loading model: %s\n", err)
		return nil
	}
	bs := &blockStateModel{
		variants: map[string]blockVariants{},
	}
	variants := data["variants"].(map[string]interface{})
	for k, v := range variants {
		var models []*blockModel
		switch v := v.(type) {
		case map[string]interface{}:
			models = append(models, parseBlockStateVariant(key.Plugin, v))
		case []interface{}:
			for _, vv := range v {
				models = append(models, parseBlockStateVariant(key.Plugin, vv.(map[string]interface{})))
			}
		default:
			fmt.Printf("Unhandled variant type: %T\n", v)
		}
		bs.variants[k] = models
	}
	return bs
}

func (bs *blockStateModel) variant(key string) blockVariants {
	return bs.variants[key]
}

func (bv blockVariants) selectModel(index int) *blockModel {
	return bv[uint(index)%uint(len(bv))]
}

type blockModel struct {
	textureVars      map[string]string
	elements         []*blockElement
	ambientOcclusion bool
	aoSet            bool

	y, x float64
}

func parseBlockStateVariant(plugin string, data map[string]interface{}) *blockModel {
	modelName := data["model"].(string)
	bdata := map[string]interface{}{}
	err := loadJSON(plugin, "models/block/"+modelName+".json", &bdata)
	if err != nil {
		fmt.Printf("Error loading model: %s\n", err)
		return nil
	}
	bm := parseBlockModel(plugin, bdata)

	bm.y, _ = data["y"].(float64)
	bm.x, _ = data["x"].(float64)
	return bm
}

func parseBlockModel(plugin string, data map[string]interface{}) *blockModel {
	var bm *blockModel
	if parent, ok := data["parent"].(string); ok {
		pdata := map[string]interface{}{}
		err := loadJSON(plugin, "models/"+parent+".json", &pdata)
		if err != nil {
			fmt.Printf("Error loading model: %s\n", err)
			return nil
		}
		bm = parseBlockModel(plugin, pdata)
	} else {
		bm = &blockModel{
			textureVars: map[string]string{},
		}
	}

	if textures, ok := data["textures"].(map[string]interface{}); ok {
		for k, v := range textures {
			bm.textureVars[k] = v.(string)
		}
	}

	if elements, ok := data["elements"].([]interface{}); ok {
		for _, e := range elements {
			bm.elements = append(bm.elements, parseBlockElement(e.(map[string]interface{})))
		}
	}

	ambientOcclusion, ok := data["ambientocclusion"].(bool)
	if ok {
		bm.ambientOcclusion = ambientOcclusion
		bm.aoSet = true
	} else if !bm.aoSet {
		bm.ambientOcclusion = true
	}

	return bm
}

type blockElement struct {
	from, to [3]float64
	shade    bool
	rotation *blockRotation

	faces [6]*blockFace
}

type blockRotation struct {
	origin  []float64
	axis    string
	angle   float64
	rescale bool
}

type blockFace struct {
	uv          [4]int
	texture     string
	textureInfo *render.TextureInfo
	cullFace    direction.Type
	rotation    int
	tintIndex   int
}

func parseBlockElement(data map[string]interface{}) *blockElement {
	be := &blockElement{}
	from := data["from"].([]interface{})
	to := data["to"].([]interface{})
	for i := 0; i < 3; i++ {
		be.from[i] = from[i].(float64)
		be.to[i] = to[i].(float64)
	}

	shade, ok := data["shade"].(bool)
	be.shade = !ok || shade

	if faces, ok := data["faces"].(map[string]interface{}); ok {
		for i, d := range direction.Values {
			if data, ok := faces[d.String()].(map[string]interface{}); ok {
				be.faces[i] = &blockFace{}
				be.faces[i].init(data)
			}
		}
	}

	if rotation, ok := data["rotation"].(map[string]interface{}); ok {
		r := &blockRotation{}
		be.rotation = r

		r.origin = []float64{8, 8, 8}
		if origin, ok := rotation["origin"].([]interface{}); ok {
			r.origin[0] = origin[0].(float64)
			r.origin[1] = origin[1].(float64)
			r.origin[2] = origin[2].(float64)
		}
		r.axis = rotation["axis"].(string)
		r.angle = rotation["angle"].(float64)
		r.rescale, _ = rotation["rescale"].(bool)
	}

	return be
}

func (bf *blockFace) init(data map[string]interface{}) {
	if uv, ok := data["uv"].([]interface{}); ok {
		for i := 0; i < 4; i++ {
			bf.uv[i] = int(uv[i].(float64))
		}
	} else {
		bf.uv = [4]int{0, 0, 16, 16}
	}
	bf.texture, _ = data["texture"].(string)
	cullFace, _ := data["cullface"].(string)
	bf.cullFace = direction.FromString(cullFace)
	rotation, ok := data["rotation"].(float64)
	if ok {
		bf.rotation = int(rotation)
	}
	bf.tintIndex = -1
	tintIndex, ok := data["tintindex"].(float64)
	if ok {
		bf.tintIndex = int(tintIndex)
	}
}

// Precomputed face vertices
var faceVertices = [6][6]chunkVertex{
	{ // Up
		{X: 0, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},
		{X: 1, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},
		{X: 0, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 1},

		{X: 1, Y: 1, Z: 1, TOffsetX: 1, TOffsetY: 1},
		{X: 0, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 1},
		{X: 1, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},
	},
	{ // Down
		{X: 0, Y: 0, Z: 0, TOffsetX: 0, TOffsetY: 0},
		{X: 0, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 1},
		{X: 1, Y: 0, Z: 0, TOffsetX: 1, TOffsetY: 0},

		{X: 1, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},
		{X: 1, Y: 0, Z: 0, TOffsetX: 1, TOffsetY: 0},
		{X: 0, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 1},
	},
	{ // North
		{X: 0, Y: 0, Z: 0, TOffsetX: 0, TOffsetY: 1},
		{X: 1, Y: 0, Z: 0, TOffsetX: 1, TOffsetY: 1},
		{X: 0, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},

		{X: 1, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},
		{X: 0, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},
		{X: 1, Y: 0, Z: 0, TOffsetX: 1, TOffsetY: 1},
	},
	{ // South
		{X: 0, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 1},
		{X: 0, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 0},
		{X: 1, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},

		{X: 1, Y: 1, Z: 1, TOffsetX: 1, TOffsetY: 0},
		{X: 1, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},
		{X: 0, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 0},
	},
	{ // West
		{X: 0, Y: 0, Z: 0, TOffsetX: 0, TOffsetY: 1},
		{X: 0, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},
		{X: 0, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},

		{X: 0, Y: 1, Z: 1, TOffsetX: 1, TOffsetY: 0},
		{X: 0, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},
		{X: 0, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},
	},
	{ // East
		{X: 1, Y: 0, Z: 0, TOffsetX: 0, TOffsetY: 1},
		{X: 1, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},
		{X: 1, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},

		{X: 1, Y: 1, Z: 1, TOffsetX: 1, TOffsetY: 0},
		{X: 1, Y: 1, Z: 0, TOffsetX: 0, TOffsetY: 0},
		{X: 1, Y: 0, Z: 1, TOffsetX: 1, TOffsetY: 1},
	},
}

var faceRotation = [...]direction.Type{
	direction.North,
	direction.East,
	direction.South,
	direction.West,
}

func (bm *blockModel) render(x, y, z int, bs *blocksSnapshot) []chunkVertex {
	this := bs.block(x, y, z)
	var out []chunkVertex
	for ei := range bm.elements {
		el := bm.elements[len(bm.elements)-1-ei]
	faceLoop:
		for i := range faceVertices {
			faceID := i
			face := el.faces[i]
			if face == nil {
				continue
			}
			cullFace := face.cullFace
			if bm.y > 0 {
				if cullFace >= 2 {
					var pos int
					for di, d := range faceRotation {
						if d == cullFace {
							pos = di
							break
						}
					}
					cullFace = faceRotation[(pos+(int(bm.y)/90))%len(faceRotation)]
				}
				if faceID >= 2 {
					var pos int
					for di, d := range faceRotation {
						if d == direction.Type(faceID) {
							pos = di
							break
						}
					}
					faceID = int(faceRotation[(pos+(int(bm.y)/90))%len(faceRotation)])
				}
			}
			if cullFace != direction.Invalid {
				ox, oy, oz := cullFace.Offset()
				if b := bs.block(x+ox, y+oy, z+oz); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			}
			vert := faceVertices[i]
			tex := face.textureInfo
			if tex == nil {
				tex = bm.lookupTexture(face.texture)
				face.textureInfo = tex
			}

			var cr, cg, cb byte
			cr = 255
			cg = 255
			cb = 255
			if this.TintImage() != nil {
				switch face.tintIndex {
				case 0:
					cr, cg, cb = calculateBiome(bs, x, z, this.TintImage())
				}
			}

			ux1 := int16(face.uv[0] * tex.Width)
			ux2 := int16(face.uv[2] * tex.Width)
			uy1 := int16(face.uv[1] * tex.Height)
			uy2 := int16(face.uv[3] * tex.Height)

			var minX, minY, minZ int16 = math.MaxInt16, math.MaxInt16, math.MaxInt16
			var maxX, maxY, maxZ int16 = math.MinInt16, math.MinInt16, math.MinInt16

			for v := range vert {
				vert[v].TX = uint16(tex.X)
				vert[v].TY = uint16(tex.Y + tex.Atlas*1024.0)
				vert[v].TW = uint16(tex.Width)
				vert[v].TH = uint16(tex.Height)
				vert[v].R = cr
				vert[v].G = cg
				vert[v].B = cb

				if vert[v].X == 0 {
					vert[v].X = int16(el.from[0] * 16)
				} else {
					vert[v].X = int16(el.to[0] * 16)
				}
				if vert[v].Y == 0 {
					vert[v].Y = int16(el.from[1] * 16)
				} else {
					vert[v].Y = int16(el.to[1] * 16)
				}
				if vert[v].Z == 0 {
					vert[v].Z = int16(el.from[2] * 16)
				} else {
					vert[v].Z = int16(el.to[2] * 16)
				}

				if bm.y > 0 {
					rotY := bm.y * (math.Pi / 180)
					c := int16(math.Cos(rotY))
					s := int16(math.Sin(rotY))
					x := vert[v].X - 8*16
					z := vert[v].Z - 8*16
					vert[v].X = 8*16 + int16(x*c-z*s)
					vert[v].Z = 8*16 + int16(z*c+x*s)
				}

				if el.rotation != nil {
					r := el.rotation
					switch r.axis {
					case "y":
						rotY := r.angle * (math.Pi / 180)
						c := math.Cos(rotY)
						s := math.Sin(rotY)
						x := float64(vert[v].X) - r.origin[0]*16
						z := float64(vert[v].Z) - r.origin[2]*16
						vert[v].X = int16(r.origin[0] + (x*c - z*s))
						vert[v].Z = int16(r.origin[2] + (z*c + x*s))
					}
				}

				vert[v].X += int16(x * 256)
				vert[v].Y += int16(y * 256)
				vert[v].Z += int16(z * 256)

				if vert[v].TOffsetX == 0 {
					vert[v].TOffsetX = int16(ux1)
				} else {
					vert[v].TOffsetX = int16(ux2)
				}
				if vert[v].TOffsetY == 0 {
					vert[v].TOffsetY = int16(uy1)
				} else {
					vert[v].TOffsetY = int16(uy2)
				}

				if el.rotation != nil && el.rotation.rescale {
					if vert[v].X < minX {
						minX = vert[v].X
					} else if vert[v].X > maxX {
						maxX = vert[v].X
					}
					if vert[v].Y < minY {
						minY = vert[v].Y
					} else if vert[v].Y > maxY {
						maxY = vert[v].Y
					}
					if vert[v].Z < minZ {
						minZ = vert[v].Z
					} else if vert[v].Z > maxZ {
						maxZ = vert[v].Z
					}
				}
			}

			if el.rotation != nil && el.rotation.rescale {
				diffX := float64(maxX - minX)
				diffY := float64(maxY - minY)
				diffZ := float64(maxZ - minZ)
				for v := range vert {
					vert[v].X = int16(x*256) + int16((float64(vert[v].X-minX)/diffX)*256)
					vert[v].Y = int16(y*256) + int16((float64(vert[v].Y-minY)/diffY)*256)
					vert[v].Z = int16(z*256) + int16((float64(vert[v].Z-minZ)/diffZ)*256)
				}
			}

			// Process lighting last, after all operations are applied
			for v := range vert {
				vert[v].BlockLight, vert[v].SkyLight = calculateLight(
					bs,
					x, y, z,
					float64(vert[v].X)/256.0,
					float64(vert[v].Y)/256.0,
					float64(vert[v].Z)/256.0,
					faceID, bm.ambientOcclusion, this.ForceShade(),
				)
			}

			out = append(out, vert[:]...)
		}
	}
	return out
}

// Takes an average of the biome colors of the surrounding area
func calculateBiome(bs *blocksSnapshot, x, z int, img *image.NRGBA) (byte, byte, byte) {
	count := 0
	var r, g, b int
	for xx := -2; xx <= 2; xx++ {
		for zz := -2; zz <= 2; zz++ {
			biome := bs.biome(x+xx, z+zz)
			ix := biome.ColorIndex & 0xFF
			iy := biome.ColorIndex >> 8
			col := img.NRGBAAt(ix, iy)
			r += int(col.R)
			g += int(col.G)
			b += int(col.B)
			count++
		}
	}
	return byte(r / count), byte(g / count), byte(b / count)
}

func calculateLight(bs *blocksSnapshot, origX, origY, origZ int,
	x, y, z float64, face int, smooth, force bool) (byte, byte) {
	blockLight := bs.blockLight(origX, origY, origZ)
	skyLight := bs.skyLight(origX, origY, origZ)
	if !smooth {
		return blockLight, skyLight
	}
	count := 1

	// TODO(Think) Document/cleanup this
	// it was taken from and older renderer of mine
	// (thinkmap).

	var pox, poy, poz, nox, noy, noz int

	switch face {
	case 0: // Up
		poz, pox = 0, 0
		noz, nox = -1, -1
		poy = 1
		noy = 0
	case 1: // Down
		poz, pox = 0, 0
		noz, nox = -1, -1
		poy = -1
		noy = -2
	case 2: // North
		poy, pox = 0, 0
		noy, nox = -1, -1
		poz = -1
		noz = -2
	case 3: // South
		poy, pox = 0, 0
		noy, nox = -1, -1
		poz = 1
		noz = 0
	case 4: // West
		poz, poy = 0, 0
		noz, noy = -1, -1
		pox = -1
		nox = -2
	case 5: // East
		poz, poy = 0, 0
		noz, noy = -1, -1
		pox = 1
		nox = 0
	}
	for ox := nox; ox <= pox; ox++ {
		for oy := noy; oy <= poy; oy++ {
			for oz := noz; oz <= poz; oz++ {
				bx := round(x + float64(ox))
				by := round(y + float64(oy))
				bz := round(z + float64(oz))
				count++
				blockLight += bs.blockLight(bx, by, bz)
				if !force {
					skyLight += bs.skyLight(bx, by, bz)
				} else if bl := bs.block(bx, by, bz); bl.Is(BlockAir) {
					skyLight += 15
				}
			}
		}

	}

	return blockLight / byte(count), skyLight / byte(count)
}

func round(f float64) int {
	return int(f + math.Copysign(0.5, f))
}

func (bm *blockModel) lookupTexture(name string) *render.TextureInfo {
	if len(name) > 0 && name[0] == '#' {
		return bm.lookupTexture(bm.textureVars[name[1:]])
	}
	if strings.HasPrefix(name, "blocks/") {
		name = name[len("blocks/"):]
	}
	return render.GetTexture(name)
}

func loadJSON(plugin, name string, target interface{}) error {
	r, err := resource.Open(plugin, name)
	if err != nil {
		return err
	}
	defer r.Close()
	d := json.NewDecoder(r)
	return d.Decode(target)
}
