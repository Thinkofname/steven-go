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

func init() {
	registerBlockType("default", &blockSimple{})
	registerBlockType("stone", &blockStone{})
	registerBlockType("grass", &blockGrass{})
	registerBlockType("planks", &blockPlanks{})
	registerBlockType("sapling", &blockSapling{})
	registerBlockType("liquid", &blockLiquid{})
	registerBlockType("log", &blockLog{})
	registerBlockType("leaves", &blockLeaves{})
	registerBlockType("sponge", &blockSponge{})
	registerBlockType("dispenser", &blockDispenser{})
	registerBlockType("bed", &blockBed{})
	registerBlockType("rail", &blockRail{})
	registerBlockType("poweredRail", &blockPoweredRail{})
	registerBlockType("piston", &blockPiston{})
	registerBlockType("pistonHead", &blockPistonHead{})
	registerBlockType("tallGrass", &blockTallGrass{})
	registerBlockType("deadBush", &blockDeadBush{})
	registerBlockType("wool", &blockWool{})
	registerBlockType("stairs", &blockStairs{})
	registerBlockType("door", &blockDoor{})
	registerBlockType("fence", &blockFence{})
	registerBlockType("fenceGate", &blockFenceGate{})
	registerBlockType("stainedGlass", &blockStainedGlass{})
	registerBlockType("stainedGlassPane", &blockStainedGlassPane{})
	registerBlockType("stainedClay", &blockStainedClay{})
	registerBlockType("connectable", &blockConnectable{})
	registerBlockType("vines", &blockVines{})
	registerBlockType("wall", &blockWall{})
	registerBlockType("slab", &blockSlab{})
	registerBlockType("slabDouble", &blockSlabDouble{})
	registerBlockType("slabDoubleSeamless", &blockSlabDoubleSeamless{})
	registerBlockType("carpet", &blockCarpet{})
	registerBlockType("torch", &blockTorch{})
	registerBlockType("wallSign", &blockWallSign{})
	registerBlockType("floorSign", &blockFloorSign{})
	registerBlockType("skull", &blockSkull{})
	registerBlockType("portal", &blockPortal{})
	registerBlockType("lilypad", &blockLilypad{})
	registerBlockType("stonebrick", &blockStoneBrick{})
}
