package models

import (
	"database/sql"
	"time"

	"github.com/imohamedsheta/xapp/app/shared/enums"
)

type PlanLimits struct {
	// ApiAccess        bool  `xqb:"api_access" json:"api_access"`
}

type PlanSetting struct {
	Id              int64                  `xqb:"id" json:"id"`
	PlanId          int64                  `xqb:"plan_id" json:"plan_id"`
	PlanLimits      PlanLimits             `xqb:"plan_limits" json:"plan_limits"`
	ExpireAction    enums.PlanExpireAction `xqb:"expire_action" json:"expire_action"`
	DowngradeToPlan sql.NullInt64          `xqb:"downgrade_to_plan" json:"downgrade_to_plan"`
	GracePeriodDays int64                  `xqb:"grace_period_days" json:"grace_period_days"`
	CreatedAt       time.Time              `xqb:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `xqb:"updated_at" json:"updated_at"`
	BaseModel
}

func (PlanSetting) Table() string {
	return "plan_settings"
}

func (m PlanSetting) Cols() []any {
	return m.BaseModel.Cols(m, "plan_setting")
}
