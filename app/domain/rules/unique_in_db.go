package rules

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/x"
)

// Custom rule to validate simple unique columns
// Format: unique_db=table-column[-scopeColumn]
func UniqueInDB(c context.Context, fl validator.FieldLevel) bool {
	log := x.Logger()

	param := fl.Param() // e.g. "pages-slug-workspace_id"
	parts := strings.Split(param, "-")
	if len(parts) < 2 || len(parts) > 3 {
		log.Error("Invalid unique validator format, expected: table-column[-scopeColumn]")
		return false
	}

	tableName := strings.TrimSpace(parts[0])
	columnName := strings.TrimSpace(parts[1])
	var scopeColumn string
	if len(parts) == 3 {
		scopeColumn = strings.TrimSpace(parts[2])
	}

	if tableName == "" || columnName == "" {
		return false
	}

	query := xqb.Table(tableName).WithContext(c).
		Where(xqb.Lower(columnName, ""), "=", strings.ToLower(fl.Field().String()))

		// If scopeColumn is defined, get its value from the struct
	if scopeColumn != "" {
		parent := fl.Parent()
		parentType := parent.Type()

		var scopeField reflect.Value
		for i := 0; i < parent.NumField(); i++ {
			sf := parentType.Field(i)
			tag := sf.Tag.Get("json")
			if tag == scopeColumn {
				scopeField = parent.Field(i)
				break
			}
		}

		if scopeField.IsValid() {
			query = query.Where(scopeColumn, "=", scopeField.Interface())
		} else {
			log.Error("Scope column not found in struct: " + scopeColumn)
			return false
		}
	}

	count, err := query.Count("id")
	if err != nil {
		log.Error("Error while checking unique constraint: " + err.Error())
		return false
	}

	return count == 0
}
