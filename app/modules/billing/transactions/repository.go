package transactions

import (
	"context"
	"database/sql"
	"time"

	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/models"
)

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

// Create inserts a new transaction into the database
func (r *TransactionRepository) Create(c context.Context, t *models.Transaction, sqlTx *sql.Tx) (*models.Transaction, error) {
	insertTime := time.Now()

	insertedId, err := xqb.Table(t.Table()).
		WithContext(c).
		WithTx(sqlTx).
		InsertGetId([]map[string]any{
			{
				"order_id":          t.OrderId,
				"tenant_id":         t.TenantId,
				"user_id":           t.UserId,
				"gateway":           nullableString(t.Gateway),
				"gateway_tx_id":     nullableString(t.GatewayTxId),
				"gateway_order_id":  nullableString(t.GatewayOrderId),
				"merchant_order_id": nullableString(t.MerchantOrderId),

				"amount_cents":    t.AmountCents,
				"currency":        nullableString(t.Currency),
				"captured_amount": t.CapturedAmount,
				"refunded_amount": t.RefundedAmount,

				"success":       t.Success,
				"status":        nullableString(t.Status),
				"response_code": nullableString(t.ResponseCode),
				"message":       nullableString(t.Message),

				"payment_method": nullableString(t.PaymentMethod),
				"card_type":      nullableString(t.CardType),
				"card_last4":     nullableString(t.CardLast4),

				"authorize_id":   nullableString(t.AuthorizeId),
				"receipt_no":     nullableString(t.ReceiptNo),
				"batch_no":       nullableString(t.BatchNo),
				"gateway_ref_id": nullableString(t.GatewayRefId),

				"customer_email": nullableString(t.CustomerEmail),
				"customer_phone": nullableString(t.CustomerPhone),
				"customer_name":  nullableString(t.CustomerName),

				"raw_response": t.RawResponse,

				"created_at": insertTime,
				"updated_at": insertTime,
			},
		})
	if err != nil {
		return nil, err
	}

	// Return a copy with ID + timestamps set
	return &models.Transaction{
		Id:       insertedId,
		TenantId: t.TenantId,
		UserId:   t.UserId,
		OrderId:  t.OrderId,

		Gateway:         t.Gateway,
		GatewayTxId:     t.GatewayTxId,
		GatewayOrderId:  t.GatewayOrderId,
		MerchantOrderId: t.MerchantOrderId,

		AmountCents:    t.AmountCents,
		Currency:       t.Currency,
		CapturedAmount: t.CapturedAmount,
		RefundedAmount: t.RefundedAmount,

		Success:      t.Success,
		Status:       t.Status,
		ResponseCode: t.ResponseCode,
		Message:      t.Message,

		PaymentMethod: t.PaymentMethod,
		CardType:      t.CardType,
		CardLast4:     t.CardLast4,

		AuthorizeId:  t.AuthorizeId,
		ReceiptNo:    t.ReceiptNo,
		BatchNo:      t.BatchNo,
		GatewayRefId: t.GatewayRefId,

		CustomerEmail: t.CustomerEmail,
		CustomerPhone: t.CustomerPhone,
		CustomerName:  t.CustomerName,

		RawResponse: t.RawResponse,

		CreatedAt: insertTime,
		UpdatedAt: insertTime,
	}, nil
}

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}
