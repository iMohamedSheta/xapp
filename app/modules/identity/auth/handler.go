package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/http/handler"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/pkg/inertia"
)

// Handler handles HTTP requests for the auth module.
// Add your handler methods here.
type AuthHandler struct {
	*handler.Handler
	authAction  *AuthAction
	permissions *PermissionService
}

func NewAuthHandler(
	base *handler.Handler,
	authAction *AuthAction,
	permissions *PermissionService,
) *AuthHandler {
	return &AuthHandler{
		Handler:     base,
		authAction:  authAction,
		permissions: permissions,
	}
}

func (h *AuthHandler) LoginView(c *gin.Context) error {
	return h.Inertia.Render(c, "Auth/Login", inertia.Props{
		"canResetPassword": false,
		"status":           "",
	})
}

func (h *AuthHandler) RegisterView(c *gin.Context) error {
	return h.Inertia.Render(c, "Auth/Register", nil)
}

func (h *AuthHandler) DashboardView(c *gin.Context) error {
	return h.Inertia.Render(c, "Dashboard", nil)
}

func (h *AuthHandler) Login(c *gin.Context) error {
	var req LoginRequest
	var err error

	if err = h.BindAndValidate(c, &req); err != nil {
		var xe *xerr.XErr
		if errors.As(err, &xe) && xe.IsType(enums.XErrValidationError) {
			return nil
		}
		return err
	}

	user, accessToken, refreshToken, err := h.authAction.Login(c, req.Email, req.Username, req.Password)
	if err != nil {
		if h.HandleValidationErrors(c, err) {
			return nil
		}
		return err
	}

	// Force clear permission cache for the user after login
	h.permissions.ClearUserPermissionCache(c, user.Id)

	if err = h.setAuthCookies(c, accessToken, refreshToken, ""); err != nil {
		return err
	}

	h.Inertia.Redirect(c.Writer, c.Request, "/dashboard", http.StatusFound)
	return nil
}

// Register a new user
func (h *AuthHandler) Register(c *gin.Context) error {
	var request RegisterRequest
	var err error

	if err = h.BindAndValidate(c, &request); err != nil {
		var xe *xerr.XErr
		if errors.As(err, &xe) && xe.IsType(enums.XErrValidationError) {
			return nil
		}
		return err
	}

	_, accessToken, refreshToken, err := h.authAction.Register(c, &request)
	if err != nil {
		return err
	}

	if err := h.setAuthCookies(c, accessToken, refreshToken, ""); err != nil {
		return err
	}

	h.Inertia.Redirect(c.Writer, c.Request, "/dashboard", http.StatusFound)
	return nil
}

func (h *AuthHandler) Logout(c *gin.Context) error {
	auth_user, xerror := User(c)
	if xerror != nil {
		return xerror
	}

	if auth_user == nil {
		return h.FlashBack(c, inertia.Props{
			"redirect_url": "/login",
		}, 200)
	}

	// Clear tokens
	ClearAuthCookies(c)

	return h.FlashBack(c, inertia.Props{
		"redirect_url": "/login",
	}, 200)
}

func (h *AuthHandler) setAuthCookies(c *gin.Context, accessToken, refreshToken string, domain string) error {
	cfg := x.Config()

	c.SetCookie(
		"access_token",
		accessToken,
		int(cfg.GetDuration("auth.jwt.access_token.expiry", 30*time.Minute)),
		cfg.GetString("auth.jwt.access_token.path", "/"),
		domain,
		false,
		true,
	)

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(cfg.GetDuration("auth.jwt.refresh_token.expiry", 160*time.Hour)),
		cfg.GetString("auth.jwt.refresh_token.path", "/auth/refresh"),
		domain,
		false,
		true,
	)

	return nil
}
