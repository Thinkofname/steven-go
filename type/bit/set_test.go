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

package bit

import "testing"

func TestSet(t *testing.T) {
	s := NewSet(200)
	for i := 0; i < 200; i++ {
		if i%3 == 0 {
			s.Set(i, true)
		}
	}
	for i := 0; i < 200; i++ {
		if s.Get(i) != (i%3 == 0) {
			t.Fail()
		}
	}
}
