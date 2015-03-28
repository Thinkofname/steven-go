package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/thinkofdeath/steven/render"
	"github.com/thinkofdeath/steven/resource"
)

var (
	blockStateModels  = map[pluginKey]*blockStateModel{}
	blockStateLock    sync.RWMutex
	blockStateWaiters = map[pluginKey]*sync.WaitGroup{}
)

type blockStateModel struct {
	variants map[string][]*blockModel
}

func findStateModel(plugin, name string) *blockStateModel {
	key := pluginKey{plugin, name}
	blockStateLock.RLock()
	if bs, ok := blockStateModels[key]; ok {
		blockStateLock.RUnlock()
		return bs
	}
	if wg := blockStateWaiters[key]; wg != nil {
		blockStateLock.RUnlock()
		wg.Wait()

		blockStateLock.RLock()
		defer blockStateLock.RUnlock()
		return blockStateModels[key]
	}
	blockStateLock.RUnlock()
	blockStateLock.Lock()

	// Re-run the above checks in case it was completed between switching
	// locks
	if bs := blockStateModels[key]; bs != nil {
		blockStateLock.Unlock()
		return bs
	}
	if wg := blockStateWaiters[key]; wg != nil {
		blockStateLock.Unlock()
		wg.Wait()

		blockStateLock.RLock()
		defer blockStateLock.RUnlock()
		return blockStateModels[key]
	}

	var wg sync.WaitGroup
	blockStateWaiters[key] = &wg
	wg.Add(1)
	// No need to continue blocking other builders which
	// don't need this model
	blockStateLock.Unlock()

	fmt.Printf("Load for %s\n", key.String())
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

	// Store the model
	blockStateLock.Lock()
	blockStateModels[key] = bs
	blockStateLock.Unlock()
	// Free anyone waiting
	wg.Done()
	// No longer need the waiter
	blockStateLock.Lock()
	delete(blockStateWaiters, key)
	blockStateLock.Unlock()
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

	faces [6]blockFace
}

type blockFace struct {
	uv        [4]int
	texture   string
	cullFace  string
	rotation  int
	tintIndex int
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
	bf.rotation, _ = data["rotation"].(int)
	bf.tintIndex, _ = data["tintindex"].(int)
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

func (bm *blockModel) render(x, y, z int, get func(x, y, z int) Block) []chunkVertex {
	this := get(x, y, z)
	var out []chunkVertex
	for _, el := range bm.elements {
	faceLoop:
		for i := range faceVertices {
			face := el.faces[i]
			switch face.cullFace {
			case "up":
				if b := get(x, y+1, z); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "down":
				if b := get(x, y-1, z); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "north":
				if b := get(x, y, z-1); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "south":
				if b := get(x, y, z+1); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "east":
				if b := get(x+1, y, z); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			case "west":
				if b := get(x-1, y, z); b.ShouldCullAgainst() || b == this {
					continue faceLoop
				}
			}
			vert := faceVertices[i]
			tex := bm.lookupTexture(face.texture)

			ux1 := int16(face.uv[0] * tex.Width)
			ux2 := int16(face.uv[2] * tex.Width)
			uy1 := int16(face.uv[1] * tex.Height)
			uy2 := int16(face.uv[3] * tex.Height)
			for v := range vert {
				vert[v].TX = uint16(tex.X)
				vert[v].TY = uint16(tex.Y + tex.Atlas*1024.0)
				vert[v].TW = uint16(tex.Width)
				vert[v].TH = uint16(tex.Height)

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

func (bm *blockModel) lookupTexture(name string) render.TextureInfo {
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
