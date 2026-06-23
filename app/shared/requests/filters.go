package requests

// Request is the base struct embedded by all request DTOs.
type Request struct {
	RequestSentFields map[string]any `json:"-"`
}

func (r *Request) GetRequestSentFields() map[string]any {
	return r.RequestSentFields
}

func (r *Request) SetRequestSentFields(requestSentFields map[string]any) {
	r.RequestSentFields = requestSentFields
}

// SearchFilters holds common search query parameters.
type SearchFilters struct {
	Search string `json:"search,omitempty" form:"search"`
}

// SortingFilters holds common sorting query parameters.
type SortingFilters struct {
	SortBy    string `json:"sort_by,omitempty" form:"sort_by"`
	SortOrder string `json:"sort_order,omitempty" form:"sort_order"`
}

// PaginationFilters holds common pagination query parameters.
type PaginationFilters struct {
	Page    int `json:"page,omitempty" form:"page"`
	PerPage int `json:"per_page,omitempty" form:"per_page"`
}

// SetPaginationDefaults applies sensible defaults when fields are zero.
func (f *PaginationFilters) SetPaginationDefaults() {
	if f.Page == 0 {
		f.Page = 1
	}
	if f.PerPage == 0 || f.PerPage > 300 {
		f.PerPage = 25
	}
}

// DashboardFilters holds common dashboard query parameters.
type DashboardFilters struct {
	Period string `json:"period" form:"period" validate:"omitempty,oneof=today week month all"`
	SearchFilters
}

func (r *DashboardFilters) Messages() map[string]string {
	return map[string]string{
		"period.oneof": "الفترة يجب أن تكون واحدة من: اليوم، الأسبوع، الشهر، الكل",
	}
}
