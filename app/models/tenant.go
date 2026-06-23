package models

import (
	"time"

	"github.com/imohamedsheta/xapp/app/shared/enums"
)

/*
* Tenant represents an isolated business entity within the system.
| A tenant can be a company, agency, organization, government entity,
| institution, network operator, or any other business owner using the platform.
| Each tenant owns its own data, users, settings, and resources, which
| are logically isolated from other tenants.
*
*/
type Tenant struct {
	Id        int64              `xqb:"id" json:"id"`
	Name      string             `xqb:"name" json:"name"`
	Status    enums.TenantStatus `xqb:"status" json:"status"`
	Balance   float64            `xqb:"balance" json:"balance"`
	CreatedAt time.Time          `xqb:"created_at" json:"created_at"`
	UpdatedAt time.Time          `xqb:"updated_at" json:"updated_at"`
	BaseModel
}

func (Tenant) Table() string {
	return "tenants"
}

func (m Tenant) Cols() []any {
	return m.BaseModel.Cols(m, "tenant")
}

func (m *Tenant) HaveEnoughBalance(amount float64) bool {
	return m.Balance >= amount
}
