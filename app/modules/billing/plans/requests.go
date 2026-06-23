package plans

import "github.com/imohamedsheta/xapp/app/shared/requests"

type PlanFilters struct {
	IsActive *bool             `form:"is_active" json:"is_active"`
	SearchBy map[string]string `json:"search_by" form:"search_by"`
	requests.SearchFilters
	requests.SortingFilters
	requests.PaginationFilters
}

// --- Price Request ---
type CreatePlanPriceRequest struct {
	Price    float64 `form:"price" json:"price" validate:"required,min=0"`
	Currency string  `form:"currency" json:"currency" validate:"required,len=3"` // EGP, USD, ...
}

// --- Settings Request ---
type CreatePlanSettingsRequest struct {
	// TODO: Add plan settings fields
}

// --- Main Plan Request ---
type CreatePlanRequest struct {
	Name     string                    `form:"name" json:"name" validate:"required,min=3,max=150"`
	Features []string                  `form:"features[]" json:"features" validate:"omitempty,dive,max=100"`
	IsActive bool                      `form:"is_active" json:"is_active" validate:"boolean"`
	Popular  bool                      `form:"popular" json:"popular" validate:"boolean"`
	Prices   []CreatePlanPriceRequest  `form:"prices" json:"prices" validate:"required,dive"`
	Settings CreatePlanSettingsRequest `form:"settings" json:"settings" validate:"required"`
	requests.Request
}

type UpdatePlanRequest struct {
	Name     string                     `form:"name" json:"name" validate:"omitempty,min=3,max=150"`
	Features []string                   `form:"features[]" json:"features" validate:"omitempty,dive,max=100"`
	IsActive *bool                      `form:"is_active" json:"is_active" validate:"omitempty,boolean"`
	Popular  bool                       `form:"popular" json:"popular" validate:"boolean"`
	Prices   []CreatePlanPriceRequest   `form:"prices" json:"prices" validate:"omitempty,dive"`
	Settings *CreatePlanSettingsRequest `form:"settings" json:"settings" validate:"omitempty"`
	requests.Request
}

func (r *CreatePlanRequest) Messages() map[string]string {
	return map[string]string{
		// Plan fields
		"name.required": "اسم الباقة مطلوب",
		"name.min":      "اسم الباقة يجب أن يكون أكبر من 3 حروف",
		"name.max":      "اسم الباقة يجب أن يكون أقل من 150 حرف",

		// Features
		"features.*.max": "كل ميزة يجب ألا تزيد عن 100 حرف",

		// IsActive
		"is_active.boolean": "حالة التفعيل يجب أن تكون true أو false",
		"popular.boolean":   " شعبية يجب أن تكون true أو false",

		// Prices
		"prices.required":            "يجب تحديد الأسعار للباقة",
		"prices.*.price.required":    "السعر مطلوب",
		"prices.*.price.min":         "السعر يجب أن يكون صفر أو أكثر",
		"prices.*.currency.required": "العملة مطلوبة",
		"prices.*.currency.len":      "العملة يجب أن تكون من 3 حروف (مثال: EGP, USD)",

		// --- Settings - limits ---
		// TODO: Add plan settings messages
	}
}

func (r *UpdatePlanRequest) Messages() map[string]string {
	return map[string]string{
		// Plan fields
		"name.required": "اسم الباقة مطلوب",
		"name.min":      "اسم الباقة يجب أن يكون أكبر من 3 حروف",
		"name.max":      "اسم الباقة يجب أن يكون أقل من 150 حرف",

		"features.*.max": "كل ميزة يجب ألا تزيد عن 100 حرف",

		"is_active.boolean": "حالة التفعيل يجب أن تكون true أو false",
		"popular.boolean":   " شعبية يجب أن تكون true أو false",

		// Prices
		"prices.required":            "يجب تحديد الأسعار للباقة",
		"prices.*.price.required":    "السعر مطلوب",
		"prices.*.price.min":         "السعر يجب أن يكون صفر أو أكثر",
		"prices.*.currency.required": "العملة مطلوبة",
		"prices.*.currency.len":      "العملة يجب أن تكون من 3 حروف (مثال: EGP, USD)",

		// Settings - limits
		// TODO: Add plan settings messages
	}
}
