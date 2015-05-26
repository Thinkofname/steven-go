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
	"log"
	"math"
	"strings"

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
		log.Printf("Error loading state %s: %s\n", key.Name, err)
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

type builtInType int

const (
	builtInFalse = iota
	builtInGenerated
	builtInEntity
	builtInCompass
	builtInClock
)

type model struct {
	textureVars      map[string]string
	elements         []*modelElement
	ambientOcclusion bool
	aoSet            bool

	uvLock bool
	y, x   float64

	// Item specific features
	display map[string]modelDisplay
	builtIn builtInType
}

type modelDisplay struct {
	Rotation    *[3]float64
	Translation *[3]float64
	Scale       *[3]float64
}

func parseBlockStateVariant(plugin string, js realjson.RawMessage) *model {
	type jsType struct {
		Model  string
		X, Y   float64
		UVLock bool
	}
	var data jsType
	err := json.Unmarshal(js, &data)
	if err != nil {
		log.Println(err)
		return nil
	}
	var bdata jsModel
	err = loadJSON(plugin, "models/block/"+data.Model+".json", &bdata)
	if err != nil {
		return nil
	}
	bm := parseModel(plugin, &bdata)

	bm.y = data.Y
	bm.x = data.X
	bm.uvLock = data.UVLock
	return bm
}

type jsModel struct {
	Parent           string
	Textures         map[string]string
	AmbientOcclusion *bool
	Elements         []*jsBlockElement
	Display          map[string]modelDisplay
}

func parseModel(plugin string, data *jsModel) *model {
	var bm *model
	if data.Parent != "" && !strings.HasPrefix(data.Parent, "builtin/") {
		var pdata jsModel
		err := loadJSON(plugin, "models/"+data.Parent+".json", &pdata)
		if err != nil {
			log.Printf("Error loading model %s: %s\n", data.Parent, err)
			return nil
		}
		bm = parseModel(plugin, &pdata)
	} else {
		bm = &model{
			textureVars: map[string]string{},
			display:     map[string]modelDisplay{},
		}
		if strings.HasPrefix(data.Parent, "builtin/") {
			switch data.Parent {
			case "builtin/generated":
				bm.builtIn = builtInGenerated
			case "builtin/entity":
				bm.builtIn = builtInEntity
			case "builtin/compass":
				bm.builtIn = builtInCompass
			case "builtin/clock":
				bm.builtIn = builtInClock
			}
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

	if data.Display != nil {
		for k, v := range data.Display {
			bm.display[k] = v
		}
	}

	return bm
}

type modelElement struct {
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

func parseBlockElement(data *jsBlockElement) *modelElement {
	be := &modelElement{}
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

func (bm *model) lookupTexture(name string) render.TextureInfo {
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
	err = realjson.NewDecoder(r).Decode(target)
	if err != nil {
		// Take the slow path through our preprocessor.
		// Hopefully this can be removed in later minecraft versions.d, err := ioutil.ReadAll(r)
		r.Close()
		r, err = resource.Open(plugin, name)
		if err != nil {
			return err
		}

		d, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		return json.Unmarshal(d, target)
	}
	return err
}
