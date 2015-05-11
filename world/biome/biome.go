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

package biome

var byId [256]*Type

// Biomes
var (
	Ocean               = newBiome(0, 0.5, 0.5)
	Plains              = newBiome(1, 0.8, 0.4)
	Desert              = newBiome(2, 2.0, 0.0)
	ExtremeHills        = newBiome(3, 0.2, 0.3)
	Forest              = newBiome(4, 0.7, 0.8)
	Taiga               = newBiome(5, 0.05, 0.8)
	Swampland           = newBiome(6, 0.8, 0.9)
	River               = newBiome(7, 0.5, 0.5)
	Hell                = newBiome(8, 2.0, 0.0)
	TheEnd              = newBiome(9, 0.5, 0.5)
	FrozenOcean         = newBiome(10, 0.0, 0.5)
	FrozenRiver         = newBiome(11, 0.0, 0.5)
	IcePlains           = newBiome(12, 0.0, 0.5)
	IceMountains        = newBiome(13, 0.0, 0.5)
	MushroomIsland      = newBiome(14, 0.9, 1.0)
	MushroomISlandShore = newBiome(15, 0.9, 1.0)
	Beach               = newBiome(16, 0.8, 0.4)
	DesertHills         = newBiome(17, 2.0, 0.0)
	ForestHills         = newBiome(18, 0.7, 0.8)
	TaigaHills          = newBiome(19, 0.2, 0.7)
	ExtremeHillsEdge    = newBiome(20, 0.2, 0.3)
	Jungle              = newBiome(21, 1.2, 0.9)
	JungleHills         = newBiome(22, 1.2, 0.9)
	JungleEdge          = newBiome(23, 0.95, 0.8)
	DeepOcean           = newBiome(24, 0.5, 0.5)
	StoneBeach          = newBiome(25, 0.2, 0.3)
	ColdBeach           = newBiome(26, 0.05, 0.3)
	BirchForest         = newBiome(27, 0.6, 0.6)
	BirchForestHills    = newBiome(28, 0.6, 0.6)
	RoofedForest        = newBiome(29, 0.7, 0.8)
	ColdTaiga           = newBiome(30, -0.5, 0.4)
	ColdTaigaHills      = newBiome(31, -0.5, 0.4)
	MegaTaiga           = newBiome(32, 0.3, 0.8)
	MegaTaigaHills      = newBiome(33, 0.3, 0.8)
	ExtremeHillsPlus    = newBiome(34, 0.2, 0.3)
	Savanna             = newBiome(35, 1.2, 0.0)
	SavannaPlateau      = newBiome(36, 1.0, 0.0)
	Mesa                = newBiome(37, 2.0, 0.0)
	MesaPlateauForest   = newBiome(38, 2.0, 0.0)
	MesaPlateau         = newBiome(39, 2.0, 0.0)

	SunflowerPlains            = newBiome(129, 0.8, 0.4)
	DesertMountain             = newBiome(130, 2.0, 0.0)
	ExtremeHillsMountains      = newBiome(131, 0.2, 0.3)
	FlowerForest               = newBiome(132, 0.7, 0.8)
	TaigaM                     = newBiome(133, 0.05, 0.8)
	SwamplandMountains         = newBiome(134, 0.8, 0.9)
	IcePlainsSpikes            = newBiome(140, 0.0, 0.5)
	JungleMountains            = newBiome(149, 1.2, 0.9)
	JungleEdgeMountains        = newBiome(151, 0.95, 0.8)
	BirchForestMountains       = newBiome(155, 0.6, 0.6)
	BirchForestHillsMountains  = newBiome(156, 0.6, 0.6)
	RoofedForestMountains      = newBiome(157, 0.7, 0.8)
	ColdTaigaMountains         = newBiome(158, -0.5, 0.4)
	MegaSpruceTaiga            = newBiome(160, 0.25, 0.8)
	MegaSpruceTaigaHills       = newBiome(161, 0.3, 0.8)
	ExtremeHillsPlusMountains  = newBiome(162, 0.2, 0.3)
	SavannaMountains           = newBiome(163, 1.2, 0.0)
	SavannaPlateauMountains    = newBiome(164, 1.0, 0.0)
	MesaBryce                  = newBiome(165, 2.0, 0.0)
	MesaPlateauForestMountains = newBiome(166, 2.0, 0.0)
	MesaPlateauMountains       = newBiome(167, 2.0, 0.0)
	//
	Invalid = newBiome(255, 0.0, 0.0)
)

func ById(id byte) *Type {
	if val := byId[id]; val != nil {
		return val
	}
	return Invalid
}

type Type struct {
	ID                    int
	Temperature, Moisture float64
	ColorIndex            int
}

func newBiome(id int, temperature, moisture float64) *Type {
	b := &Type{
		ID:          id,
		Temperature: clamp(temperature, 0, 1),
		Moisture:    clamp(moisture, 0, 1),
	}
	b.Moisture *= b.Temperature
	bx := int((1.0 - b.Temperature) * 255.0)
	by := int((1.0 - b.Moisture) * 255.0)
	b.ColorIndex = bx | (by << 8)
	byId[id] = b
	return b
}

func clamp(x, l, h float64) float64 {
	if x < l {
		return l
	}
	if x > h {
		return h
	}
	return x
}

func init() {
	for i := range byId {
		if byId[i] == nil {
			byId[i] = Invalid
		}
	}
}
