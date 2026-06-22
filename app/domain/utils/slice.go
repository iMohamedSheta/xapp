package utils

import (
	"fmt"
	"sort"
	"strings"
)

func UniqueSliceUInts(input []uint) []uint {
	seen := make(map[uint]bool)
	result := make([]uint, 0)

	for _, val := range input {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}

	return result
}

func ToAnySlice[T any](in []T) []any {
	out := make([]any, len(in))
	for i, v := range in {
		out[i] = v
	}
	return out
}

func ToStringSlice[T any](in []T) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = fmt.Sprint(v)
	}
	return out
}

func SortMaps(maps []map[string]any, key string, order string) {
	sort.Slice(maps, func(i, j int) bool {
		valI := ToInt64(maps[i][key])
		valJ := ToInt64(maps[j][key])
		if strings.ToUpper(order) == "DESC" {
			return valI > valJ
		}
		return valI < valJ
	})
}
