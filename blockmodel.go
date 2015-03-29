package main

import (
	"encoding/json"
	"fmt"
	"image"
	"strings"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
)

var (
	blockStateModels = map[pluginKey]*blockStateModel{}
)

type blockStateModel struct {
	variants map[string][]*blockModel
}

func findStateModel(plugin, name string) *blockStateModel {
	key := pluginKey{plugin, name}
	if bs, ok := blockStateModels[key]; ok {
		return bs
	}

	if plugin == "steven" && name == "missing_block" {
		key.Plugin = "minecraft"
		key.Name = "clay"
	}
	bs := loadStateModel(key)
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
		variants: map[string][]*blockModel{},
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

func (bs *blockStateModel) variant(key string, seed int) *blockModel {
	v := bs.variants[key]
	if v == nil {
		return nil
	}
	return v[uint(seed)%uint(len(v))]
}

type blockModel struct {
	textureVars map[string]string
	elements    []*blockElement
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

	return bm
}

type blockElement struct {
	from, to [3]int
	shade    bool

	faces [6]*blockFace
}

type blockFace struct {
	uv          [4]int
	texture     string
	textureInfo *render.TextureInfo
	cullFace    string
	rotation    int
	tintIndex   int
}

var faceNames = []string{"up", "down", "north", "south", "west", "east"}

func parseBlockElement(data map[string]interface{}) *blockElement {
	be := &blockElement{}
	from := data["from"].([]interface{})
	to := data["to"].([]interface{})
	for i := 0; i < 3; i++ {
		be.from[i] = int(from[i].(float64))
		be.to[i] = int(to[i].(float64))
	}

	shade, ok := data["shade"].(bool)
	be.shade = !ok || shade

	if faces, ok := data["faces"].(map[string]interface{}); ok {
		for i := range faceNames {
			if data, ok := faces[faceNames[i]].(map[string]interface{}); ok {
				be.faces[i] = &blockFace{}
				be.faces[i].init(data)
			}
		}
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
	bf.cullFace, _ = data["cullface"].(string)
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

func (bm *blockModel) render(x, y, z int, bs *blocksSnapshot) []chunkVertex {
	this := bs.block(x, y, z)
	var out []chunkVertex
	for ei := range bm.elements {
		el := bm.elements[len(bm.elements)-1-ei]
	faceLoop:
		for i := range faceVertices {
			face := el.faces[i]
			if face == nil {
				continue
			}
			switch face.cullFace {
			case "up":
				if b := bs.block(x, y+1, z); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "down":
				if b := bs.block(x, y-1, z); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "north":
				if b := bs.block(x, y, z-1); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "south":
				if b := bs.block(x, y, z+1); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "east":
				if b := bs.block(x+1, y, z); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "west":
				if b := bs.block(x-1, y, z); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			}
			vert := faceVertices[i]
			tex := face.textureInfo
			if tex == nil {
				tex = bm.lookupTexture(face.texture)
			}

			var cr, cg, cb byte
			switch face.tintIndex {
			case 0:
				cr, cg, cb = calculateBiome(bs, x, z, this.TintImage())
			default:
				cr = 255
				cg = 255
				cb = 255
			}

			ux1 := int16(face.uv[0] * tex.Width)
			ux2 := int16(face.uv[2] * tex.Width)
			uy1 := int16(face.uv[1] * tex.Height)
			uy2 := int16(face.uv[3] * tex.Height)
			for v := range vert {
				vert[v].TX = uint16(tex.X)
				vert[v].TY = uint16(tex.Y + tex.Atlas*1024.0)
				vert[v].TW = uint16(tex.Width)
				vert[v].TH = uint16(tex.Height)
				vert[v].R = cr
				vert[v].G = cg
				vert[v].B = cb

				if vert[v].X == 0 {
					vert[v].X = int16(el.from[0]*16 + x*256)
				} else {
					vert[v].X = int16(el.to[0]*16 + x*256)
				}
				if vert[v].Y == 0 {
					vert[v].Y = int16(el.from[1]*16 + y*256)
				} else {
					vert[v].Y = int16(el.to[1]*16 + y*256)
				}
				if vert[v].Z == 0 {
					vert[v].Z = int16(el.from[2]*16 + z*256)
				} else {
					vert[v].Z = int16(el.to[2]*16 + z*256)
				}

				vert[v].BlockLight, vert[v].SkyLight = calculateLight(
					bs,
					x, y, z,
					float64(vert[v].X)/256.0,
					float64(vert[v].Y)/256.0,
					float64(vert[v].Z)/256.0,
					i, true, this.ForceShade(),
				)

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
			}

			out = append(out, vert[:]...)
		}
	}
	return out
}

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
				bx := int(x + float64(ox))
				by := int(y + float64(oy))
				bz := int(z + float64(oz))
				count++
				blockLight += bs.blockLight(bx, by, bz)
				if !force {
					skyLight += bs.skyLight(bx, by, bz)
				} else {
					if bl := bs.block(bx, by, bz); bl.Is(BlockAir) {
						skyLight += 15
					}
				}
			}
		}

	}

	return blockLight / byte(count), skyLight / byte(count)
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
