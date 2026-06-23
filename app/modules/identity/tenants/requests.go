package tenants

import (
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/requests"
)

type TenantManagerFilters struct {
	Status   enums.TenantStatus `json:"status" form:"status"`
	SearchBy map[string]string  `json:"search_by" form:"search_by"`
	requests.SearchFilters
	requests.SortingFilters
	requests.PaginationFilters
}

type UpdateTenantStatusRequest struct {
	Status enums.TenantStatus `json:"status" form:"status" validate:"required,oneof=1 2"`
	requests.Request
}

func (r *UpdateTenantStatusRequest) Messages() map[string]string {
	return map[string]string{
		"status.required": "الحالة مطلوبة",
		"status.oneof":    "الحالة يجب أن تكون مفعلة أو غير مفعلة",
	}
}

type TenantImpersonateRequest struct {
	TenantId int64 `json:"tenant_id" form:"tenant_id" validate:"required,gt=0,exists_db=tenants-id"`
	requests.Request
}

func (r *TenantImpersonateRequest) Messages() map[string]string {
	return map[string]string{
		"tenant_id.required":  "الوكيل مطلوب",
		"tenant_id.gt":        "الوكيل يجب أن يكون أكبر من 0",
		"tenant_id.exists_db": "الوكيل غير موجود",
	}
}

type UserImpersonateRequest struct {
	UserId int64 `json:"user_id" form:"user_id" validate:"required,gt=0,exists_db=users-id"`
	requests.Request
}

func (r *UserImpersonateRequest) Messages() map[string]string {
	return map[string]string{
		"user_id.required":  "المستخدم مطلوب",
		"user_id.gt":        "المستخدم يجب أن يكون أكبر من 0",
		"user_id.exists_db": "المستخدم غير موجود",
	}
}

type TenantImpersonateLeaveRequest struct {
	Password string `json:"password" validate:"required,min=8,max=30" example:"123456789"`
	requests.Request
}

func (r *TenantImpersonateLeaveRequest) Messages() map[string]string {
	return map[string]string{
		"password.required": "كلمة المرور مطلوبة",
		"password.min":      "كلمة المرور يجب أن تكون 8 أحرف على الأقل",
		"password.max":      "كلمة المرور يجب أن تكون 30 حرف على الأكثر",
	}
}
