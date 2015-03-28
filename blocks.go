package main

// Valid blocks.
var (
	BlockAir          = initSimple("air", simpleConfig{})
	BlockStone        = initSimple("stone", simpleConfig{Color: 0xC0C0C0})
	BlockGrass        = initSimple("grass", simpleConfig{Color: 0x00FF00})
	BlockDirt         = initSimple("dirt", simpleConfig{Color: 0xAE6000})
	BlockCobblestone  = initSimple("cobblestone", simpleConfig{Color: 0x474747})
	BlockPlanks       = initSimple("planks", simpleConfig{Color: 0xE3BC62})
	BlockSapling      = initSimple("sapling", simpleConfig{Color: 0xFFFFFF})
	BlockBedrock      = initSimple("bedrock", simpleConfig{Color: 0x000000})
	BlockFlowingWater = initLiquid("flowing_water", false)
	BlockWater        = initLiquid("water", false)
	BlockFlowingLava  = initLiquid("flowing_lava", true)
	BlockLava         = initLiquid("lava", true)
)
