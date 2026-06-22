package enums

type AuditLogAction string

const (
	AuditLogActionCreate        AuditLogAction = "create"
	AuditLogActionUpdate        AuditLogAction = "update"
	AuditLogActionDelete        AuditLogAction = "delete"
	AuditLogActionRestore       AuditLogAction = "restore"
	AuditLogActionLogin         AuditLogAction = "login"
	AuditLogActionLogout        AuditLogAction = "logout"
	AuditLogActionToggleStatus  AuditLogAction = "toggle_status"
	AuditLogActionChangeBalance AuditLogAction = "change_balance"
	AuditLogActionRenew         AuditLogAction = "renew"
	AuditLogActionSubscribe     AuditLogAction = "subscribe"
	AuditLogActionChangePlan    AuditLogAction = "change_plan"
)

type AuditableType string

const (
	AuditableTypeUser AuditableType = "user"
	AuditableTypeNas  AuditableType = "nas"
)
