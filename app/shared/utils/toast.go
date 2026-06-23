package utils

func ToastSuccess(msg string, description string) map[string]any {
	return map[string]any{
		"title":       msg,
		"type":        "success",
		"richColors":  true,
		"description": description,
		"classes": map[string]any{
			"description": "text-xs !text-gray-400 px-2",
			"title":       "px-2",
			"closeButton": "!bg-[#1f1f1f] !outline-none !border-none !text-white",
		},
		"closeButton": true,
		"cancelButtonStyle": map[string]any{
			"backgroundColor": "#1f1f1f",
		},
		"style": map[string]any{
			"backgroundColor": "#0a0a0a",
			"borderColor":     "#1f1f1f",
			"direction":       "rtl",
			"textAlign":       "right",
			"fontFamily":      "Cairo, sans-serif",
			"paddingRight":    "10px",
		},
		"position":  "bottom-center",
		"duration":  5000,
		"invert":    true,
		"important": false,
	}
}

func ToastInfo(msg string, description string) map[string]any {
	return map[string]any{
		"title":       msg,
		"type":        "info",
		"richColors":  true,
		"description": description,
		"classes": map[string]any{
			"description": "text-xs !text-gray-400 px-2",
			"title":       "px-2",
			"closeButton": "!bg-[#1f1f1f] !outline-none !border-none !text-white",
		},
		"closeButton": true,
		"cancelButtonStyle": map[string]any{
			"backgroundColor": "#1f1f1f",
		},
		"style": map[string]any{
			"backgroundColor": "#0a0a0a",
			"borderColor":     "#1f1f1f",
			"direction":       "rtl",
			"textAlign":       "right",
			"fontFamily":      "Cairo, sans-serif",
			"paddingRight":    "10px",
		},
		"position":  "bottom-center",
		"duration":  5000,
		"invert":    true,
		"important": false,
	}
}

func ToastError(msg string, description string) map[string]any {
	return map[string]any{
		"title":       msg,
		"type":        "error",
		"richColors":  true,
		"description": description,
		"classes": map[string]any{
			"description": "text-xs !text-gray-400 px-2",
			"title":       "px-2 !text-red-500",
			"icon":        "!text-red-500",
			"closeButton": "!bg-[#1f1f1f] !outline-none !border-none !text-white",
		},
		"closeButton": true,
		"cancelButtonStyle": map[string]any{
			"backgroundColor": "#1f1f1f",
		},
		"style": map[string]any{
			"backgroundColor": "#0a0a0a",
			"borderColor":     "#1f1f1f",
			"direction":       "rtl",
			"textAlign":       "right",
			"fontFamily":      "Cairo, sans-serif",
			"paddingRight":    "10px",
		},
		"position":  "bottom-center",
		"duration":  5000,
		"invert":    true,
		"important": false,
	}
}
