package rules

import (
	"context"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/x"
)

func ExistsInDB(c context.Context, fl validator.FieldLevel) bool {
	log := x.Logger()

	param := fl.Param() // expected format: "table-column"
	parts := strings.Split(param, "-")
	if len(parts) != 2 {
		log.Error("Invalid unique validator format, expected: table-column")
		return false
	}

	tableName := strings.TrimSpace(parts[0])
	columnName := strings.TrimSpace(parts[1])
	if tableName == "" || columnName == "" {
		return false
	}

	count, err := xqb.Table(tableName).
		WithContext(c).
		Where(columnName, "=", fl.Field().Interface()).
		Count(columnName)
	if err != nil {
		log.Error("Error while checking unique constraint: " + err.Error())
		return false
	}

	return count == 1
}

func ExistsNonSensitiveInDB(c context.Context, fl validator.FieldLevel) bool {
	log := x.Logger()

	param := fl.Param() // expected format: "table-column"
	parts := strings.Split(param, "-")
	if len(parts) != 2 {
		log.Error("Invalid unique validator format, expected: table-column")
		return false
	}

	tableName := strings.TrimSpace(parts[0])
	columnName := strings.TrimSpace(parts[1])
	if tableName == "" || columnName == "" {
		return false
	}

	count, err := xqb.Table(tableName).
		WithContext(c).
		Where(columnName, "=", fl.Field().Interface()).
		Count(columnName)
	if err != nil {
		log.Error("Error while checking unique constraint: " + err.Error())
		return false
	}

	return count == 1
}
