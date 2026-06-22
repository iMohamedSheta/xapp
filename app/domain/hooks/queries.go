//go:build !dev

package hooks

type QueryLog struct{}

// InitQueryHooks wires all query builder hooks
func InitQueryHooks() {}

// GetQueries returns all logged queries for a request
func GetQueries(reqID string) []QueryLog {
	return nil
}

// ClearQueries removes all logged queries for a request
func ClearQueries(reqID string) {}
