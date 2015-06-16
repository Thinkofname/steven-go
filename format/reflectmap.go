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

package format

import (
	"encoding/json"
	"reflect"
	"strings"
)

func mapStruct(val interface{}, m map[string]json.RawMessage) error {
	v := reflect.ValueOf(val).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type.Kind() == reflect.Struct {
			mapStruct(v.Field(i).Addr().Interface(), m)
			continue
		}
		name := f.Tag.Get("json")
		if name == "" {
			name = f.Name
		}
		if strings.ContainsRune(name, ',') {
			name = name[:strings.IndexRune(name, ',')]
		}
		mv, ok := m[name]
		if !ok {
			continue
		}
		json.Unmarshal([]byte(mv), v.Field(i).Addr().Interface())
	}
	return nil
}
