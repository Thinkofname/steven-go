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

func (ce *clientEntities) registerBlockEntities() {

}

type BlockEntity interface {
	BlockComponent
}

type blockComponent struct {
	Location Position
}

func (bc *blockComponent) Position() Position {
	return bc.Location
}

func (bc *blockComponent) SetPosition(p Position) {
	bc.Location = p
}

type BlockComponent interface {
	Position() Position
	SetPosition(p Position)
}
