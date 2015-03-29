package biome

var byId [256]*Type

// Biomes
var (
	// Snowy
	FrozenOcean        = newBiome(10, 0.0, 0.5)
	FrozenRiver        = newBiome(11, 0.0, 0.5)
	IcePlains          = newBiome(12, 0.0, 0.5)
	IcePlainsSpikes    = newBiome(140, 0.0, 0.5)
	ColdBeach          = newBiome(26, 0.05, 0.3)
	ColdTaiga          = newBiome(30, 0.0, 0.4)
	ColdTaigaMountains = newBiome(158, 0.0, 0.4)
	// Cold
	ExtremeHills              = newBiome(3, 0.2, 0.3)
	ExtremeHillsMountains     = newBiome(131, 0.2, 0.3)
	Taiga                     = newBiome(5, 0.25, 0.8)
	TaigaM                    = newBiome(133, 0.25, 0.8)
	TheEnd                    = newBiome(9, 0.5, 0.5)
	MegaTaiga                 = newBiome(32, 0.3, 0.8)
	MegaSpruceTaiga           = newBiome(160, 0.5, 0.5)
	ExtremeHillsPlus          = newBiome(34, 0.2, 0.3)
	ExtremeHillsPlusMountains = newBiome(162, 0.2, 0.3)
	StoneBeach                = newBiome(25, 0.2, 0.3)
	// Medium/Lush
	Plains                = newBiome(1, 0.5, 0.5)
	SunflowerPlains       = newBiome(129, 0.5, 0.5)
	Forest                = newBiome(4, 0.5, 0.5)
	FlowerForest          = newBiome(132, 0.5, 0.5)
	Swampland             = newBiome(6, 0.8, 0.9)
	SwamplandMountains    = newBiome(134, 0.8, 0.9)
	River                 = newBiome(7, 0.5, 0.5)
	MushroomIsland        = newBiome(14, 0.9, 1.0)
	MushroomISlandShore   = newBiome(15, 0.9, 1.0)
	Beach                 = newBiome(16, 0.8, 0.4)
	Jungle                = newBiome(21, 0.95, 0.8)
	JungleMountains       = newBiome(149, 0.95, 0.9)
	JungleEdge            = newBiome(23, 0.95, 0.8)
	JungleEdgeMountains   = newBiome(151, 0.95, 0.8)
	BirchForest           = newBiome(27, 0.5, 0.5)
	BirchForestMountains  = newBiome(155, 0.5, 0.5)
	RoofedForest          = newBiome(29, 0.5, 0.5)
	RoofedForestMountains = newBiome(157, 0.5, 0.5)
	// Dry/Warm
	Desert                     = newBiome(2, 1.0, 0.0)
	DesertMountain             = newBiome(130, 1.0, 0.0)
	Hell                       = newBiome(8, 1.0, 0.0)
	Savanna                    = newBiome(35, 1.0, 0.0)
	SavannaMountains           = newBiome(163, 1.0, 0.0)
	Mesa                       = newBiome(37, 0.5, 0.5)
	MesaBryce                  = newBiome(165, 0.5, 0.5)
	SavannaPlateau             = newBiome(36, 1.0, 0.0)
	MesaPlateauForest          = newBiome(38, 0.5, 0.5)
	MesaPlateau                = newBiome(39, 0.5, 0.5)
	SavannaPlateauMountains    = newBiome(164, 1.0, 0.0)
	MesaPlateauForestMountains = newBiome(166, 0.5, 0.5)
	MesaPlateauMountains       = newBiome(167, 0.5, 0.5)
	// Neutral
	Ocean                     = newBiome(0, 0.5, 0.5)
	DeepOcean                 = newBiome(24, 0.5, 0.5)
	IceMountains              = newBiome(13, 0.0, 0.5)
	DesertHills               = newBiome(17, 1.0, 0.0)
	ForestHills               = newBiome(18, 0.8, 0.9)
	TaigaHills                = newBiome(19, 0.25, 0.8)
	JungleHills               = newBiome(22, 0.95, 0.9)
	BirchForestHills          = newBiome(28, 0.5, 0.5)
	ColdTaigaHills            = newBiome(31, 0.5, 0.5)
	MegaTaigaHills            = newBiome(33, 0.3, 0.8)
	BirchForestHillsMountains = newBiome(156, 0.5, 0.5)
	MegaSpruceTaigaHills      = newBiome(161, 0.5, 0.5)
	// Custom
	Invalid = newBiome(255, 0.0, 0.0)
)

func ById(id byte) *Type {
	return byId[id]
}

type Type struct {
	ID                    int
	Temperature, Moisture float64
	ColorIndex            int
}

func newBiome(id int, temperature, moisture float64) *Type {
	b := &Type{
		ID:          id,
		Temperature: temperature,
		Moisture:    moisture * temperature,
	}
	bx := int((1.0 - temperature) * 255.0)
	by := int((1.0 - moisture) * 255.0)
	b.ColorIndex = bx | (by << 8)
	byId[id] = b
	return b
}

func init() {
	for i := range byId {
		if byId[i] == nil {
			byId[i] = Invalid
		}
	}
}
