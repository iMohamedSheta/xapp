package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"

	"github.com/imohamedsheta/xnotify"
	"github.com/imohamedsheta/xws"
)

// WebsocketChannelHandler is a handler for websocket channel
func WebsocketChannelHandler(ctx context.Context, task *xnotify.NotificationTask) error {
	data := task.Data
	websocket := x.WS()

	payload := data["payload"].(map[string]any)

	utils.PrintErr(payload)
	notificationData := models.FrontendNotificationDTO{
		Id:      fmt.Sprintf("%.0f", payload["id"].(float64)),
		Title:   payload["title"].(string),
		Message: payload["message"].(string),
		Type:    payload["type"].(string),
		Link:    payload["link"].(string),
		Time:    payload["time"].(string),
		Read:    payload["read"].(bool),
	}

	jsonBytes, _ := json.Marshal(notificationData)

	websocket.Hub.Broadcast(&xws.WSMessage{
		Type:    "notification",
		Channel: data["channel"].(string),
		Data:    jsonBytes,
		From:    "server",
	})
	return nil
}
