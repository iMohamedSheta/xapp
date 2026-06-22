package utils

type RequestHasRequestSentFields interface {
	GetRequestSentFields() map[string]any
}

type ExtractedFields map[string]any

// ExtractUpdatedFields filters the sent fields and keeps only valid JSON keys.
func ExtractUpdatedFields(req RequestHasRequestSentFields, excludeKeys ...string) ExtractedFields {
	sentFields := req.GetRequestSentFields()
	validKeys := GetJSONKeys(req)

	exclusions := make(map[string]bool, len(excludeKeys))
	for _, key := range excludeKeys {
		exclusions[key] = true
	}

	updated := make(map[string]any)
	for key, value := range sentFields {
		if validKeys[key] && !exclusions[key] {
			updated[key] = value
		}
	}

	return updated
}

func (f ExtractedFields) Exists(key string) bool {
	return (f)[key] != nil
}

func (f ExtractedFields) Get(key string) any {
	if f.Exists(key) {
		return (f)[key]
	}
	return nil
}

func (f ExtractedFields) GetMap(key string) map[string]any {
	if val := f.Get(key); val != nil {
		if m, ok := val.(map[string]any); ok {
			return m
		}
	}
	return nil
}

// GetMapFromRequest extracts a nested map directly from the request
func GetMapFromRequest(req RequestHasRequestSentFields, key string) map[string]any {
	if val := req.GetRequestSentFields()[key]; val != nil {
		if m, ok := val.(map[string]any); ok {
			return m
		}
	}
	return nil
}

func PrefixErrorKeys(errors map[string]any, prefix string) map[string]any {
	prefixed := make(map[string]any)
	for key, value := range errors {
		prefixed[prefix+"."+key] = value
	}
	return prefixed
}
