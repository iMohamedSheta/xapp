package plans

import (
	"database/sql"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/shared/enums"
)

type PlanSettingRepository struct {
}

func NewPlanSettingRepository() *PlanSettingRepository {
	return &PlanSettingRepository{}
}

func (t *PlanSettingRepository) Create(c *gin.Context, planSetting []*models.PlanSetting, tx *sql.Tx) (insertedId int64, err error) {
	if len(planSetting) == 0 {
		return 0, xerr.New("can not create plan with empty plan list", enums.XErrBadRequestError, nil)
	}

	insertedValues := make([]map[string]any, 0)
	for _, setting := range planSetting {
		insertTime := time.Now()
		insertedValues = append(insertedValues, map[string]any{
			"plan_id":           setting.PlanId,
			"plan_limits":       setting.PlanLimits,
			"expire_action":     setting.ExpireAction,
			"downgrade_to_plan": setting.DowngradeToPlan,
			"grace_period_days": setting.GracePeriodDays,
			"created_at":        insertTime,
			"updated_at":        insertTime,
		})
	}

	return xqb.Table("plan_settings").WithContext(c).WithTx(tx).InsertGetId(insertedValues)
}

func (r *PlanSettingRepository) UpdateByPlanID(c *gin.Context, planId int64, updatedFields map[string]any, tx *sql.Tx) error {
	if planId == 0 {
		return xerr.New("invalid plan id", enums.XErrBadRequestError, nil)
	}

	if len(updatedFields) == 0 {
		return nil
	}

	_, err := xqb.Table("plan_settings").
		WithContext(c).
		WithTx(tx).
		Where("plan_id", "=", planId).
		Update(updatedFields)
	return err
}

func (t *PlanSettingRepository) DeleteByPlanId(c *gin.Context, planId int64, tx *sql.Tx) error {
	if planId == 0 {
		return xerr.New("invalid plan id", enums.XErrBadRequestError, nil)
	}

	_, err := xqb.Table("plan_prices").
		WithContext(c).
		WithTx(tx).
		Where("plan_id", "=", planId).
		Delete()

	if err != nil {
		return err
	}

	return nil
}
