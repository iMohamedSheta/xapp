package notifications

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

type TestNotification struct{}

func NewTestNotification() *TestNotification {
	return &TestNotification{}
}

func (n *TestNotification) Type() string {
	return "TestNotification"
}

func (n *TestNotification) Channels() []string {
	return []string{"database", "websocket"}
}

func (n *TestNotification) ShouldQueue() bool {
	return true
}

func (n *TestNotification) ScheduledAt() *time.Time {
	return nil
}

func (n *TestNotification) Data(channel string, notifiableType string, notifiableId int64) map[string]any {
	t := time.Now().Format(time.RFC3339)

	// Database channel data
	if channel == "database" {
		return map[string]any{
			"title":             "تم انشاء طلب جديد",
			"message":           "تم استلام طلب تذاكر جديد بواسطة مستخدم بالرجاء مراجعة الطلب!",
			"notification_type": "notification",
			"link":              "",
			"timestamp":         t,
		}
	}

	// Websocket channel data
	if channel == "websocket" {
		return map[string]any{
			"channel": fmt.Sprintf("user_notifications.%d", notifiableId),
			"payload": map[string]any{
				"id":      time.Now().UnixMilli(),
				"title":   "تم انشاء طلب جديد",
				"message": "تم استلام طلب تذاكر جديد بواسطة مستخدم بالرجاء مراجعة الطلب!",
				"time":    t,
				"type":    "notification",
				"read":    false,
			},
		}
	}

	return nil
}

func (n *TestNotification) QueueOpts(channel string) []asynq.Option {
	if channel == "websocket" {
		return []asynq.Option{
			asynq.Queue("ws_notifications"),
		}
	}
	return nil
}
