package hooks

import (
	"context"

	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/modules/notifications"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"

	"github.com/gin-gonic/gin"
	"github.com/imohamedsheta/xapp/pkg/inertia"
)

func InitInertiaHooks(i *inertia.Inertia, flash *inertia.InmemFlashProvider) {
	i.OnBeforeRender(func(ctx context.Context) {
		ginCtx, ok := ctx.(*gin.Context)
		if !ok || ginCtx == nil {
			logError("OnBeforeRender: ctx is not *gin.Context")
			return
		}

		reqID, _ := ctx.Value(enums.ContextKeyRequestId.String()).(string)
		// logError("Calling OnBeforeRender Inertia Hook")

		if reqID == "" {
			logError("Failed to get request id from context")
			return
		}

		// logError("Calling OnBeforeRender Inertia Hook")

		if reqID == "" {
			logError("Failed to get request id from context")
			return
		}

		attachFlashErrors(ginCtx, flash)
		attachUserProps(ginCtx)
		// Always run this code after the request is done

		if utils.IsDebug() {
			attachDebugProps(ctx, ginCtx, reqID)
		}
	})
}

// ---------------- Helpers ----------------
func attachUserProps(c *gin.Context) {
	user, xerr := x.AuthUser(c)
	if xerr != nil {
		// if the user is not logged in yet ignore logging the error of missing user
		if !xerr.IsType(enums.XErrUnAuthorizedError) {
			x.Logger().Error(xerr.Error())
		}

		x.Inertia().ShareProp("auth", map[string]any{
			"user":             nil,
			"is_impersonating": false,
			"permissions":      nil,
		})
		x.Inertia().ShareProp("notification", map[string]any{
			"count":         0,
			"notifications": []models.FrontendNotificationDTO{},
		})
		x.Inertia().ShareProp("settings", nil)
		return
	}

	if user != nil && utils.IsSuperAdminOrAdmin(c) {
		if err := user.LoadTenant(c); err != nil {
			x.Logger().Error(err.Error())
		}

	} else {
		x.Inertia().ShareProp("settings", nil)
	}

	is_impersonating := false
	if impersonatorId, exists := c.Get(string(enums.ContextKeyImpersonatorId)); exists && impersonatorId != nil {
		if id, ok := impersonatorId.(*int64); ok && id != nil && *id > 0 {
			is_impersonating = true
		}
	}

	notifRepo := notifications.NewNotificationRepository()
	count, err := notifRepo.UnreadCount(c, user.Id)
	if err != nil {
		x.Logger().Error(err.Error())
	}

	var notifications, errList = notifRepo.ListUnreadForUser(c, user.Id, 100)
	if errList != nil {
		x.Logger().Error(errList.Error())
		notifications = nil
	}

	var frontendNotifications []models.FrontendNotificationDTO
	for _, notification := range notifications {
		frontendNotifications = append(frontendNotifications, notification.ToFrontend())
	}

	x.Inertia().ShareProp("notification", map[string]any{
		"count":         count,
		"notifications": frontendNotifications,
	})
	var permissions map[string]bool
	if perms, exists := c.Get(string(enums.ContextKeyPermissions)); exists {
		if p, ok := perms.(map[string]bool); ok {
			permissions = p
		}
	}

	x.Inertia().ShareProp("auth", map[string]any{
		"user":             user,
		"is_impersonating": is_impersonating,
		"permissions":      permissions,
	})
}

func attachFlashErrors(c *gin.Context, flash *inertia.InmemFlashProvider) {
	errors, err := flash.GetErrors(c)
	if err != nil {
		x.Logger().Error(err.Error())
	}
	x.Inertia().ShareProp("errors", errors)

	props, err := flash.GetFlash(c)
	if err != nil {
		x.Logger().Error(err.Error())
	}
	x.Inertia().ShareProp("flash", props)
}

func logError(msg string) {
	utils.PrintErr(msg)
	x.Logger().Error(msg)
}
