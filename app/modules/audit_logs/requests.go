package audit_logs

import (
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/requests"
)

type AuditLogFilters struct {
	requests.SearchFilters
	requests.SortingFilters
	requests.PaginationFilters
	Action          string              `json:"action" query:"action" form:"action"`
	AuditableType   enums.AuditableType `json:"auditable_type" query:"auditable_type" form:"auditable_type"`
	AuditableId     int64               `json:"auditable_id" query:"auditable_id" form:"auditable_id"`
	FromDate        string              `json:"from_date" query:"from_date" form:"from_date"`
	ToDate          string              `json:"to_date" query:"to_date" form:"to_date"`
	Group           string              `json:"group" query:"group" form:"group"`
	UserId          int64               `json:"user_id" query:"user_id" form:"user_id"`
	SearchBy        map[string]string   `json:"search_by" form:"search_by"`
	IsClientView    bool                `json:"-"`
	OwnUserId       int64               `json:"-"`
	OwnRadiusUserId int64               `json:"-"`
	requests.Request
}

func (r *AuditLogFilters) Messages() map[string]string {
	return map[string]string{}
}
