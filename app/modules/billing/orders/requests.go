package orders

import (
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/requests"
)

type OrderFilters struct {
	Status   enums.OrderStatus `json:"status" form:"status"`
	SearchBy map[string]string `json:"search_by" form:"search_by"`
	requests.SearchFilters
	requests.SortingFilters
	requests.PaginationFilters
}

type OrderUpdateRequest struct {
	Status enums.OrderStatus `json:"status" form:"status" validate:"required,oneof=1 2 3"`
	requests.Request
}

func (r *OrderUpdateRequest) Messages() map[string]string {
	return map[string]string{
		"status.required": "الحالة مطلوبة",
		"status.oneof":    "الحالة يجب أن تكون واحدة من معلق، مدفوع، ملغي",
	}
}
