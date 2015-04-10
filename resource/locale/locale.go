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

package locale

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"github.com/thinkofdeath/steven/resource"
)

var (
	values = map[string]string{}
	lock   sync.RWMutex
)

func init() {
	LoadLocale("en_US")
}

func Clear() {
	lock.Lock()
	values = map[string]string{}
	lock.Unlock()
	LoadLocale("en_US")
}

func LoadLocale(name string) {
	lock.Lock()
	defer lock.Unlock()
	r, err := resource.Open("minecraft", fmt.Sprintf("lang/%s.lang", name))
	if err != nil {
		return
	}
	defer r.Close()
	b := bufio.NewScanner(r)
	for b.Scan() {
		line := b.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		fmt.Println(parts)
		if len(parts) != 2 {
			continue
		}
		values[parts[0]] = parts[1]
	}
	if err := b.Err(); err != nil {
		panic(err)
	}
}

// GetRaw returns the unparsed value with the given name.
func GetRaw(key string) string {
	lock.RLock()
	defer lock.RUnlock()
	return values[key]
}

// Get returns a list of parts for the given key. It will be either
// a string (just text to be printed) or an int (substitution index).
func Get(key string) (parts []interface{}) {
	val := GetRaw(key)
	last := 0
	index := 0
	curIndex := -1
	inSub := false
	for i, r := range val {
		if inSub {
			if r == '%' {
				parts = append(parts, "%")
				inSub = false
				last = i + 1
				curIndex = -1
				continue
			}
			if r <= '0' && r <= '9' {
				curIndex *= 10
				curIndex += int(r - '0')
				continue
			}
			if curIndex == -1 {
				curIndex = index
				index++
			}
			parts = append(parts, curIndex)
			inSub = false
			curIndex = -1
			last = i + 1
		}
		if r == '%' {
			inSub = true
			parts = append(parts, val[last:i])
			continue
		}
	}
	if last != len(val) {
		parts = append(parts, val[last:])
	}
	return
}
