package requests

import (
	"github.com/imohamedsheta/xapp/app/shared/requests"
)

type SubscribeRequest struct {
	PlanId int64 `form:"plan_id" json:"plan_id" validate:"required,min=1,exists_db=plans-id"`
	requests.Request
}

func (r *SubscribeRequest) Messages() map[string]string {
	return map[string]string{
		"plan_id.required":  "الباقة المرجعة مطلوبة",
		"plan_id.min":       "الباقة المرجعة يجب أن تكون أكبر من 1",
		"plan_id.exists_db": "الباقة المرجعة غير موجودة",
	}
}

type AddBalanceRequest struct {
	Amount         float64 `form:"amount" json:"amount" validate:"required,min=1"`
	PaymentGateway string  `form:"payment_gateway" json:"payment_gateway" validate:"required,oneof= paymob paypal"`
	PaymentMethod  string  `form:"payment_method" json:"payment_method" validate:"required,oneof= credit mobile"`
	requests.Request
}

func (r *AddBalanceRequest) Messages() map[string]string {
	return map[string]string{
		"amount.required":          "المبلغ المطلوب إضافته مطلوب",
		"amount.min":               "المبلغ المطلوب إضافته يجب أن يكون أكبر من 0",
		"payment_gateway.required": "بوابة الدفع مطلوبة",
		"payment_gateway.oneof":    "بوابة الدفع غير مدعومة",
		"payment_method.required":  "طريقة الدفع مطلوبة",
		"payment_method.oneof":     "طريقة الدفع غير مدعومة",
	}
}
