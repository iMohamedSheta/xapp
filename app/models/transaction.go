package models

import (
	"database/sql"
	"time"
)

// Transaction model matches the `transactions` table
type Transaction struct {
	Id       int64         `xqb:"id" json:"id"`
	TenantId int64         `xqb:"tenant_id" json:"tenant_id"`
	UserId   sql.NullInt64 `xqb:"user_id" json:"user_id"`
	OrderId  int64         `xqb:"order_id" json:"order_id"`

	// Gateway info
	Gateway         string `xqb:"gateway" json:"gateway"`
	GatewayTxId     string `xqb:"gateway_tx_id" json:"gateway_tx_id,omitempty"`
	GatewayOrderId  string `xqb:"gateway_order_id" json:"gateway_order_id,omitempty"`
	MerchantOrderId string `xqb:"merchant_order_id" json:"merchant_order_id,omitempty"`

	// Amounts
	AmountCents    int64  `xqb:"amount_cents" json:"amount_cents"`
	Currency       string `xqb:"currency" json:"currency"`
	CapturedAmount int64  `xqb:"captured_amount" json:"captured_amount"`
	RefundedAmount int64  `xqb:"refunded_amount" json:"refunded_amount"`

	// Status
	Success      bool   `xqb:"success" json:"success"`
	Status       string `xqb:"status" json:"status,omitempty"`
	ResponseCode string `xqb:"response_code" json:"response_code,omitempty"`
	Message      string `xqb:"message" json:"message,omitempty"`

	// Payment Method Info (non-sensitive)
	PaymentMethod string `xqb:"payment_method" json:"payment_method,omitempty"`
	CardType      string `xqb:"card_type" json:"card_type,omitempty"`
	CardLast4     string `xqb:"card_last4" json:"card_last4,omitempty"`

	// Gateway Metadata
	AuthorizeId  string `xqb:"authorize_id" json:"authorize_id,omitempty"`
	ReceiptNo    string `xqb:"receipt_no" json:"receipt_no,omitempty"`
	BatchNo      string `xqb:"batch_no" json:"batch_no,omitempty"`
	GatewayRefId string `xqb:"gateway_ref_id" json:"gateway_ref_id,omitempty"`

	// Customer Snapshot
	CustomerEmail string `xqb:"customer_email" json:"customer_email,omitempty"`
	CustomerPhone string `xqb:"customer_phone" json:"customer_phone,omitempty"`
	CustomerName  string `xqb:"customer_name" json:"customer_name,omitempty"`

	// Raw JSON response
	RawResponse []byte `xqb:"raw_response" json:"raw_response"`

	// Timestamps
	CreatedAt time.Time `xqb:"created_at" json:"created_at"`
	UpdatedAt time.Time `xqb:"updated_at" json:"updated_at"`

	// Relations
	Order *Order `table:"orders" json:"order,omitempty"`

	BaseModel
}

func (Transaction) Table() string {
	return "transactions"
}

func (m Transaction) Cols() []any {
	return m.BaseModel.Cols(m, "transaction")
}
