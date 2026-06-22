package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// Convert value type (array of strings) to CSV string
func ToCSV(val any, defaultValue string) string {
	switch v := val.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ", ")
	default:
		return defaultValue
	}
}

func ToArrayOfStrings(val any, defaultValue []string) []string {
	switch v := val.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	default:
		return defaultValue
	}
}

func ToString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func ToInt64(v any) int64 {
	switch val := v.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	case float32:
		return int64(val)
	case uint:
		return int64(val)
	case uint64:
		return int64(val)
	}
	return 0
}

func ToFloat64(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	}
	return 0
}
func ToBool(v any) bool {
	if v == nil {
		return false
	}
	switch b := v.(type) {
	case bool:
		return b
	case string:
		return b == "true" || b == "1"
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%v", b) == "1"
	case float32, float64:
		return fmt.Sprintf("%v", b) == "1"
	default:
		return false
	}
}

// ParsePGArrayToInt64Slice parses a PostgreSQL array_agg string like "{1,2,3}" into []int64
func ParsePGArrayToInt64Slice(val any) []int64 {
	if val == nil {
		return nil
	}

	str, ok := val.(string)
	if !ok {
		return nil
	}

	str = strings.TrimPrefix(str, "{")
	str = strings.TrimSuffix(str, "}")

	if str == "" {
		return nil
	}

	parts := strings.Split(str, ",")
	result := make([]int64, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err != nil {
			continue
		}
		result = append(result, n)
	}
	return result
}
