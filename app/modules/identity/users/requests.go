package users

import (
	"context"

	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/requests"
)

// UserManagerFilters - Filters for manager's user list with role, status, and tenant filtering
type UserManagerFilters struct {
	UserFilters
	Role     enums.UserRole   `json:"role" form:"role"`
	Status   enums.UserStatus `json:"status" form:"status"`
	TenantId int64            `json:"tenant_id" form:"tenant_id"`
}

type UserFilters struct {
	SearchBy  map[string]string `json:"search_by" form:"search_by"`
	UsersRole []enums.UserRole  `json:"users_role" form:"users_role"`
	requests.SearchFilters
	requests.SortingFilters
	requests.PaginationFilters
}

type CreateUserRequest struct {
	Name                 string           `json:"name" validate:"required,min=3,max=30" example:"John"`
	Username             string           `json:"username" validate:"omitempty,alphanum,min=3,max=255,unique_db=users-username" example:"mohamedsheta"`
	Email                string           `json:"email" validate:"omitempty,email,unique_db=users-email" example:"test@test.com"`
	Password             string           `json:"password" validate:"required,min=8,max=30" example:"123456789"`
	PasswordConfirmation string           `json:"password_confirmation" validate:"required" example:"123456789"`
	Role                 string           `json:"role" validate:"required,oneof=admin distributor" example:"admin"`
	Phone                string           `json:"phone" validate:"omitempty,max=20"`
	NasSerials           []int64          `json:"nas_serials"`
	UserCardsLimit       int64            `json:"user_cards_limit"`
	RechargeCardsLimit   int64            `json:"recharge_cards_limit"`
	Notes                string           `json:"notes" validate:"omitempty,max=500"`
	Status               enums.UserStatus `json:"status"`
	requests.Request
}

func (r *CreateUserRequest) Messages() map[string]string {
	return map[string]string{
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
		"password_confirmation.required": "تأكيد كلمة المرور مطلوب",
		"role.required":                  "الدور مطلوب",
		"role.oneof":                     "الدور يجب أن يكون إما 'مدير' أو 'موزع'.",
		"phone.max":                      "رقم الهاتف يجب أن لا يتجاوز 20 رقم",
		"notes.max":                      "الملاحظات يجب أن لا تتجاوز 500 حرف",
	}
}

func (r *CreateUserRequest) Validate(ctx context.Context) (map[string]any, error) {
	errors := make(map[string]any)

	if r.Email == "" && r.Username == "" {
		errors["username"] = "يرجى إدخال اسم المستخدم أو البريد الإلكتروني."
	}

	if r.Password != r.PasswordConfirmation {
		errors["password_confirmation"] = "تأكيد كلمة المرور غير صحيح."
	}

	return errors, nil
}

type UpdateUserRequest struct {
	Id                   int64            `json:"id" validate:"required"`
	Name                 string           `json:"name" validate:"required,min=3,max=30" example:"John"`
	Username             string           `json:"username" validate:"omitempty,alphanum,min=3,max=255" example:"mohamedsheta"`
	Email                string           `json:"email" validate:"omitempty,email" example:"test@test.com"`
	Password             string           `json:"password" validate:"omitempty,min=8,max=30" example:"123456789"`
	PasswordConfirmation string           `json:"password_confirmation" validate:"omitempty" example:"123456789"`
	Role                 string           `json:"role" validate:"required,oneof=admin distributor" example:"admin"`
	Phone                string           `json:"phone" validate:"omitempty,max=20"`
	NasSerials           []int64          `json:"nas_serials"`
	UserCardsLimit       int64            `json:"user_cards_limit"`
	RechargeCardsLimit   int64            `json:"recharge_cards_limit"`
	Notes                string           `json:"notes" validate:"omitempty,max=500"`
	Status               enums.UserStatus `json:"status"`
	requests.Request
}

func (r *UpdateUserRequest) Messages() map[string]string {
	return map[string]string{
		"id.required":                   "معرف المستخدم مطلوب",
		"name.required":                 "الاسم مطلوب",
		"name.min":                      "الاسم يجب أن يكون 3 أحرف على الأقل",
		"name.max":                      "الاسم يجب أن يكون 30 حرف على الأكثر",
		"username.required":             "اسم المستخدم مطلوب",
		"username.alphanum":             "اسم المستخدم يجب أن يكون أحرف وأرقام فقط",
		"username.min":                  "اسم المستخدم يجب أن يكون 3 أحرف على الأقل",
		"username.max":                  "اسم المستخدم يجب أن يكون 255 حرف على الأكثر",
		"email.required":                "البريد الإلكتروني مطلوب",
		"email.email":                   "البريد الإلكتروني غير صحيح",
		"password.min":                  "كلمة المرور يجب أن تكون 8 أحرف على الأقل",
		"password.max":                  "كلمة المرور يجب أن تكون 30 حرف على الأكثر",
		"password_confirmation.eqfield": "تأكيد كلمة المرور غير صحيح",
		"role.required":                 "الدور مطلوب",
		"role.oneof":                    "الدور يجب أن يكون إما 'مدير' أو 'موزع'.",
		"phone.max":                     "رقم الهاتف يجب أن لا يتجاوز 20 رقم",
		"notes.max":                     "الملاحظات يجب أن لا تتجاوز 500 حرف",
	}
}

func (r *UpdateUserRequest) Validate(ctx context.Context) (map[string]any, error) {
	errors := make(map[string]any)

	if r.Email == "" && r.Username == "" {
		errors["username"] = "يرجى إدخال اسم المستخدم أو البريد الإلكتروني."
	}

	if r.Password != r.PasswordConfirmation {
		errors["password_confirmation"] = "تأكيد كلمة المرور غير صحيح."
	}

	return errors, nil
}

type UpdateUserPermissionsRequest struct {
	Id          int64    `json:"id" validate:"required"`
	Permissions []string `json:"permissions" validate:"required"`
	requests.Request
}

func (r *UpdateUserPermissionsRequest) Messages() map[string]string {
	return map[string]string{
		"id.required":          "معرف المستخدم مطلوب",
		"permissions.required": "الصلاحيات مطلوبة",
	}
}

func (r *UpdateUserPermissionsRequest) Validate(ctx context.Context) (map[string]any, error) {
	return nil, nil
}

type UpdateProfileRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=30"`
	Username string `json:"username" validate:"required,alphanum,min=3,max=255"`
	Email    string `json:"email" validate:"required,email"`
	requests.Request
}

func (r *UpdateProfileRequest) Messages() map[string]string {
	return map[string]string{
		"name.required":     "الاسم مطلوب",
		"name.min":          "الاسم يجب أن يكون على الأقل 3 أحرف",
		"username.required": "اسم المستخدم مطلوب",
		"username.alphanum": "اسم المستخدم يجب أن يحتوي على أحرف وأرقام فقط",
		"email.required":    "البريد الإلكتروني مطلوب",
		"email.email":       "البريد الإلكتروني غير صحيح",
	}
}

type UpdatePasswordRequest struct {
	CurrentPassword      string `json:"current_password" validate:"required"`
	Password             string `json:"password" validate:"required,min=8,max=30"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required"`
	requests.Request
}

func (r *UpdatePasswordRequest) Messages() map[string]string {
	return map[string]string{
		"current_password.required":      "كلمة المرور الحالية مطلوبة",
		"password.required":              "كلمة المرور الجديدة مطلوبة",
		"password.min":                   "كلمة المرور الجديدة يجب أن تكون على الأقل 8 أحرف",
		"password_confirmation.required": "تأكيد كلمة المرور مطلوبة",
		"password_confirmation.eqfield":  "تأكيد كلمة المرور غير متطابق",
	}
}

func (r *UpdatePasswordRequest) Validate(ctx context.Context) (map[string]any, error) {
	msgs := r.Messages()
	errs := make(map[string]any)
	if r.Password != r.PasswordConfirmation {
		errs["password_confirmation"] = msgs["password_confirmation.eqfield"]
	}

	return errs, nil
}
