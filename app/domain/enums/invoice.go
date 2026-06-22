package enums

// invoice creator type (who created the invoice)
type InvoiceCreatorType string

const (
	InvoiceCreatorTypeUser InvoiceCreatorType = "user"
)

// invoice user type (for whom the invoice is created for)
type InvoiceUserType string

const (
	InvoiceUserTypeUser InvoiceUserType = "user"
)

// invoiceable type (what the invoice is related to)
type InvoiceableType string

const (
	InvoiceableTypeSubscription InvoiceableType = "subscription"
)

// invoice type (what the invoice is for)
type InvoiceType string

const (
	InvoiceTypeSubscription      InvoiceType = "subscription"
	InvoiceTypeRenewSubscription InvoiceType = "renew_subscription"
	InvoiceTypeUpgrade           InvoiceType = "upgrade"
	InvoiceTypeDowngrade         InvoiceType = "downgrade"
	InvoiceTypeManualCharge      InvoiceType = "manual_charge"
	InvoiceTypeManualDebit       InvoiceType = "manual_debit"
	InvoiceTypeRefund            InvoiceType = "refund"
)

func IsValidInvoiceType(t InvoiceType) bool {
	switch t {
	case InvoiceTypeSubscription,
		InvoiceTypeRenewSubscription,
		InvoiceTypeUpgrade,
		InvoiceTypeDowngrade,
		InvoiceTypeManualCharge,
		InvoiceTypeManualDebit,
		InvoiceTypeRefund:
		return true
	default:
		return false
	}
}

// invoice status (what the invoice status is)
type InvoiceStatus string

const (
	InvoiceStatusPending       InvoiceStatus = "pending"
	InvoiceStatusPaid          InvoiceStatus = "paid"
	InvoiceStatusPartiallyPaid InvoiceStatus = "partially_paid"
	InvoiceStatusFailed        InvoiceStatus = "failed"
	InvoiceStatusRefunded      InvoiceStatus = "refunded"
	InvoiceStatusCanceled      InvoiceStatus = "canceled"
)

func IsValidInvoiceStatus(status InvoiceStatus) bool {
	switch status {
	case InvoiceStatusPending,
		InvoiceStatusPaid,
		InvoiceStatusPartiallyPaid,
		InvoiceStatusFailed,
		InvoiceStatusRefunded,
		InvoiceStatusCanceled:
		return true
	default:
		return false
	}
}

// statistic trend (what the statistic trend is)
type StatisticTrend string

const (
	StatisticTrendGrowing  StatisticTrend = "growing"
	StatisticTrendLowering StatisticTrend = "lowering"
	StatisticTrendStable   StatisticTrend = "stable"
)
