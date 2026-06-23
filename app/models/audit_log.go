package models

import (
	"database/sql"
	"time"

	"github.com/imohamedsheta/xapp/app/shared/enums"
)

/*
*
| AuditLog represents an auditable action performed by a user on a resource.
| It tracks who did what, when, and on which resource, providing a complete
| audit trail for security, compliance, and debugging purposes.
*
*/
type AuditLog struct {
	Id            int64                `xqb:"id" json:"id"`
	TenantId      int64                `xqb:"tenant_id" json:"tenant_id"`
	UserId        sql.NullInt64        `xqb:"user_id" json:"user_id"`
	UserType      sql.NullString       `xqb:"user_type" json:"user_type"` // max:20
	AuditableId   sql.NullInt64        `xqb:"auditable_id" json:"auditable_id"`
	AuditableType enums.AuditableType  `xqb:"auditable_type" json:"auditable_type"` // max:50
	Action        enums.AuditLogAction `xqb:"action" json:"action"`
	Summary       string               `xqb:"summary" json:"summary"` // max:255
	Details       map[string]any       `xqb:"details" json:"details"`
	CreatedAt     time.Time            `xqb:"created_at" json:"created_at"`
	UpdatedAt     time.Time            `xqb:"updated_at" json:"updated_at"`
	BaseModel
}

func (AuditLog) Table() string {
	return "audit_logs"
}
