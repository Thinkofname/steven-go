package chat

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
