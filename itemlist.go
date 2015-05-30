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

var itemsByID = map[int]func() ItemType{
	256: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.shovelIron.name"
		i.itemNamed.name = "iron_shovel"
		return i
	},
	257: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.pickaxeIron.name"
		i.itemNamed.name = "iron_pickaxe"
		return i
	},
	258: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.axeIron.name"
		i.itemNamed.name = "iron_axe"
		return i
	},
	259: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.flintAndSteel.name"
		i.itemNamed.name = "flint_and_steel"
		return i
	},
	260: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.apple.name"
		i.itemNamed.name = "apple"
		return i
	},
	261: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bow.name"
		i.itemNamed.name = "bow"
		return i
	},
	262: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.arrow.name"
		i.itemNamed.name = "arrow"
		return i
	},
	263: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.coal.name"
		i.itemNamed.name = "coal"
		return i
	},
	264: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.diamond.name"
		i.itemNamed.name = "diamond"
		return i
	},
	265: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.ingotIron.name"
		i.itemNamed.name = "iron_ingot"
		return i
	},
	266: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.ingotGold.name"
		i.itemNamed.name = "gold_ingot"
		return i
	},
	267: func() ItemType {
		i := &itemSword{}
		i.locale = "item.swordIron.name"
		i.itemNamed.name = "iron_sword"
		i.maxDamage = 251
		return i
	},
	268: func() ItemType {
		i := &itemSword{}
		i.locale = "item.swordWood.name"
		i.itemNamed.name = "wooden_sword"
		i.maxDamage = 260
		return i
	},
	269: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.shovelWood.name"
		i.itemNamed.name = "wooden_shovel"
		return i
	},
	270: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.pickaxeWood.name"
		i.itemNamed.name = "wooden_pickaxe"
		return i
	},
	271: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.axeWood.name"
		i.itemNamed.name = "wooden_axe"
		return i
	},
	272: func() ItemType {
		i := &itemSword{}
		i.locale = "item.swordStone.name"
		i.itemNamed.name = "stone_sword"
		i.maxDamage = 132
		return i
	},
	273: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.shovelStone.name"
		i.itemNamed.name = "stone_shovel"
		return i
	},
	274: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.pickaxeStone.name"
		i.itemNamed.name = "stone_pickaxe"
		return i
	},
	275: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.axeStone.name"
		i.itemNamed.name = "stone_axe"
		return i
	},
	276: func() ItemType {
		i := &itemSword{}
		i.locale = "item.swordDiamond.name"
		i.itemNamed.name = "diamond_sword"
		i.maxDamage = 1562
		return i
	},
	277: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.shovelDiamond.name"
		i.itemNamed.name = "diamond_shovel"
		return i
	},
	278: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.pickaxeDiamond.name"
		i.itemNamed.name = "diamond_pickaxe"
		return i
	},
	279: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.axeDiamond.name"
		i.itemNamed.name = "diamond_axe"
		return i
	},
	280: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.stick.name"
		i.itemNamed.name = "stick"
		return i
	},
	281: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bowl.name"
		i.itemNamed.name = "bowl"
		return i
	},
	282: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.mushrromStew.name"
		i.itemNamed.name = "mushroom_stew"
		return i
	},
	283: func() ItemType {
		i := &itemSword{}
		i.locale = "item.swordGold.name"
		i.itemNamed.name = "golden_sword"
		i.maxDamage = 33
		return i
	},
	284: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.shovelGold.name"
		i.itemNamed.name = "golden_sword"
		return i
	},
	285: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.pickaxeGold.name"
		i.itemNamed.name = "golden_pickaxe"
		return i
	},
	286: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.axeGold.name"
		i.itemNamed.name = "golden_axe"
		return i
	},
	287: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.string.name"
		i.itemNamed.name = "string"
		return i
	},
	288: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.feather.name"
		i.itemNamed.name = "feather"
		return i
	},
	289: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.sulpher.name"
		i.itemNamed.name = "gunpowder"
		return i
	},
	290: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.hoeWood.name"
		i.itemNamed.name = "wooden_hoe"
		return i
	},
	291: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.hoeStone.name"
		i.itemNamed.name = "stone_hoe"
		return i
	},
	292: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.hoeIron.name"
		i.itemNamed.name = "iron_hoe"
		return i
	},
	293: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.hoeDiamond.name"
		i.itemNamed.name = "diamond_hoe"
		return i
	},
	294: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.hoeGold.name"
		i.itemNamed.name = "golden_hoe"
		return i
	},
	295: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.seeds.name"
		i.itemNamed.name = "wheat_seeds"
		return i
	},
	296: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.wheat.name"
		i.itemNamed.name = "wheat"
		return i
	},
	297: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bread.name"
		i.itemNamed.name = "bread"
		return i
	},
	298: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.helmetCloth.name"
		i.itemNamed.name = "leather_helmet"
		return i
	},
	299: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.chestplateCloth.name"
		i.itemNamed.name = "leather_chestplate"
		return i
	},
	300: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.leggingsCloth.name"
		i.itemNamed.name = "leather_leggings"
		return i
	},
	301: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bootsCloth.name"
		i.itemNamed.name = "leather_boots"
		return i
	},
	302: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.helmetChain.name"
		i.itemNamed.name = "chainmail_helmet"
		return i
	},
	303: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.chestplateChain.name"
		i.itemNamed.name = "chainmail_chestplate"
		return i
	},
	304: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.leggingsChain.name"
		i.itemNamed.name = "chainmail_leggings"
		return i
	},
	305: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bootsChain.name"
		i.itemNamed.name = "chainmail_boots"
		return i
	},
	306: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.helmetIron.name"
		i.itemNamed.name = "iron_helmet"
		return i
	},
	307: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.chestplateIron.name"
		i.itemNamed.name = "iron_chestplate"
		return i
	},
	308: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.leggingsIron.name"
		i.itemNamed.name = "iron_leggings"
		return i
	},
	309: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bootsIron.name"
		i.itemNamed.name = "iron_boots"
		return i
	},
	310: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.helmetDiamond.name"
		i.itemNamed.name = "diamond_helmet"
		return i
	},
	311: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.chestplateDiamond.name"
		i.itemNamed.name = "diamond_chestplate"
		return i
	},
	312: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.leggingsDiamond.name"
		i.itemNamed.name = "diamond_leggings"
		return i
	},
	313: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bootsDiamond.name"
		i.itemNamed.name = "diamond_boots"
		return i
	},
	314: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.helmetGold.name"
		i.itemNamed.name = "golden_helmet"
		return i
	},
	315: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.chestplateGold.name"
		i.itemNamed.name = "golden_chestplate"
		return i
	},
	316: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.leggingsGold.name"
		i.itemNamed.name = "golden_leggings"
		return i
	},
	317: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bootsGold.name"
		i.itemNamed.name = "golden_boots"
		return i
	},
	318: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.flint.name"
		i.itemNamed.name = "flint"
		return i
	},
	319: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.porkchopRaw.name"
		i.itemNamed.name = "porkchop"
		return i
	},
	320: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.porkchopCooked.name"
		i.itemNamed.name = "cooked_porkchop"
		return i
	},
	321: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.painting.name"
		i.itemNamed.name = "painting"
		return i
	},
	322: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.appleGold.name"
		i.itemNamed.name = "golden_apple"
		return i
	},
	323: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.sign.name"
		i.itemNamed.name = "sign"
		return i
	},
	324: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.doorOak.name"
		i.itemNamed.name = "wooden_door"
		return i
	},
	325: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bucket.name"
		i.itemNamed.name = "bucket"
		return i
	},
	326: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bucketWater.name"
		i.itemNamed.name = "water_bucket"
		return i
	},
	327: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bucketLava.name"
		i.itemNamed.name = "lava_bucket"
		return i
	},
	328: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.minecart.name"
		i.itemNamed.name = "minecart"
		return i
	},
	329: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.saddle.name"
		i.itemNamed.name = "saddle"
		return i
	},
	330: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.doorIron.name"
		i.itemNamed.name = "iron_door"
		return i
	},
	331: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.redstone.name"
		i.itemNamed.name = "redstone"
		return i
	},
	332: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.snowball.name"
		i.itemNamed.name = "snowball"
		return i
	},
	333: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.boat.name"
		i.itemNamed.name = "boat"
		return i
	},
	334: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.leather.name"
		i.itemNamed.name = "leather"
		return i
	},
	335: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.milk.name"
		i.itemNamed.name = "milk_bucket"
		return i
	},
	336: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.brick.name"
		i.itemNamed.name = "brick"
		return i
	},
	337: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.clay.name"
		i.itemNamed.name = "clay_ball"
		return i
	},
	338: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.reeds.name"
		i.itemNamed.name = "reeds"
		return i
	},
	339: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.paper.name"
		i.itemNamed.name = "paper"
		return i
	},
	340: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.book.name"
		i.itemNamed.name = "book"
		return i
	},
	341: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.slimeball.name"
		i.itemNamed.name = "slime_ball"
		return i
	},
	342: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.minecartChest.name"
		i.itemNamed.name = "chest_minecart"
		return i
	},
	343: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.minecartFurnace.name"
		i.itemNamed.name = "furnace_minecart"
		return i
	},
	344: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.egg.name"
		i.itemNamed.name = "egg"
		return i
	},
	345: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.compass.name"
		i.itemNamed.name = "compass"
		return i
	},
	346: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.fishingRod.name"
		i.itemNamed.name = "fishing_rod"
		return i
	},
	347: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.clock.name"
		i.itemNamed.name = "clock"
		return i
	},
	348: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.yellowDust.name"
		i.itemNamed.name = "glowstone_dust"
		return i
	},
	349: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.fish.name"
		i.itemNamed.name = "fish"
		return i
	},
	350: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.fish.name"
		i.itemNamed.name = "cooked_fish"
		return i
	},
	351: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.dyePowder.name"
		i.itemNamed.name = "dye"
		return i
	},
	352: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bone.name"
		i.itemNamed.name = "bone"
		return i
	},
	353: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.sugar.name"
		i.itemNamed.name = "sugar"
		return i
	},
	354: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.cake.name"
		i.itemNamed.name = "cake"
		return i
	},
	355: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.bed.name"
		i.itemNamed.name = "bed"
		return i
	},
	356: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.diode.name"
		i.itemNamed.name = "repeater"
		return i
	},
	357: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.cookie.name"
		i.itemNamed.name = "cookie"
		return i
	},
	358: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.map.name"
		i.itemNamed.name = "filled_map"
		return i
	},
	359: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.shears.name"
		i.itemNamed.name = "shears"
		return i
	},
	360: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.melon.name"
		i.itemNamed.name = "melon"
		return i
	},
	361: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.seeds_pumpkin.name"
		i.itemNamed.name = "pumpkin_seeds"
		return i
	},
	362: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.seeds_melon.name"
		i.itemNamed.name = "melon_seeds"
		return i
	},
	363: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.beefRaw.name"
		i.itemNamed.name = "beef"
		return i
	},
	364: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.beefCooked.name"
		i.itemNamed.name = "cooked_beef"
		return i
	},
	365: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.chickenRaw.name"
		i.itemNamed.name = "chicken"
		return i
	},
	366: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.chickenCooked.name"
		i.itemNamed.name = "cooked_chicken"
		return i
	},
	367: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.rottenFlesh.name"
		i.itemNamed.name = "rotten_flesh"
		return i
	},
	368: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.enderPearl.name"
		i.itemNamed.name = "ender_pearl"
		return i
	},
	369: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.blazeRod.name"
		i.itemNamed.name = "blaze_rod"
		return i
	},
	370: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.ghastTear.name"
		i.itemNamed.name = "ghast_tear"
		return i
	},
	371: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.goldNugget.name"
		i.itemNamed.name = "gold_nugget"
		return i
	},
	372: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.netherStalkSeeds.name"
		i.itemNamed.name = "nether_wart"
		return i
	},
	373: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.potion.name"
		i.itemNamed.name = "potion"
		return i
	},
	374: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.glassBottle.name"
		i.itemNamed.name = "glass_bottle"
		return i
	},
	375: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.spiderEye.name"
		i.itemNamed.name = "spider_eye"
		return i
	},
	376: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.fermentedSpiderEye.name"
		i.itemNamed.name = "fermented_spider_eye"
		return i
	},
	377: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.blazePowder.name"
		i.itemNamed.name = "blaze_powder"
		return i
	},
	378: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.magmaCream.name"
		i.itemNamed.name = "magma_cream"
		return i
	},
	379: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.brewingStand.name"
		i.itemNamed.name = "brewing_stand"
		return i
	},
	380: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.cauldron.name"
		i.itemNamed.name = "cauldron"
		return i
	},
	381: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.eyeOfEnder.name"
		i.itemNamed.name = "ender_eye"
		return i
	},
	382: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.speckledMelon.name"
		i.itemNamed.name = "speckled_melon"
		return i
	},
	383: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.monsterPlacer.name"
		i.itemNamed.name = "spawn_egg"
		return i
	},
	384: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.expBottle.name"
		i.itemNamed.name = "experience_bottle"
		return i
	},
	385: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.fireball.name"
		i.itemNamed.name = "fire_charge"
		return i
	},
	386: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.writingBook.name"
		i.itemNamed.name = "writable_book"
		return i
	},
	387: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.writtenBook.name"
		i.itemNamed.name = "written_book"
		return i
	},
	388: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.emerald.name"
		i.itemNamed.name = "emerald"
		return i
	},
	389: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.frame.name"
		i.itemNamed.name = "item_frame"
		return i
	},
	390: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.flowerPot.name"
		i.itemNamed.name = "flower_pot"
		return i
	},
	391: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.carrots.name"
		i.itemNamed.name = "carrot"
		return i
	},
	392: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.potato.name"
		i.itemNamed.name = "potato"
		return i
	},
	393: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.potatoBaked.name"
		i.itemNamed.name = "baked_potato"
		return i
	},
	394: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.potatoPoisonous.name"
		i.itemNamed.name = "poisonous_potato"
		return i
	},
	395: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.emptyMap.name"
		i.itemNamed.name = "map"
		return i
	},
	396: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.carrotGolden.name"
		i.itemNamed.name = "golden_carrot"
		return i
	},
	397: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.skull.name"
		i.itemNamed.name = "skull"
		return i
	},
	398: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.carrotOnAStick.name"
		i.itemNamed.name = "carrot_on_a_stick"
		return i
	},
	399: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.netherStar.name"
		i.itemNamed.name = "nether_star"
		return i
	},
	400: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.pumpkinPie.name"
		i.itemNamed.name = "pumpkin_pie"
		return i
	},
	401: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.fireworks.name"
		i.itemNamed.name = "fireworks"
		return i
	},
	402: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.fireworksCharge.name"
		i.itemNamed.name = "firework_charge"
		return i
	},
	403: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.enchantedBook.name"
		i.itemNamed.name = "enchanted_book"
		return i
	},
	404: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.comparator.name"
		i.itemNamed.name = "comparator"
		return i
	},
	405: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.netherbrick.name"
		i.itemNamed.name = "netherbrick"
		return i
	},
	406: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.netherquartz.name"
		i.itemNamed.name = "quartz"
		return i
	},
	407: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.minecartTnt.name"
		i.itemNamed.name = "tnt_minecart"
		return i
	},
	408: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.minecartHopper.name"
		i.itemNamed.name = "hopper_minecart"
		return i
	},
	409: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.prismarineShard.name"
		i.itemNamed.name = "prismarine_shard"
		return i
	},
	410: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.prismarineCrystals.name"
		i.itemNamed.name = "prismarine_crystals"
		return i
	},
	411: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.rabbitRaw.name"
		i.itemNamed.name = "rabbit"
		return i
	},
	412: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.rabbitCooked.name"
		i.itemNamed.name = "cooked_rabbit"
		return i
	},
	413: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.rabbitStew.name"
		i.itemNamed.name = "rabbit_stew"
		return i
	},
	414: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.rabbitFoot.name"
		i.itemNamed.name = "rabbit_foot"
		return i
	},
	415: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.rabbitHide.name"
		i.itemNamed.name = "rabbit_hide"
		return i
	},
	416: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.armorStand.name"
		i.itemNamed.name = "armor_stand"
		return i
	},
	417: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.horsearmormetal.name"
		i.itemNamed.name = "iron_horse_armor"
		return i
	},
	418: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.horsearmorgold.name"
		i.itemNamed.name = "golden_horse_armor"
		return i
	},
	419: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.horsearmordiamond.name"
		i.itemNamed.name = "diamond_horse_armor"
		return i
	},
	420: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.leash.name"
		i.itemNamed.name = "lead"
		return i
	},
	421: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.nameTag.name"
		i.itemNamed.name = "name_tag"
		return i
	},
	422: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.minecartCommandBlock.name"
		i.itemNamed.name = "command_block_minecart"
		return i
	},
	423: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.muttonRaw.name"
		i.itemNamed.name = "mutton"
		return i
	},
	424: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.muttonCooked.name"
		i.itemNamed.name = "cooked_mutton"
		return i
	},
	425: func() ItemType {
		i := &itemBasic{}
		i.locale = "tile.banner.name"
		i.itemNamed.name = "banner"
		return i
	},
	427: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.doorSpruce.name"
		i.itemNamed.name = "spruce_door"
		return i
	},
	428: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.doorBirch.name"
		i.itemNamed.name = "birch_door"
		return i
	},
	429: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.doorJungle.name"
		i.itemNamed.name = "jungle_door"
		return i
	},
	430: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.doorAcacia.name"
		i.itemNamed.name = "acacia_door"
		return i
	},
	431: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.doorDarkOak.name"
		i.itemNamed.name = "dark_oak_door"
		return i
	},
	2256: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_13"
		return i
	},
	2257: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_cat"
		return i
	},
	2258: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_blocks"
		return i
	},
	2259: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_chirp"
		return i
	},
	2260: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_far"
		return i
	},
	2261: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_mall"
		return i
	},
	2262: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_mellohi"
		return i
	},
	2263: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_stal"
		return i
	},
	2264: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_strad"
		return i
	},
	2265: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_ward"
		return i
	},
	2266: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_11"
		return i
	},
	2267: func() ItemType {
		i := &itemBasic{}
		i.locale = "item.record.name"
		i.itemNamed.name = "record_wait"
		return i
	},
}

type itemBasic struct {
	displayTag
	itemSimpleLocale
	itemNamed
}

func (i *itemBasic) Stackable() bool {
	return true
}

func (i *itemBasic) ParseDamage(d int16) {}

type itemSword struct {
	displayTag
	itemDamagable
	itemSimpleLocale
	itemNamed
}
