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
	"bytes"
	"reflect"
	"unicode"
)

// Valid blocks.
var Blocks = struct {
	Air                        *BlockSet `cullAgainst:"false" collidable:"false" renderable:"false"`
	Stone                      *BlockSet `type:"stone"`
	Grass                      *BlockSet `type:"grass"`
	Dirt                       *BlockSet
	Cobblestone                *BlockSet
	Planks                     *BlockSet `type:"planks"`
	Sapling                    *BlockSet `type:"sapling"`
	Bedrock                    *BlockSet `hardness:"Inf"`
	FlowingWater               *BlockSet `type:"liquid"`
	Water                      *BlockSet `type:"liquid"`
	FlowingLava                *BlockSet `type:"liquid" lava:"true"`
	Lava                       *BlockSet `type:"liquid" lava:"true"`
	Sand                       *BlockSet `hardness:"0.5"`
	Gravel                     *BlockSet
	GoldOre                    *BlockSet
	IronOre                    *BlockSet
	CoalOre                    *BlockSet
	Log                        *BlockSet `type:"log"`
	Leaves                     *BlockSet `type:"leaves"`
	Sponge                     *BlockSet `type:"sponge"`
	Glass                      *BlockSet `cullAgainst:"false"`
	LapisOre                   *BlockSet
	LapisBlock                 *BlockSet
	Dispenser                  *BlockSet `type:"dispenser"`
	Sandstone                  *BlockSet
	Note                       *BlockSet
	Bed                        *BlockSet `type:"bed"`
	GoldenRail                 *BlockSet `type:"poweredRail"`
	DetectorRail               *BlockSet `type:"poweredRail"`
	StickyPiston               *BlockSet `type:"piston"`
	Web                        *BlockSet `cullAgainst:"false" collidable:"false"`
	TallGrass                  *BlockSet `type:"tallGrass" mc:"tallgrass"`
	DeadBush                   *BlockSet `type:"deadBush" mc:"deadbush"`
	Piston                     *BlockSet `type:"piston"`
	PistonHead                 *BlockSet `type:"pistonHead"`
	Wool                       *BlockSet `type:"wool"`
	PistonExtension            *BlockSet `renderable:"false"`
	YellowFlower               *BlockSet `cullAgainst:"false" collidable:"false"`
	RedFlower                  *BlockSet `cullAgainst:"false" collidable:"false"`
	BrownMushroom              *BlockSet `cullAgainst:"false" collidable:"false"`
	RedMushrrom                *BlockSet `cullAgainst:"false" collidable:"false"`
	GoldBlock                  *BlockSet
	IronBlock                  *BlockSet
	DoubleStoneSlab            *BlockSet `type:"slabDoubleSeamless" variant:"stone"`
	StoneSlab                  *BlockSet `type:"slab" variant:"stone"`
	BrickBlock                 *BlockSet
	TNT                        *BlockSet `mc:"tnt"`
	BookShelf                  *BlockSet
	MossyCobblestone           *BlockSet
	Obsidian                   *BlockSet
	Torch                      *BlockSet `type:"torch" model:"torch"`
	Fire                       *BlockSet
	MobSpawner                 *BlockSet
	OakStairs                  *BlockSet `type:"stairs"`
	Chest                      *BlockSet
	RedstoneWire               *BlockSet
	DiamondOre                 *BlockSet
	DiamondBlock               *BlockSet
	CraftingTable              *BlockSet
	Wheat                      *BlockSet
	Farmland                   *BlockSet
	Furnace                    *BlockSet
	FurnaceLit                 *BlockSet
	StandingSign               *BlockSet `type:"floorSign"`
	WoodenDoor                 *BlockSet `type:"door"`
	Ladder                     *BlockSet
	Rail                       *BlockSet `type:"rail"`
	StoneStairs                *BlockSet `type:"stairs"`
	WallSign                   *BlockSet `type:"wallSign"`
	Lever                      *BlockSet
	StonePressurePlate         *BlockSet
	IronDoor                   *BlockSet `type:"door"`
	WoodenPressurePlate        *BlockSet
	RedstoneOre                *BlockSet
	RedstoneOreLit             *BlockSet
	RedstoneTorchUnlit         *BlockSet `type:"torch" model:"unlit_redstone_torch"`
	RedstoneTorch              *BlockSet `type:"torch" model:"redstone_torch"`
	StoneButton                *BlockSet
	SnowLayer                  *BlockSet
	Ice                        *BlockSet
	Snow                       *BlockSet
	Cactus                     *BlockSet `cullAgainst:"false"`
	Clay                       *BlockSet
	Reeds                      *BlockSet `cullAgainst:"false" collidable:"false"`
	Jukebox                    *BlockSet
	Fence                      *BlockSet `type:"fence"`
	Pumpkin                    *BlockSet
	Netherrack                 *BlockSet
	SoulSand                   *BlockSet
	Glowstone                  *BlockSet
	Portal                     *BlockSet `type:"portal"`
	PumpkinLit                 *BlockSet
	Cake                       *BlockSet
	RepeaterUnpowered          *BlockSet
	RepeaterPowered            *BlockSet
	StainedGlass               *BlockSet `type:"stainedGlass"`
	TrapDoor                   *BlockSet
	MonsterEgg                 *BlockSet
	StoneBrick                 *BlockSet `mc:"stonebrick" type:"stonebrick"`
	BrownMushroomBlock         *BlockSet
	RedMushroomBlock           *BlockSet
	IronBars                   *BlockSet `type:"connectable"`
	GlassPane                  *BlockSet `type:"connectable"`
	MelonBlock                 *BlockSet
	PumpkinStem                *BlockSet
	MelonStem                  *BlockSet
	Vine                       *BlockSet `type:"vines"`
	FenceGate                  *BlockSet `type:"fenceGate"`
	BrickStairs                *BlockSet `type:"stairs"`
	StoneBrickStairs           *BlockSet `type:"stairs"`
	Mycelium                   *BlockSet
	Waterlily                  *BlockSet `type:"lilypad"`
	NetherBrick                *BlockSet
	NetherBrickFence           *BlockSet `type:"fence" wood:"false"`
	NetherBrickStairs          *BlockSet `type:"stairs"`
	NetherWart                 *BlockSet
	EnchantingTable            *BlockSet
	BrewingStand               *BlockSet
	Cauldron                   *BlockSet
	EndPortal                  *BlockSet
	EndPortalFrame             *BlockSet
	EndStone                   *BlockSet
	DragonEgg                  *BlockSet
	RedstoneLamp               *BlockSet
	RedstoneLampLit            *BlockSet
	DoubleWoodenSlab           *BlockSet `type:"slabDouble" variant:"wood"`
	WoodenSlab                 *BlockSet `type:"slab" variant:"wood"`
	Cocoa                      *BlockSet
	SandstoneStairs            *BlockSet `type:"stairs"`
	EmeraldOre                 *BlockSet
	EnderChest                 *BlockSet
	TripwireHook               *BlockSet
	Tripwire                   *BlockSet
	EmeraldBlock               *BlockSet
	SpruceStairs               *BlockSet `type:"stairs"`
	BirchStairs                *BlockSet `type:"stairs"`
	JungleStairs               *BlockSet `type:"stairs"`
	CommandBlock               *BlockSet
	Beacon                     *BlockSet `cullAgainst:"false"`
	CobblestoneWall            *BlockSet `type:"wall"`
	FlowerPot                  *BlockSet
	Carrots                    *BlockSet
	Potatoes                   *BlockSet
	WoodenButton               *BlockSet
	Skull                      *BlockSet `type:"skull"`
	Anvil                      *BlockSet
	TrappedChest               *BlockSet
	LightWeightedPressurePlate *BlockSet
	HeavyWeightedPressurePlate *BlockSet
	ComparatorUnpowered        *BlockSet
	ComparatorPowered          *BlockSet
	DaylightDetector           *BlockSet
	RedstoneBlock              *BlockSet
	QuartzOre                  *BlockSet
	Hopper                     *BlockSet
	QuartzBlock                *BlockSet
	QuartzStairs               *BlockSet `type:"stairs"`
	ActivatorRail              *BlockSet `type:"poweredRail"`
	Dropper                    *BlockSet `type:"dispenser"`
	StainedHardenedClay        *BlockSet `type:"stainedClay"`
	StainedGlassPane           *BlockSet `type:"stainedGlassPane"`
	Leaves2                    *BlockSet `type:"leaves" second:"true"`
	Log2                       *BlockSet `type:"log" second:"true"`
	AcaciaStairs               *BlockSet `type:"stairs"`
	DarkOakStairs              *BlockSet `type:"stairs"`
	Slime                      *BlockSet
	Barrier                    *BlockSet `cullAgainst:"false" renderable:"false"`
	IronTrapDoor               *BlockSet
	Prismarine                 *BlockSet
	SeaLantern                 *BlockSet
	HayBlock                   *BlockSet
	Carpet                     *BlockSet `type:"carpet"`
	HardenedClay               *BlockSet
	CoalBlock                  *BlockSet
	PackedIce                  *BlockSet
	DoublePlant                *BlockSet
	StandingBanner             *BlockSet
	WallBanner                 *BlockSet
	DaylightDetectorInverted   *BlockSet
	RedSandstone               *BlockSet
	RedSandstoneStairs         *BlockSet `type:"stairs"`
	DoubleStoneSlab2           *BlockSet `type:"slabDoubleSeamless" variant:"stone2"`
	StoneSlab2                 *BlockSet `type:"slab" variant:"stone2"`
	SpruceFenceGate            *BlockSet `type:"fenceGate"`
	BirchFenceGate             *BlockSet `type:"fenceGate"`
	JungleFenceGate            *BlockSet `type:"fenceGate"`
	DarkOakFenceGate           *BlockSet `type:"fenceGate"`
	AcaciaFenceGate            *BlockSet `type:"fenceGate"`
	SpruceFence                *BlockSet `type:"fence"`
	BirchFence                 *BlockSet `type:"fence"`
	JungleFence                *BlockSet `type:"fence"`
	DarkOakFence               *BlockSet `type:"fence"`
	AcaciaFence                *BlockSet `type:"fence"`
	SpruceDoor                 *BlockSet `type:"door"`
	BirchDoor                  *BlockSet `type:"door"`
	JungleDoor                 *BlockSet `type:"door"`
	AcaciaDoor                 *BlockSet `type:"door"`
	DarkOakDoor                *BlockSet `type:"door"`

	MissingBlock *BlockSet `mc:"steven:missing_block"`
}{}

var blockTypes = map[string]reflect.Type{}

func init() {
	type loadable interface {
		load(tag reflect.StructTag)
	}
	v := reflect.ValueOf(&Blocks).Elem()
	t := v.Type()
	bsType := reflect.TypeOf(&BlockSet{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)
		if !f.Type.AssignableTo(bsType) {
			continue
		}
		tag := f.Tag

		ty := tag.Get("type")
		if ty == "" {
			ty = "default"
		}

		name := tag.Get("mc")
		if name == "" {
			name = formatFieldName(f.Name)
		}

		rT, ok := blockTypes[ty]
		if !ok {
			panic("invalid block type " + ty)
		}
		nv := reflect.New(rT)
		block := nv.Interface().(Block)
		block.init(name)
		if l, ok := block.(loadable); ok {
			l.load(tag)
		}
		set := alloc(block)
		fv.Set(reflect.ValueOf(set))
	}
}

func formatFieldName(name string) string {
	var buf bytes.Buffer
	for _, r := range name {
		if unicode.IsUpper(r) {
			r = unicode.ToLower(r)
			if buf.Len() > 0 {
				buf.WriteRune('_')
			}
		}
		buf.WriteRune(r)
	}
	return buf.String()
}

func registerBlockType(name string, v Block) {
	blockTypes[name] = reflect.TypeOf(v).Elem()
}

func wrapTagBool(tag reflect.StructTag) func(name string, def bool) bool {
	return func(name string, def bool) bool {
		v := tag.Get(name)
		if v == "" {
			return def
		}
		return v == "true"
	}
}
