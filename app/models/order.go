package models

import (
	"database/sql"
	"time"

	"github.com/imohamedsheta/xapp/app/shared/enums"
)

// Order model matches the `orders` table
type Order struct {
	Id       int64         `xqb:"id" json:"id"`
	TenantId int64         `xqb:"tenant_id" json:"tenant_id"`
	UserId   sql.NullInt64 `xqb:"user_id" json:"user_id"`

	// Polymorphic relation
	OrderableType enums.OrderableType `xqb:"orderable_type" json:"orderable_type"`
	OrderableId   int64               `xqb:"orderable_id" json:"orderable_id"`

	// Order details
	Quantity   int64             `xqb:"quantity" json:"quantity"`
	UnitPrice  float64           `xqb:"unit_price" json:"unit_price"`
	TotalPrice float64           `xqb:"total_price" json:"total_price"`
	Currency   string            `xqb:"currency" json:"currency"`
	Status     enums.OrderStatus `xqb:"status" json:"status"`

	// Timestamps
	CreatedAt time.Time `xqb:"created_at" json:"created_at"`
	UpdatedAt time.Time `xqb:"updated_at" json:"updated_at"`

	// Relations
	Tenant *Tenant ` table:"tenants" json:"tenant,omitempty"`
	User   *User   `table:"users" json:"user,omitempty"`

	BaseModel
}

func (Order) Table() string {
	return "orders"
}

func (Order) CalculateFullPrice(unitPrice float64, quantity int64) float64 {
	return unitPrice * float64(quantity)
}

func (m Order) Cols() []any {
	return m.BaseModel.Cols(m, "order")
}
