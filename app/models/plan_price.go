package models

import (
	"time"
)

type PlanPrice struct {
	Id        int64     `xqb:"id" json:"id"`
	PlanId    int64     `xqb:"plan_id" json:"plan_id"`
	Price     float64   `xqb:"price" json:"price"`
	Discount  float64   `xqb:"discount" json:"discount"`
	Currency  string    `xqb:"currency" json:"currency"`
	CreatedAt time.Time `xqb:"created_at" json:"created_at"`
	UpdatedAt time.Time `xqb:"updated_at" json:"updated_at"`
	BaseModel
}

func (PlanPrice) Table() string {
	return "plan_prices"
}

func (m PlanPrice) Cols() []any {
	return m.BaseModel.Cols(m, "plan_prices")
}
