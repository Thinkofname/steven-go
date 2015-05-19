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
	272: func() ItemType {
		i := &itemSword{}
		i.locale = "item.swordStone.name"
		i.itemNamed.name = "stone_sword"
		i.maxDamage = 132
		return i
	},
	276: func() ItemType {
		i := &itemSword{}
		i.locale = "item.swordDiamond.name"
		i.itemNamed.name = "diamond_sword"
		i.maxDamage = 1562
		return i
	},

	283: func() ItemType {
		i := &itemSword{}
		i.locale = "item.swordGold.name"
		i.itemNamed.name = "golden_sword"
		i.maxDamage = 33
		return i
	},
}

type itemSword struct {
	displayTag
	itemDamagable
	itemSimpleLocale
	itemNamed
}

func (i *itemSword) Stackable() bool {
	return false
}
