package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/models"
)

type NotificationListFilters struct {
	Page    int
	PerPage int
	Unread  bool
}

type NotificationRepository struct {
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) UnreadCount(ctx context.Context, userId int64) (int64, error) {
	return xqb.Table("notifications").
		WithContext(ctx).
		WhereNull("read_at").
		Where("notifiable_type", "=", models.User{}.GetNotifiableType()).
		Where("notifiable_id", "=", userId).
		Count("id")
}

func (r *NotificationRepository) ListUnreadForUser(ctx context.Context, userId int64, limit int) ([]models.Notification, error) {
	if limit <= 0 {
		limit = 100
	}
	return xqb.Model[models.Notification]().
		WithContext(ctx).
		WhereNull("read_at").
		Where("notifiable_type", "=", models.User{}.GetNotifiableType()).
		Where("notifiable_id", "=", userId).
		OrderByDesc("id").
		Limit(limit).
		Get()
}

func (r *NotificationRepository) Create(
	ctx context.Context,
	notificationType string,
	data map[string]any,
	notifiableID int64,
	notifiableType string,
) error {
	now := time.Now()
	if err := xqb.Table("notifications").WithContext(ctx).Insert([]map[string]any{
		{
			"type":            notificationType,
			"data":            data,
			"notifiable_id":   notifiableID,
			"notifiable_type": notifiableType,
			"created_at":      now,
			"updated_at":      now,
		},
	}); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) ReadNotification(c *gin.Context, userId int64, notificationId int64) error {
	_, err := xqb.Table("notifications").
		WithContext(c).
		Where("id", "=", notificationId).
		Where("notifiable_id", "=", userId).
		Where("notifiable_type", "=", models.User{}.GetNotifiableType()).
		Update(map[string]any{
			"read_at": time.Now(),
		})
	return err
}

func (r *NotificationRepository) ListForUser(c *gin.Context, userId int64, filters NotificationListFilters) ([]models.Notification, map[string]any, error) {
	if filters.PerPage <= 0 {
		filters.PerPage = 25
	}
	if filters.Page <= 0 {
		filters.Page = 1
	}

	q := xqb.Model[models.Notification]().
		WithContext(c).
		Where("notifiable_id", "=", userId).
		Where("notifiable_type", "=", models.User{}.GetNotifiableType()).
		OrderByDesc("id")

	if filters.Unread {
		q = q.WhereNull("read_at")
	}

	return q.Paginate(filters.PerPage, filters.Page, "notifications.id")
}

func (r *NotificationRepository) MarkAllRead(c *gin.Context, userId int64) error {
	_, err := xqb.Table("notifications").
		WithContext(c).
		Where("notifiable_id", "=", userId).
		Where("notifiable_type", "=", models.User{}.GetNotifiableType()).
		WhereNull("read_at").
		Update(map[string]any{
			"read_at": time.Now(),
		})
	return err
}

func (r *NotificationRepository) OpenNotification(c *gin.Context, userId int64, notificationId int64) error {
	_, err := xqb.Table("notifications").
		WithContext(c).
		Where("id", "=", notificationId).
		Where("notifiable_id", "=", userId).
		Where("notifiable_type", "=", models.User{}.GetNotifiableType()).
		Update(map[string]any{
			"opened_at": time.Now(),
		})
	return err
}
