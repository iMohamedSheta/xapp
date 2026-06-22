package models

import (
	"database/sql"
	"time"

	"github.com/imohamedsheta/xapp/app/domain/enums"
)

type Invoice struct {
	Id       int64 `xqb:"id" json:"id"`
	TenantId int64 `xqb:"tenant_id" json:"tenant_id"`

	CreatorId   sql.NullInt64            `xqb:"creator_id" json:"creator_id"`
	CreatorType enums.InvoiceCreatorType `xqb:"creator_type" json:"creator_type"` // could be user in the system
	CreatorRole sql.NullString           `xqb:"creator_role" json:"creator_role"`

	UserId   sql.NullInt64         `xqb:"user_id" json:"user_id"`
	UserType enums.InvoiceUserType `xqb:"user_type" json:"user_type"`

	OrderId       sql.NullInt64 `xqb:"order_id" json:"order_id"`
	TransactionId sql.NullInt64 `xqb:"transaction_id" json:"transaction_id"`

	InvoiceableType sql.NullString `xqb:"invoiceable_type" json:"invoiceable_type"` // enums.InvoiceableType
	InvoiceableId   sql.NullInt64  `xqb:"invoiceable_id" json:"invoiceable_id"`

	InvoiceNumber string              `xqb:"invoice_number" json:"invoice_number"`
	Type          enums.InvoiceType   `xqb:"type" json:"type"`
	Status        enums.InvoiceStatus `xqb:"status" json:"status"`

	Amount   float64 `xqb:"amount" json:"amount"`
	Paid     float64 `xqb:"paid" json:"paid"`
	Currency string  `xqb:"currency" json:"currency"`

	DueDate sql.NullTime `xqb:"due_date" json:"due_date"`
	PaidAt  sql.NullTime `xqb:"paid_at" json:"paid_at"`

	Notes    sql.NullString `xqb:"notes" json:"notes"`
	Metadata []byte         `xqb:"metadata" json:"metadata"` // JSONB

	CreatedAt time.Time `xqb:"created_at" json:"created_at"`
	UpdatedAt time.Time `xqb:"updated_at" json:"updated_at"`

	// Relationships
	Tenant      *Tenant      `table:"tenants" json:"tenant,omitempty"`
	User        *User        `table:"users" json:"user,omitempty"`
	Order       *Order       `table:"orders" json:"order,omitempty"`
	Transaction *Transaction `table:"transactions" json:"transaction,omitempty"`

	BaseModel
}

func (Invoice) Table() string {
	return "invoices"
}

func (m Invoice) Cols() []any {
	return m.BaseModel.Cols(m, "invoice")
}
