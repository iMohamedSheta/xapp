package auth

import (
	"context"

	"github.com/imohamedsheta/xapp/app/shared/requests"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"omitempty,email" example:"test@test.com"`
	Username string `json:"username" validate:"omitempty,alphanum,min=3,max=255" example:"mohamedsheta"`
	Password string `json:"password" validate:"required,min=8,max=30" example:"123456789"`
	Remember bool   `json:"remember" validate:"boolean" example:"true"`
	requests.Request
}

func (r *LoginRequest) Messages() map[string]string {
	return map[string]string{
		"email.email":       "البريد الإلكتروني غير صحيح",
		"username.alphanum": "اسم المستخدم يجب أن يكون أحرف وأرقام فقط",
		"username.min":      "اسم المستخدم يجب أن يكون 3 أحرف على الأقل",
		"username.max":      "اسم المستخدم يجب أن يكون 255 حرف على الأكثر",
		"password.required": "كلمة المرور مطلوبة",
		"password.min":      "كلمة المرور يجب أن تكون 8 أحرف على الأقل",
		"password.max":      "كلمة المرور يجب أن تكون 30 حرف على الأكثر",
	}
}

// Validate performs custom conditional validation
func (r *LoginRequest) Validate(ctx context.Context) (map[string]any, error) {
	errors := make(map[string]any)

	// Validate email or username
	if r.Email == "" && r.Username == "" {
		errors["email"] = "البريد الإلكتروني أو اسم المستخدم مطلوب"
	}

	return errors, nil
}

type RegisterRequest struct {
	TenantName           string `json:"tenant_name" validate:"required,min=3,max=150" example:"tenant_name"`
	Name                 string `json:"name" validate:"required,min=3,max=30" example:"John"`
	Username             string `json:"username" validate:"required,alphanum,min=3,max=255,unique_db=users-username" example:"mohamedsheta"`
	Email                string `json:"email" validate:"required,email,unique_db=users-email" example:"test@test.com"`
	Password             string `json:"password" validate:"required,min=8,max=30" example:"123456789"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password" example:"123456789"`
	requests.Request
}

func (r *RegisterRequest) Messages() map[string]string {
	return map[string]string{
		"tenant_name.required":           "الوكالة مطلوبة",
		"tenant_name.min":                "الوكالة يجب أن تكون 3 أحرف على الأقل",
		"tenant_name.max":                "الوكالة يجب أن تكون 150 حرف على الأكثر",
		"name.required":                  "الاسم مطلوب",
		"name.min":                       "الاسم يجب أن يكون 3 أحرف على الأقل",
		"name.max":                       "الاسم يجب أن يكون 30 حرف على الأكثر",
		"username.required":              "اسم المستخدم مطلوب",
		"username.alphanum":              "اسم المستخدم يجب أن يكون أحرف وأرقام فقط",
		"username.min":                   "اسم المستخدم يجب أن يكون 3 أحرف على الأقل",
		"username.max":                   "اسم المستخدم يجب أن يكون 255 حرف على الأكثر",
		"username.unique_db":             "اسم المستخدم موجود بالفعل",
		"email.required":                 "البريد الإلكتروني مطلوب",
		"email.email":                    "البريد الإلكتروني غير صحيح",
		"email.unique_db":                "البريد الإلكتروني موجود بالفعل",
		"password.required":              "كلمة المرور مطلوبة",
		"password.min":                   "كلمة المرور يجب أن تكون 8 أحرف على الأقل",
		"password.max":                   "كلمة المرور يجب أن تكون 30 حرف على الأكثر",
		"password_confirmation.required": "كلمة المرور مطلوبة",
		"password_confirmation.eqfield":  "كلمة المرور غير متطابقة",
	}
}
