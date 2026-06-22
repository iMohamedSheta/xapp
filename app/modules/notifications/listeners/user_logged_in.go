package listeners

import (
	"context"
	"encoding/json"

	"github.com/imohamedsheta/xapp/app/domain/events"
	"github.com/imohamedsheta/xapp/app/domain/notifications"
	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xnotify"
)

type UserRepository interface {
	ListManagers(ctx context.Context) ([]models.User, error)
}

type UserLoggedInListener struct {
	notify         *xnotify.Notify
	userRepository UserRepository
}

func NewUserLoggedInListener(notify *xnotify.Notify, userRepository UserRepository) *UserLoggedInListener {
	return &UserLoggedInListener{
		notify:         notify,
		userRepository: userRepository,
	}
}

func (l *UserLoggedInListener) Handle(ctx context.Context, payload []byte) error {
	var eventPayload events.UserLoggedInPayload
	if err := json.Unmarshal(payload, &eventPayload); err != nil {
		return err
	}

	notification := notifications.NewSystemEventNotification(
		"New user registered",
		"A new user has been registered in the system",
		"",
		"user_registered",
	)

	managers, err := l.userRepository.ListManagers(ctx)
	if err != nil {
		return err
	}

	notifiable := make([]xnotify.Notifiable, 0, len(managers))
	for _, manager := range managers {
		notifiable = append(notifiable, &manager)
	}

	return l.notify.Send(ctx, notification, notifiable...)
}

func (l *UserLoggedInListener) TaskName() string {
	return events.EventUserLoggedIn
}

func (l *UserLoggedInListener) ShouldQueue() bool {
	return false
}
