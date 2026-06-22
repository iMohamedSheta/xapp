package models

import (
	"context"
	"time"

	"github.com/iMohamedSheta/xqb"
)

type Plan struct {
	Id        int64     `xqb:"id" json:"id"`
	Name      string    `xqb:"name" json:"name"`
	Features  []string  `xqb:"features" json:"features"`
	IsActive  bool      `xqb:"is_active" json:"is_active"`
	Popular   bool      `xqb:"popular" json:"popular"`
	CreatedAt time.Time `xqb:"created_at" json:"created_at"`
	UpdatedAt time.Time `xqb:"updated_at" json:"updated_at"`

	PlanSetting *PlanSetting `table:"plan_setting" json:"plan_setting,omitempty"`
	PlanPrices  []*PlanPrice `table:"plan_prices" json:"plan_prices,omitempty"`
	BaseModel
}

func (Plan) Table() string {
	return "plans"
}

func (m Plan) Cols() []any {
	return m.BaseModel.Cols(m, "plan")
}

func (m *Plan) LoadPlanPrices(c context.Context) error {
	if m.PlanPrices != nil {
		return nil
	}

	pricesData, err := xqb.Table("plan_prices").WithContext(c).Where("plan_id", "=", m.Id).Get()
	if err != nil {
		return err
	}

	var prices []PlanPrice
	if err = xqb.Bind(pricesData, &prices); err != nil {
		return err
	}

	for i := range prices {
		m.PlanPrices = append(m.PlanPrices, &prices[i])
	}

	return nil
}
