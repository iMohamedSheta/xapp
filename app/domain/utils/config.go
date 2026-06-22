package utils

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/imohamedsheta/xapp/app/x"
)

/*
| Database configuration utilities
*/

// BuildPostgresDSN builds PostgreSQL DSN for use in Go libraries (key=value format)
func BuildPostgresDSN(cfg map[string]any) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%v sslmode=%s TimeZone=%s",
		cfg["host"], cfg["user"], cfg["pass"], cfg["database"],
		cfg["port"], cfg["sslmode"], cfg["timezone"],
	)
}

// BuildPostgresURL builds PostgreSQL DSN in URL format (for goose CLI)
func BuildPostgresURL(cfg map[string]any) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=%s&TimeZone=%s",
		url.PathEscape(fmt.Sprint(cfg["user"])),
		url.PathEscape(fmt.Sprint(cfg["pass"])),
		cfg["host"], cfg["port"],
		cfg["database"], cfg["sslmode"],
		url.QueryEscape(fmt.Sprint(cfg["timezone"])),
	)
}

// IsDebug checks if the app is in debug mode
func IsDebug() bool {
	return x.Config().GetBool("app.debug", false)
}

// IsDevEnv checks if the app is in development mode
func IsDevEnv() bool {
	return x.Config().GetString("app.env", "production") == "dev"
}

// Is Production env
func IsProductionEnv() bool {
	return x.Config().GetString("app.env", "production") == "production"
}

// IsDemo returns true when DEMO_MODE=true is set in the environment.
// In demo mode, all write operations on sensitive services
// (FreeRADIUS, PostgreSQL config, backups) are blocked.
func IsDemo() bool {
	return x.Config().GetBool("app.demo", false)
}

// Is Https enabled
func IsHttpsEnabled() bool {
	return x.Config().GetBool("app.https.enabled", false)
}

// MapToStruct converts a map[string]any to a struct using JSON marshaling/unmarshaling
// This handles type conversions automatically
func MapToStruct(data map[string]any, result any) error {
	// Marshal the map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal map: %w", err)
	}

	// Unmarshal JSON to the target struct
	if err := json.Unmarshal(jsonData, result); err != nil {
		return fmt.Errorf("failed to unmarshal to struct: %w", err)
	}

	return nil
}

// structToMap converts a struct to map[string]any using JSON marshaling
func StructToMap(data any) map[string]any {
	result := make(map[string]any)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return result
	}

	_ = json.Unmarshal(jsonData, &result)
	return result
}

func GetAppSecret() string {
	return x.Config().GetString("app.secret", "hxdCTfhtkyJBVE01k8vvtaMHbzTmr401QqGl1111")
}
