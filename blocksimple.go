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
	"reflect"
	"strconv"
)

type blockSimple struct {
	baseBlock
}

func (b *blockSimple) load(tag reflect.StructTag) {
	getBool := wrapTagBool(tag)
	b.cullAgainst = getBool("cullAgainst", true)
	b.collidable = getBool("collidable", true)
	b.renderable = getBool("renderable", true)
	b.hardness, _ = strconv.ParseFloat(tag.Get("hardness"), 64)
	b.translucent = getBool("translucent", false)
}

func (b *blockSimple) toData() int {
	if b == b.Parent.Base {
		return 0
	}
	return -1
}
