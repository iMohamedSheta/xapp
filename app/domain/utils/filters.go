package utils

import (
	"strings"

	"github.com/iMohamedSheta/xqb"
)

type FilterOptions struct {
	// Search    string
	SearchBy  map[string]string
	SortBy    string
	SortOrder string
}

func ApplyFilters(
	q *xqb.QueryBuilder,
	filters *FilterOptions,
	mappedFields map[string][]string,
	defaults *FilterOptions,
) {
	filters = mergeDefaults(filters, defaults)

	// --- SearchBy (field specific) ---
	if filters.SearchBy != nil {
		for key, value := range filters.SearchBy {
			if fieldsToSearch, ok := mappedFields[key]; ok {
				q.WhereGroup(func(qb *xqb.QueryBuilder) {
					for i, field := range fieldsToSearch {
						cond := "CAST(" + field + " AS TEXT)"
						if i == 0 {
							qb.Where(cond, "LIKE", "%"+value+"%")
						} else {
							qb.OrWhere(cond, "LIKE", "%"+value+"%")
						}
					}
				})
			}
		}
	}

	// --- Global search ---
	// if filters.Search != "" {
	// 	q.WhereGroup(func(qb *xqb.QueryBuilder) {
	// 		for _, fields := range mappedFields {
	// 			for _, field := range fields {
	// 				qb.OrWhere("CAST("+field+" AS TEXT)", "ILIKE", "%"+filters.Search+"%")
	// 			}
	// 		}
	// 	})
	// }

	// --- Sorting ---
	if filters.SortBy != "" {
		if fieldsToSort, ok := mappedFields[filters.SortBy]; ok {
			order := strings.ToUpper(filters.SortOrder)
			if order != "ASC" && order != "DESC" {
				order = "ASC"
			}
			q.OrderBy(fieldsToSort[0], order)
		}
	}
}

// mergeDefaults merges missing fields from defaults into filters
func mergeDefaults(filters, defaults *FilterOptions) *FilterOptions {
	if filters == nil {
		return defaults
	}
	if filters.SearchBy == nil {
		filters.SearchBy = defaults.SearchBy
	}
	if filters.SortBy == "" {
		filters.SortBy = defaults.SortBy
	}
	if filters.SortOrder == "" {
		filters.SortOrder = defaults.SortOrder
	}
	return filters
}
