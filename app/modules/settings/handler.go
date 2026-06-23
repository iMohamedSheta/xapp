// modules/settings/settings_handler.go
package settings

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/http/handler"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/x"
)

type SettingsHandler struct {
	*handler.Handler
	settingAction *SettingAction
}

func NewSettingsHandler(base *handler.Handler, settingAction *SettingAction) *SettingsHandler {
	return &SettingsHandler{
		Handler:       base,
		settingAction: settingAction,
	}
}

func (h *SettingsHandler) AppearanceView(c *gin.Context) error {
	return h.Inertia.Render(c, "Settings/Appearance", nil)
}

func (h *SettingsHandler) SaveSettings(c *gin.Context) error {
	var req SaveSettingRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		var xe *xerr.XErr
		if errors.As(err, &xe) && xe.IsType(enums.XErrValidationError) {
			return nil
		}
		return err
	}

	user, err := x.AuthUser(c)
	if err != nil {
		return err
	}

	if err := h.settingAction.Save(c, &user.TenantId, &req); err != nil {
		return err
	}

	return h.BackWithFlash(c, "تم تعديل الاعدادات بنجاح", 201)
}
