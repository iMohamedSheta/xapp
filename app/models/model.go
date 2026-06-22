package models

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/iMohamedSheta/xqb"
)

type BaseModel struct{}

// GenerateColumns builds a list of table columns with table_column aliases, ignoring relations
func (BaseModel) Cols(model xqb.ModelInterface, aliasTable string, ignoredColumns ...string) []any {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if aliasTable == "" {
		aliasTable = model.Table()
	}

	var columns []any

	for field := range t.Fields() {
		xqbTag := field.Tag.Get("xqb")
		if xqbTag == "" {
			continue
		}

		kind := field.Type.Kind()

		// skip slices
		if kind == reflect.Slice {
			continue
		}

		// skip nested structs that implement ModelInterface (1:1 or nested relations)
		if kind == reflect.Struct {
			newVal := reflect.New(field.Type).Interface()
			if _, ok := newVal.(xqb.ModelInterface); ok {
				continue
			}
		}

		if slices.Contains(ignoredColumns, xqbTag) {
			continue
		}

		// add normal column
		columns = append(columns, fmt.Sprintf("%s.%s AS %s.%s", model.Table(), xqbTag, aliasTable, xqbTag))
	}

	return columns
}
