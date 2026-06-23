package handlers

import (
	"context"

	"github.com/imohamedsheta/xapp/app/shared/enums"

	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xnotify"
)

// WhatsappChannelHandler is a handler for whatsapp channel
func WhatsappChannelHandler(ctx context.Context, task *xnotify.NotificationTask) error {
	// data := task.Data
	// phone := data["phone"].(string)
	// payload := data["payload"].(map[string]any)

	// notificationData := dto.FrontendNotification{
	// 	Message: payload["message"].(string),
	// }

	// whatsappAction := actions.NewWhatsappService(x.WhatsappManager())

	// user, err := xqb.Model[models.User]().
	// 	WithContext(ctx).
	// 	Where("id", "=", task.NotifiableID).
	// 	First()
	// if err != nil {
	// 	return err
	// }

	// return whatsappAction.SendMessage(ctx, user.AgencyId, phone, notificationData.Message)
	return xerr.New("not implemented", enums.XErrServerError, nil)
}
