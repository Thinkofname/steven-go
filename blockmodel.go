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

package main

import (
	"encoding/json"
	"fmt"
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

type blockVariants []*processedModel

func findStateModel(plugin, name string) *blockStateModel {
	key := pluginKey{plugin, name}
	if bs, ok := blockStateModels[key]; ok {
		return bs
	}
	bs := loadStateModel(key)

	if bs == nil {
		blockStateModels[key] = nil
		return nil
	}

	blockStateModels[key] = bs
	return bs
}

func loadStateModel(key pluginKey) *blockStateModel {
	data := map[string]interface{}{}
	err := loadJSON(key.Plugin, fmt.Sprintf("blockstates/%s.json", key.Name), &data)
	if err != nil {
		fmt.Printf("Error loading state %s: %s\n", key.Name, err)
		return nil
	}
	bs := &blockStateModel{
		variants: map[string]blockVariants{},
	}
	variants := data["variants"].(map[string]interface{})
	for k, v := range variants {
		var models blockVariants
		switch v := v.(type) {
		case map[string]interface{}:
			models = append(models, precomputeModel(parseBlockStateVariant(key.Plugin, v)))
		case []interface{}:
			for _, vv := range v {
				models = append(models, precomputeModel(parseBlockStateVariant(key.Plugin, vv.(map[string]interface{}))))
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

func (bv blockVariants) selectModel(index int) *processedModel {
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
		fmt.Printf("Error loading model %s: %s\n", modelName, err)
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
			fmt.Printf("Error loading model %s: %s\n", parent, err)
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
		{X: 1, Y: 0, Z: 0, TOffsetX: 1, TOffsetY: 1},
		{X: 1, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 1},
		{X: 1, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},

		{X: 1, Y: 1, Z: 1, TOffsetX: 0, TOffsetY: 0},
		{X: 1, Y: 1, Z: 0, TOffsetX: 1, TOffsetY: 0},
		{X: 1, Y: 0, Z: 1, TOffsetX: 0, TOffsetY: 1},
	},
}

var faceRotation = [...]direction.Type{
	direction.North,
	direction.East,
	direction.South,
	direction.West,
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
