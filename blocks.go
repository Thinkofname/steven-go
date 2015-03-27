package main

// Valid blocks.
var (
	BlockAir     = initSimple("air", simpleConfig{})
	Stone        = initSimple("stone", simpleConfig{Color: 0xC0C0C0})
	Grass        = initSimple("grass", simpleConfig{Color: 0x00FF00})
	Dirt         = initSimple("dirt", simpleConfig{Color: 0xAE6000})
	Cobblestone  = initSimple("cobblestone", simpleConfig{Color: 0x474747})
	Planks       = initSimple("planks", simpleConfig{Color: 0xE3BC62})
	Sapling      = initSimple("sapling", simpleConfig{Color: 0xFFFFFF})
	Bedrock      = initSimple("bedrock", simpleConfig{Color: 0x000000})
	FlowingWater = initLiquid("flowing_water", false)
	Water        = initLiquid("water", false)
	FlowingLava  = initLiquid("flowing_lava", true)
	Lava         = initLiquid("lava", true)
)
