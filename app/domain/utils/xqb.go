package utils

import "strings"

// fields that should be sanitized if found
var sensitiveKeys = []string{
	"password",
	"secret",
	"api_key",
	"token",
	"access_token",
	"refresh_token",
}

// sanitize replaces sensitive values with a mask
func sanitize(key string, value any) any {
	lower := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lower, s) {
			return ""
		}
	}
	return value
}

// NormalizeRows converts []map[string]any with dot keys into nested maps
// and sanitizes sensitive fields automatically.
func NormalizeRows(rows ...map[string]any) []map[string]any {
	out := make([]map[string]any, len(rows))

	for i, row := range rows {
		root := map[string]any{}

		for k, v := range row {
			parts := strings.Split(k, ".")
			current := root

			for j, part := range parts {
				if j == len(parts)-1 {
					// last part, set value with sanitization
					current[part] = sanitize(part, v)
				} else {
					// ensure nested map exists
					if _, ok := current[part]; !ok {
						current[part] = map[string]any{}
					}
					current = current[part].(map[string]any)
				}
			}
		}

		out[i] = root
	}

	return out
}
