package settings

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/domain/enums"
)

type SettingAction struct {
	settingRepo   *SettingRepository
	settingPolicy *SettingPolicy
}

func NewSettingAction(
	settingRepo *SettingRepository,
	settingPolicy *SettingPolicy,
) *SettingAction {
	return &SettingAction{
		settingRepo:   settingRepo,
		settingPolicy: settingPolicy,
	}
}

func (a *SettingAction) Save(c *gin.Context, tenantId *int64, req *SaveSettingRequest) error {
	// Check permission
	if !a.settingPolicy.CanUpdate(c) {
		return xerr.New("You do not have permission to update settings", enums.XErrForbiddenError, nil)
	}

	// // Get the validated typed data based on setting type
	// var settingsData map[string]any

	// switch req.SettingType {

	// default:
	// 	return fmt.Errorf("unknown setting type: %s", req.SettingType)
	// }

	// if err := a.settingRepo.SaveSetting(c, "tenant", tenantId, req.SettingType, settingsData, nil); err != nil {
	// 	return err
	// }

	return nil
}

func (a *SettingAction) Get(c context.Context, tenantId *int64, settingType enums.SettingType) (map[string]any, error) {
	return a.settingRepo.GetSettingByType(c, "tenant", tenantId, settingType)
}

func (a *SettingAction) CanViewAppSystemPage(c *gin.Context) error {
	if !a.settingPolicy.CanView(c) {
		return xerr.New("You do not have permission to view app settings page", enums.XErrForbiddenError, nil).
			WithPublicMessage("ليس لديك صلاحية للوصول لهذه الصفحة")
	}

	return nil
}
