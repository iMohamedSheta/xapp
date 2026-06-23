package settings

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/imohamedsheta/xapp/app/shared/requests"
	"github.com/imohamedsheta/xapp/app/x"

	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xvalid"
)

// SaveSettingRequest is the main request wrapper
type SaveSettingRequest struct {
	SettingType enums.SettingType `form:"setting_type" json:"setting_type" validate:"required,oneof=admin_api_settings admin_app_settings admin_backup_settings"`
	SettingData map[string]any    `form:"setting_data" json:"setting_data" validate:"required"`
	requests.Request

	// Private field to cache validated data
	validatedData any `json:"-"`
}

func (r *SaveSettingRequest) Messages() map[string]string {
	return map[string]string{
		"setting_type.required": "نوع الاعدادات مطلوب",
		"setting_type.oneof":    "يجب ان يكون نوع الاعدادات صحيح",
		"setting_data.required": "البيانات مطلوبة",
	}
}

// Validate performs custom conditional validation
func (r *SaveSettingRequest) Validate(ctx context.Context) (map[string]any, error) {
	errors := make(map[string]any)

	switch r.SettingType {

	// Example:
	// case enums.AdminAppSettings:
	// 	data := &AdminAppSettingsData{}
	// 	if validationErrors := r.validateTypedData(ctx, data); validationErrors != nil {
	// 		errors = validationErrors
	// 	} else {
	// 		r.validatedData = data // Cache validated data
	// 	}

	default:
		errors["setting_type"] = r.Messages()["setting_type.oneof"]
	}

	if len(errors) > 0 {
		return errors, nil
	}

	return nil, nil
}

func (r *SaveSettingRequest) validateTypedData(ctx context.Context, target xvalid.Validatable) map[string]any {
	// Convert map to JSON then to struct
	jsonData, err := json.Marshal(r.SettingData)
	if err != nil {
		return map[string]any{"data": "تنسيق البيانات غير صحيح"}
	}

	if err := json.Unmarshal(jsonData, target); err != nil {
		return map[string]any{"data": "تنسيق البيانات غير صحيح"}
	}

	// Use your validator instance to validate the struct
	if xerror := x.Validator().ValidateRequest(ctx, target); xerror != nil {
		var ve *xvalid.ValidationError
		if errors.As(xerror, &ve) {
			// Prefix keys with "setting_data." because i am doing multiple requests validations Careful!!!!!!!!!!!!!!!!
			return utils.PrefixErrorKeys(ve.Fields, "setting_data")
		}

		return map[string]any{"data": "خطأ في التحقق من صحة البيانات"}
	}

	return nil
}

// GetValidatedData returns the validated data as interface{}
func (r *SaveSettingRequest) GetValidatedData() any {
	return r.validatedData
}
