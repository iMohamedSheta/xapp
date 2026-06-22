package invoices

import (
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/requests"
)

type InvoiceFilters struct {
	UserId   *int64                `json:"user_id" form:"user_id"`
	Type     enums.InvoiceType     `json:"type" form:"type"`
	Status   enums.InvoiceStatus   `json:"status" form:"status"`
	SearchBy map[string]string     `json:"search_by" form:"search_by"`
	FromDate string                `json:"from_date" form:"from_date"`
	ToDate   string                `json:"to_date" form:"to_date"`
	UserType enums.InvoiceUserType `json:"user_type" form:"user_type"`
	requests.Request
	requests.SearchFilters
	requests.SortingFilters
	requests.PaginationFilters
}

func (req *InvoiceFilters) Messages() map[string]string {
	return nil
}

type PayInvoiceRequest struct {
	InvoiceId int64   `json:"invoice_id" form:"invoice_id" validate:"required"`
	Amount    float64 `json:"amount" form:"amount" validate:"required,min=1"`
	requests.Request
}

func (req *PayInvoiceRequest) Messages() map[string]string {
	return map[string]string{
		"invoice_id": "رقم الفاتورة مطلوب",
		"amount":     "المبلغ مطلوب",
	}
}
