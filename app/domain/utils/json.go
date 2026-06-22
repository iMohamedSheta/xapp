package utils

import (
	"reflect"
)

// GetJSONKeys returns the JSON keys of the given struct
func GetJSONKeys(i any) map[string]bool {
	keys := make(map[string]bool)

	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for field := range t.Fields() {
		field := field

		if field.Anonymous {
			continue
		}

		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		name := tag
		if idx := indexComma(tag); idx != -1 {
			name = tag[:idx]
		}

		keys[name] = true
	}

	return keys
}

func indexComma(tag string) int {
	for i, r := range tag {
		if r == ',' {
			return i
		}
	}
	return -1
}
