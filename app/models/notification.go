package models

import (
	"database/sql"
	"fmt"

	"time"
)

/*
* Notification model matches the `notifications` table
 */
type Notification struct {
	Id             int64             `xqb:"id" json:"id"`
	TenantId       int64             `xqb:"tenant_id" json:"tenant_id"`
	NotifiableId   int64             `xqb:"notifiable_id" json:"notifiable_id"`
	NotifiableType string            `xqb:"notifiable_type" json:"notifiable_type"`
	Type           string            `xqb:"type" json:"type"`
	Data           *NotificationData `xqb:"data" json:"data"` // JSON field
	ReadAt         sql.NullTime      `xqb:"read_at" json:"read_at"`
	OpenedAt       sql.NullTime      `xqb:"opened_at" json:"opened_at"`
	CreatedAt      time.Time         `xqb:"created_at" json:"created_at"`
	UpdatedAt      time.Time         `xqb:"updated_at" json:"updated_at"`
}

/*
* NotificationData struct matches the `notification_data` column in the `notifications` table
 */
type NotificationData struct {
	Title   string `xqb:"title" json:"title"`
	Message string `xqb:"message" json:"message"`
	Link    string `xqb:"link" json:"link"`
	Type    string `xqb:"type" json:"type"`
}

func (Notification) Table() string {
	return "notifications"
}

func (n Notification) Cols() []any {
	return []any{
		n.Id,
		n.NotifiableId,
		n.NotifiableType,
		n.Type,
		n.Data,
		n.ReadAt,
		n.CreatedAt,
		n.UpdatedAt,
	}
}

type FrontendNotificationDTO struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Link    string `json:"link"`
	Type    string `json:"type"`
	Time    string `json:"time"`
	Read    bool   `json:"read"`
	Open    bool   `json:"open"`
}

/*
* ToFrontend converts the Notification model to a FrontendNotification DTO
 */
func (n Notification) ToFrontend() FrontendNotificationDTO {
	title := ""
	message := ""
	link := ""
	displayType := "notification"

	if n.Data != nil {
		title = n.Data.Title
		message = n.Data.Message
		link = n.Data.Link
		if n.Data.Type != "" {
			displayType = n.Data.Type
		}
	}
	t := n.CreatedAt.Format(time.RFC3339)

	return FrontendNotificationDTO{
		Id:      fmt.Sprintf("%d", n.Id),
		Title:   title,
		Message: message,
		Link:    link,
		Time:    t,
		Read:    n.ReadAt.Valid,
		Type:    displayType,
		Open:    n.OpenedAt.Valid,
	}
}
