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
	"fmt"
	"io/ioutil"
	"math"

	realjson "encoding/json"

	"github.com/thinkofdeath/steven/encoding/json"

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
	type jsType struct {
		Variants map[string]realjson.RawMessage
	}

	var data jsType
	err := loadJSON(key.Plugin, fmt.Sprintf("blockstates/%s.json", key.Name), &data)
	if err != nil {
		fmt.Printf("Error loading state %s: %s\n", key.Name, err)
		return nil
	}
	bs := &blockStateModel{
		variants: map[string]blockVariants{},
	}
	variants := data.Variants
	for k, v := range variants {
		var models blockVariants
		switch v[0] {
		case '[':
			var list []realjson.RawMessage
			json.Unmarshal(v, &list)
			for _, vv := range list {
				mdl := parseBlockStateVariant(key.Plugin, vv)
				if mdl != nil {
					models = append(models, precomputeModel(mdl))
				}
			}
		default:
			mdl := parseBlockStateVariant(key.Plugin, v)
			if mdl != nil {
				models = append(models, precomputeModel(mdl))
			}
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

	uvLock bool
	y, x   float64
}

func parseBlockStateVariant(plugin string, js realjson.RawMessage) *blockModel {
	type jsType struct {
		Model  string
		X, Y   float64
		UVLock bool
	}
	var data jsType
	err := json.Unmarshal(js, &data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var bdata jsBlockModel
	err = loadJSON(plugin, "models/block/"+data.Model+".json", &bdata)
	if err != nil {
		return nil
	}
	bm := parseBlockModel(plugin, &bdata)

	bm.y = data.Y
	bm.x = data.X
	bm.uvLock = data.UVLock
	return bm
}

type jsBlockModel struct {
	Parent           string
	Textures         map[string]string
	AmbientOcclusion *bool
	Elements         []*jsBlockElement
}

func parseBlockModel(plugin string, data *jsBlockModel) *blockModel {
	var bm *blockModel
	if data.Parent != "" {
		var pdata jsBlockModel
		err := loadJSON(plugin, "models/"+data.Parent+".json", &pdata)
		if err != nil {
			fmt.Printf("Error loading model %s: %s\n", data.Parent, err)
			return nil
		}
		bm = parseBlockModel(plugin, &pdata)
	} else {
		bm = &blockModel{
			textureVars: map[string]string{},
		}
	}

	if data.Textures != nil {
		for k, v := range data.Textures {
			bm.textureVars[k] = v
		}
	}

	for _, e := range data.Elements {
		bm.elements = append(bm.elements, parseBlockElement(e))
	}

	if data.AmbientOcclusion != nil {
		bm.ambientOcclusion = *data.AmbientOcclusion
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
	uv          [4]float64
	texture     string
	textureInfo *render.TextureInfo
	cullFace    direction.Type
	rotation    int
	tintIndex   int
}

type jsBlockElement struct {
	From, To [3]float64
	Shade    *bool
	Faces    map[string]*jsBlockFace
	Rotation *struct {
		Origin  *[3]float64
		Axis    string
		Angle   float64
		Rescale bool
	}
}

func parseBlockElement(data *jsBlockElement) *blockElement {
	be := &blockElement{}
	be.from, be.to = data.From, data.To

	be.shade = data.Shade == nil || *data.Shade

	if data.Faces != nil {
		for i, d := range direction.Values {
			if data, ok := data.Faces[d.String()]; ok {
				be.faces[i] = &blockFace{}
				be.faces[i].init(data)
				if math.IsNaN(be.faces[i].uv[0]) {
					be.faces[i].uv = [4]float64{0, 0, 16, 16}
					switch d {
					case direction.North, direction.South:
						be.faces[i].uv[0] = be.from[0]
						be.faces[i].uv[2] = be.to[0]
						be.faces[i].uv[1] = 16 - be.to[1]
						be.faces[i].uv[3] = 16 - be.from[1]
					case direction.West, direction.East:
						be.faces[i].uv[0] = be.from[2]
						be.faces[i].uv[2] = be.to[2]
						be.faces[i].uv[1] = 16 - be.to[1]
						be.faces[i].uv[3] = 16 - be.from[1]
					case direction.Down, direction.Up:
						be.faces[i].uv[0] = be.from[0]
						be.faces[i].uv[2] = be.to[0]
						be.faces[i].uv[1] = 16 - be.to[2]
						be.faces[i].uv[3] = 16 - be.from[2]
					}
				}
			}
		}
	}

	if data.Rotation != nil {
		r := &blockRotation{}
		be.rotation = r
		rot := data.Rotation

		r.origin = []float64{8, 8, 8}
		if rot.Origin != nil {
			r.origin = rot.Origin[:]
		}
		r.axis = rot.Axis
		r.angle = rot.Angle
		r.rescale = rot.Rescale
	}

	return be
}

type jsBlockFace struct {
	UV        *[4]float64
	Texture   string
	CullFace  string
	Rotation  int
	TintIndex *int
}

func (bf *blockFace) init(data *jsBlockFace) {
	if data.UV != nil {
		bf.uv = *data.UV
	} else {
		bf.uv = [4]float64{math.NaN(), 0, 0, 0}
	}
	bf.texture = data.Texture
	bf.cullFace = direction.FromString(data.CullFace)
	bf.rotation = data.Rotation
	bf.tintIndex = -1
	if data.TintIndex != nil {
		bf.tintIndex = *data.TintIndex
	}
}

func (bm *blockModel) lookupTexture(name string) *render.TextureInfo {
	if len(name) > 0 && name[0] == '#' {
		return bm.lookupTexture(bm.textureVars[name[1:]])
	}
	return render.GetTexture(name)
}

func loadJSON(plugin, name string, target interface{}) error {
	r, err := resource.Open(plugin, name)
	if err != nil {
		return err
	}
	defer r.Close()
	d, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, target)
}
