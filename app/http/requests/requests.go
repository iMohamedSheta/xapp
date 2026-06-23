package requests

import (
	"github.com/imohamedsheta/xapp/app/shared/requests"
)

type PasswordConfirmationRequest struct {
	Password string `json:"password" validate:"required,min=8,max=30" example:"123456789"`
	requests.Request
}

func (r *PasswordConfirmationRequest) Messages() map[string]string {
	return map[string]string{
		"password.required": "كلمة المرور مطلوبة",
		"password.min":      "كلمة المرور يجب أن تكون 8 أحرف على الأقل",
		"password.max":      "كلمة المرور يجب أن تكون 30 حرف على الأكثر",
	}
}
