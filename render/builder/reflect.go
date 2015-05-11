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

package builder

import "reflect"

// Struct returns a function that will serialize
// structs of the type passed to Struct originally.
// It also returns an array of types that can be
// passed to New
func Struct(i interface{}) (func(*Buffer, interface{}), []Type) {
	t := reflect.TypeOf(i)
	isPtr := false
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		isPtr = true
	}
	l := t.NumField()

	var funcs []func(buf *Buffer, v reflect.Value)
	var types []Type

	for j := 0; j < l; j++ {
		f := t.Field(j)

		var fu func(buf *Buffer, v reflect.Value)
		var t Type

		ii := j
		switch f.Type.Kind() {
		case reflect.Float32:
			fu = func(buf *Buffer, v reflect.Value) {
				buf.Float(float32(v.Field(ii).Float()))
			}
			t = Float
		case reflect.Uint16:
			fu = func(buf *Buffer, v reflect.Value) {
				buf.UnsignedShort(uint16(v.Field(ii).Uint()))
			}
			t = UnsignedShort
		case reflect.Int16:
			fu = func(buf *Buffer, v reflect.Value) {
				buf.Short(int16(v.Field(ii).Int()))
			}
			t = Short
		case reflect.Uint8:
			fu = func(buf *Buffer, v reflect.Value) {
				buf.UnsignedByte(uint8(v.Field(ii).Uint()))
			}
			t = UnsignedByte
		case reflect.Int8:
			fu = func(buf *Buffer, v reflect.Value) {
				buf.Byte(int8(v.Field(ii).Int()))
			}
			t = Byte
		default:
			panic("unsupported type " + f.Type.String())
		}
		funcs = append(funcs, fu)
		types = append(types, t)
	}

	return func(buf *Buffer, i interface{}) {
		v := reflect.ValueOf(i)
		if isPtr {
			v = v.Elem()
		}

		for _, f := range funcs {
			f(buf, v)
		}
	}, types
}
