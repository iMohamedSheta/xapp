package listeners

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/events"
	"github.com/imohamedsheta/xnotify"
)

type AuditLogRepository interface {
	Create(c context.Context, log *models.AuditLog, tx *sql.Tx) error
}

type UserLoggedInListener struct {
	notify             *xnotify.Notify
	auditLogRepository AuditLogRepository
}

func NewUserLoggedInListener(notify *xnotify.Notify, auditLogRepository AuditLogRepository) *UserLoggedInListener {
	return &UserLoggedInListener{
		notify:             notify,
		auditLogRepository: auditLogRepository,
	}
}

func (l *UserLoggedInListener) Handle(ctx context.Context, payload []byte) error {
	var event events.UserLoggedInPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return err
	}

	log := &models.AuditLog{
		UserId: sql.NullInt64{
			Int64: event.UserId,
			Valid: true,
		},
		UserType: sql.NullString{
			String: string(enums.AuditableTypeUser),
			Valid:  true,
		},
		AuditableId: sql.NullInt64{
			Int64: event.UserId,
			Valid: true,
		},
		AuditableType: enums.AuditableTypeUser,
		Action:        enums.AuditLogActionLogin,
		Summary:       "user logged in",
	}

	return l.auditLogRepository.Create(ctx, log, nil)
}

func (l *UserLoggedInListener) TaskName() string {
	return events.EventUserLoggedIn
}

func (l *UserLoggedInListener) ShouldQueue() bool {
	return false
}
