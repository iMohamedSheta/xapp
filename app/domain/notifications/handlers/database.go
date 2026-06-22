package handlers

import (
	"context"

	"github.com/imohamedsheta/xnotify"
)

func DatabaseChannelHandler(ctx context.Context, task *xnotify.NotificationTask) error {
	// repo := repository.NewNotificationRepository()
	// return repo.Create(ctx, task.NotificationType, task.Data, task.NotifiableID, task.NotifiableType)
	return nil
}
