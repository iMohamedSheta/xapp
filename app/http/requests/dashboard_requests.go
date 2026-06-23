package requests

import "github.com/imohamedsheta/xapp/app/shared/requests"

type DashboardFilters struct {
	Period string `json:"period" form:"period" validate:"omitempty,oneof=today week month all"`
	requests.SearchFilters
}

func (r *DashboardFilters) Messages() map[string]string {
	return map[string]string{
		"period.oneof": "الفترة يجب أن تكون واحدة من: اليوم، الأسبوع، الشهر، الكل",
	}
}
