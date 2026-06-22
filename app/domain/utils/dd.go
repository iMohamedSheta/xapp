package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func DD(v ...any) {
	Dump(v...)
	Die()
}

func Dump(v ...any) {
	for _, val := range v {
		fmt.Println(formatOutput(val))
	}
}

func Die() {
	os.Exit(1)
}

func formatOutput(value any) string {
	// Handle different types
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%s%s%s%s", BlackBG, Red, v, Reset)
	case error:
		return fmt.Sprintf("%s%s%s%s", BlackBG, Red, v.Error(), Reset)
	default:
		// Use JSON for complex structures
		jsonData, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			return fmt.Sprintf("%s%s%#v%s", BlackBG, Red, value, Reset)
		}
		return fmt.Sprintf("%s%s%s%s", BlackBG, Red, jsonData, Reset)
	}
}
