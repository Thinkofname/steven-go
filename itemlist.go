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
	276: func() ItemType {
		i := &itemSword{
			itemSimpleLocale: itemSimpleLocale{
				locale: "item.swordDiamond.name",
			},
			itemNamed: itemNamed{
				name: "diamond_sword",
			},
			itemDamagable: itemDamagable{
				maxDamage: 1562,
			},
		}
		return i
	},

	283: func() ItemType {
		i := &itemSword{
			itemSimpleLocale: itemSimpleLocale{
				locale: "item.swordGold.name",
			},
			itemNamed: itemNamed{
				name: "golden_sword",
			},
			itemDamagable: itemDamagable{
				maxDamage: 33,
			},
		}
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
