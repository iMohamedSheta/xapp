package enums

type OrderStatus int16

const (
	OrderPending OrderStatus = iota + 1
	OrderPaid
	OrderCancelled
)

type OrderableType string

const (
	OrderableTypeBalance OrderableType = "balance"
)
