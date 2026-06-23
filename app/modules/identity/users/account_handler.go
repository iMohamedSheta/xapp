package users

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/http/handler"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/inertia"
)

type AccountHandler struct {
	*handler.Handler
	userRepo *UserRepository
}

func NewAccountHandler(base *handler.Handler, userRepo *UserRepository) *AccountHandler {
	return &AccountHandler{
		Handler:  base,
		userRepo: userRepo,
	}
}

func (h *AccountHandler) ProfileView(c *gin.Context) error {
	return h.Inertia.Render(c, "Settings/Profile", inertia.Props{
		"mustVerifyEmail": false,
	})
}

func (h *AccountHandler) PasswordView(c *gin.Context) error {
	return h.Inertia.Render(c, "Settings/Password", nil)
}

func (h *AccountHandler) UpdateProfile(c *gin.Context) error {
	var req UpdateProfileRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		var xe *xerr.XErr
		if errors.As(err, &xe) && xe.IsType(enums.XErrValidationError) {
			return nil
		}
		return err
	}

	user, xerrVal := x.AuthUser(c)
	if xerrVal != nil {
		return xerrVal
	}

	fields := map[string]any{
		"name":     req.Name,
		"username": req.Username,
		"email":    req.Email,
	}

	if err := h.userRepo.Update(c, user.Id, fields, nil); err != nil {
		if h.HandleValidationErrors(c, err) {
			return nil
		}
		return err
	}

	return h.BackWithFlash(c, "تم تحديث الملف الشخصي بنجاح", 200)
}

func (h *AccountHandler) UpdatePassword(c *gin.Context) error {
	var req UpdatePasswordRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		var xe *xerr.XErr
		if errors.As(err, &xe) && xe.IsType(enums.XErrValidationError) {
			return nil
		}
		return err
	}

	authUser, xerrVal := x.AuthUser(c)
	if xerrVal != nil {
		return xerrVal
	}

	// Fetch user because AuthUser might be a snapshot
	user, err := h.userRepo.FindById(c, nil, authUser.Id)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(req.CurrentPassword, string(user.Password)) {
		err := xerr.New("كلمة المرور الحالية غير صحيحة", enums.XErrValidationError, nil).
			WithDetails(map[string]any{"current_password": "كلمة المرور الحالية غير صحيحة"})
		if h.HandleValidationErrors(c, err) {
			return nil
		}
		return err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	if err := h.userRepo.Update(c, user.Id, map[string]any{"password": hashedPassword}, nil); err != nil {
		if h.HandleValidationErrors(c, err) {
			return nil
		}
		return err
	}

	return h.BackWithFlash(c, "تم تحديث كلمة المرور بنجاح", 200)
}
