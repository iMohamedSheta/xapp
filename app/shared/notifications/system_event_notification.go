package notifications

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

// SystemEventNotification is a generic in-app + websocket notification.
type SystemEventNotification struct {
	Title    string
	Message  string
	Link     string
	NType    string // warning, error, success, info, notification
	EventKey string // dedupe key for queue
}

func NewSystemEventNotification(title, message, link, nType string) *SystemEventNotification {
	if nType == "" {
		nType = "notification"
	}
	return &SystemEventNotification{
		Title:   title,
		Message: message,
		Link:    link,
		NType:   nType,
	}
}

func (n *SystemEventNotification) Type() string {
	if n.EventKey != "" {
		return "system_event:" + n.EventKey
	}
	return "system_event"
}

func (n *SystemEventNotification) Channels() []string {
	return []string{"database", "websocket"}
}

func (n *SystemEventNotification) ShouldQueue() bool {
	return true
}

func (n *SystemEventNotification) ScheduledAt() *time.Time {
	return nil
}

func (n *SystemEventNotification) Data(channel string, notifiableType string, notifiableId int64) map[string]any {
	t := time.Now().Format(time.RFC3339)
	ntype := n.NType
	if ntype == "" {
		ntype = "notification"
	}

	if channel == "database" {
		return map[string]any{
			"title":             n.Title,
			"message":           n.Message,
			"notification_type": ntype,
			"link":              n.Link,
			"timestamp":         t,
		}
	}

	if channel == "websocket" {
		return map[string]any{
			"channel": fmt.Sprintf("user_notifications.%d", notifiableId),
			"payload": map[string]any{
				"id":      time.Now().UnixMilli(),
				"title":   n.Title,
				"message": n.Message,
				"time":    t,
				"link":    n.Link,
				"type":    ntype,
				"read":    false,
			},
		}
	}

	return nil
}

// QueueOpts defines queue options for the notification.
func (n *SystemEventNotification) QueueOpts(channel string) []asynq.Option {
	if channel == "websocket" {
		return []asynq.Option{asynq.Queue("ws_notifications")}
	}
	if n.EventKey != "" {
		return []asynq.Option{asynq.Unique(2 * time.Minute)}
	}
	return nil
}
