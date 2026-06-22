package enums

type TenantStatus int8

const (
	TenantStatusActive TenantStatus = iota + 1
	TenantStatusInactive
)
