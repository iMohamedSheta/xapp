package auth

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/imohamedsheta/xapp/app/http/middleware"
	"github.com/imohamedsheta/xapp/app/modules/settings"
	"github.com/imohamedsheta/xapp/app/x"

	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"

	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"
	"github.com/imohamedsheta/xapp/pkg/inertia"
)

type AuthMiddleware struct {
	jwt         *JwtService
	permissions *PermissionService
	settingRepo *settings.SettingRepository
}

func NewAuthMiddleware(jwt *JwtService, permissionService *PermissionService, settingRepo *settings.SettingRepository) *AuthMiddleware {
	return &AuthMiddleware{
		jwt:         jwt,
		permissions: permissionService,
		settingRepo: settingRepo,
	}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		isInertiaReq := isInertiaRequestOrBrowserVisit(c)
		// Try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
			tokenStr = after
		}

		// Fallback: try to get token from cookie
		if tokenStr == "" {
			var err error
			tokenStr, err = c.Cookie("access_token")
			if err != nil || tokenStr == "" {
				middleware.ErrorHandler(c, xerr.New("missing access token in cookie", enums.XErrUnAuthorizedError, err), isInertiaReq)
				return
			}
		}

		// Validate token
		claims, err := m.jwt.ValidateAuthToken(enums.AccessToken, tokenStr)
		if err != nil {
			middleware.ErrorHandler(c, xerr.New("Invalid access token in the request", enums.XErrUnAuthorizedError, err), isInertiaReq)
			return
		}

		// Extract user ID (subject)
		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			middleware.ErrorHandler(c, xerr.New("Invalid or missing user ID in token in context", enums.XErrUnAuthorizedError, err), isInertiaReq)
			return
		}

		// Extract user role
		role, err := claims.GetRole()
		if err != nil || role == "" {
			middleware.ErrorHandler(c, xerr.New("Invalid or missing role in token in context", enums.XErrForbiddenError, err), isInertiaReq)
			return
		}

		userIDInt, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			middleware.ErrorHandler(c, xerr.New("Invalid or missing user ID in token in context", enums.XErrServerError, err), isInertiaReq)
			return
		}

		user, err := xqb.Model[models.User]().WithContext(c).Where("id", "=", userIDInt).First()
		if err != nil {
			middleware.ErrorHandler(c, xerr.New(fmt.Sprintf("Failed to get user from database user_id:  %s", userID), enums.XErrServerError, err), isInertiaReq)
			return
		}

		if user.TenantId != 0 {
			if err = user.LoadTenant(c); err != nil {
				middleware.ErrorHandler(c, xerr.New("missing access token", enums.XErrUnAuthorizedError, err), false)
				c.Abort()
				return
			}
		}

		if user.IsBlockedFromLogin() {
			x.ClearAuthCookies(c)
			return
		}

		// Load user permissions
		var permissions map[string]bool
		if m.permissions != nil {
			perms, err := m.permissions.GetGlobalPermissions(c, userIDInt, role)
			if err != nil {
				// Log error but continue with empty permissions
				permissions = make(map[string]bool)
			} else {
				permissions = perms
			}
		} else {
			// If no permission service is set, use empty permissions
			permissions = make(map[string]bool)
		}

		c.Set(string(enums.ContextKeyAuthId), userID)
		c.Set(string(enums.ContextKeyRole), enums.UserRole(role))
		c.Set(string(enums.ContextKeyImpersonatorId), claims.ImpersonatedBy)
		c.Set(string(enums.ContextKeyAuthUser), user)
		c.Set(string(enums.ContextKeySessionId), claims.SessionId)
		c.Set(string(enums.ContextKeyPermissions), permissions)

		c.Next()
	}
}

func (m *AuthMiddleware) RedirectToDashboardIfAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// Try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
			tokenStr = after
		}

		// Fallback to cookie
		if tokenStr == "" {
			var err error
			tokenStr, err = c.Cookie("access_token")
			if err != nil || tokenStr == "" {
				// No token: this is a guest, allow access
				c.Next()
				return
			}
		}

		// Validate token
		claims, err := m.jwt.ValidateAuthToken(enums.AccessToken, tokenStr)
		if err != nil {
			// Invalid token: treat as guest
			c.Next()
			return
		}

		var redirectTo string

		switch claims.Role {
		case enums.RoleSuperManager, enums.RoleManager:
			redirectTo = "/manager/dashboard"
		case enums.RoleSuperAdmin, enums.RoleAdmin:
			redirectTo = "/dashboard"
		case enums.RoleClient:
			redirectTo = "/client/dashboard"
		default:
			redirectTo = "/" // optional fallback
		}

		if inertia.IsInertiaRequest(c.Request) {
			i := x.Inertia()
			i.Redirect(c.Writer, c.Request, redirectTo)
			c.Abort()
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, redirectTo)
		c.Abort()
	}
}

func (m *AuthMiddleware) Web() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// Try Authorization header
		authHeader := c.GetHeader("Authorization")
		if after, ok := strings.CutPrefix(authHeader, "Bearer "); ok {
			tokenStr = after
		}

		// Fallback: try cookie
		if tokenStr == "" {
			token, err := c.Cookie("access_token")
			if err == nil && token != "" {
				tokenStr = token
			}
		}

		// If no token → guest, just continue
		if tokenStr == "" {
			c.Next()
			return
		}

		// Validate token
		claims, err := m.jwt.ValidateAuthToken(enums.AccessToken, tokenStr)
		if err != nil {
			// invalid token → guest
			c.Next()
			return
		}

		// Get subject (user ID)
		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			c.Next()
			return
		}

		// Get role
		role, _ := claims.GetRole()
		if role != "" {
			c.Set(string(enums.ContextKeyRole), enums.UserRole(role))
		}

		// Load user from DB
		userIDInt, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			c.Next()
			return
		}

		var user *models.User

		// Load user from DB
		user, err = xqb.Model[models.User]().WithContext(c).Where("id", "=", userIDInt).First()
		if err != nil {
			// DB error or user not found → guest
			c.Next()
			return
		}

		// Put into context
		c.Set(string(enums.ContextKeyAuthId), userID)
		c.Set(string(enums.ContextKeyAuthUser), user)
		c.Set(string(enums.ContextKeySessionId), claims.SessionId)

		// Continue request
		c.Next()
	}
}

func (m *AuthMiddleware) WebSocketAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string
		var err error

		// 1️⃣ Try reading token from Sec-WebSocket-Protocol (for JS clients using subprotocol)
		webSocketProtocolHeader := c.GetHeader("Sec-WebSocket-Protocol")
		if webSocketProtocolHeader != "" {
			tokens := strings.Split(webSocketProtocolHeader, ",")
			if len(tokens) == 2 {
				header := strings.TrimSpace(tokens[0])
				authHeader, err := decodeBase64URL(strings.TrimSpace(tokens[1]))
				if err == nil && header == "Authorization" && strings.HasPrefix(authHeader, "Bearer ") {
					tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}
		}

		// 2️⃣ If not found, try cookie
		if tokenStr == "" {
			tokenStr, err = c.Cookie("access_token")
			if err != nil {
				middleware.ErrorHandler(c, xerr.New("missing access token", enums.XErrUnAuthorizedError, err), false)
				c.Abort()
				return
			}
		}

		// 3️⃣ Validate token
		claims, err := m.jwt.ValidateAuthToken(enums.AccessToken, tokenStr)
		if err != nil {
			middleware.ErrorHandler(c, xerr.New("missing access token", enums.XErrUnAuthorizedError, err), false)
			c.Abort()
			return
		}

		userID, err := claims.GetSubject()
		if err != nil || userID == "" {
			middleware.ErrorHandler(c, xerr.New("missing access token", enums.XErrUnAuthorizedError, err), false)
			c.Abort()
			return
		}

		role, err := claims.GetRole()
		if err != nil || role == "" {
			middleware.ErrorHandler(c, xerr.New("missing access token", enums.XErrUnAuthorizedError, err), false)
			c.Abort()
			return
		}

		userIDInt, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			middleware.ErrorHandler(c, xerr.New("missing access token", enums.XErrUnAuthorizedError, err), false)
			c.Abort()
			return
		}

		var user *models.User
		user, err = xqb.Model[models.User]().WithContext(c).Where("id", "=", userIDInt).First()
		if err != nil {
			middleware.ErrorHandler(c, xerr.New("missing access token", enums.XErrUnAuthorizedError, err), false)
			c.Abort()
			return
		}

		if user.TenantId != 0 {
			if err = user.LoadTenant(c); err != nil {
				middleware.ErrorHandler(c, xerr.New("missing access token", enums.XErrUnAuthorizedError, err), false)
				c.Abort()
				return
			}
		}

		// Load user permissions for WebSocket
		var permissions map[string]bool
		if m.permissions != nil {
			perms, err := m.permissions.GetGlobalPermissions(c, userIDInt, role)
			if err != nil {
				permissions = make(map[string]bool)
			} else {
				permissions = perms
			}
		} else {
			permissions = make(map[string]bool)
		}

		c.Set(string(enums.ContextKeyAuthId), userID)
		c.Set(string(enums.ContextKeyRole), enums.UserRole(role))
		c.Set(string(enums.ContextKeyAuthUser), user)
		c.Set(string(enums.ContextKeySessionId), claims.SessionId)
		c.Set(string(enums.ContextKeyImpersonatorId), claims.ImpersonatedBy)
		c.Set(string(enums.ContextKeyPermissions), permissions)

		c.Next()
	}
}

func (m *AuthMiddleware) SuperAdminOrAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsSuperAdminOrAdmin(c) {
			middleware.ErrorHandler(c, xerr.New("Unauthorized", enums.XErrUnAuthorizedError, nil), isInertiaRequestOrBrowserVisit(c))
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) SuperAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsSuperAdmin(c) {
			middleware.ErrorHandler(c, xerr.New("Unauthorized", enums.XErrUnAuthorizedError, nil), isInertiaRequestOrBrowserVisit(c))
			c.Abort()
			return
		}
		c.Next()
	}
}

func (m *AuthMiddleware) ClientOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsClient(c) {
			middleware.ErrorHandler(c, xerr.New("Unauthorized", enums.XErrUnAuthorizedError, nil), isInertiaRequestOrBrowserVisit(c))
			c.Abort()
			return
		}
		c.Next()
	}
}

func (m *AuthMiddleware) ManagerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsManager(c) {
			middleware.ErrorHandler(c, xerr.New("Unauthorized", enums.XErrUnAuthorizedError, nil), isInertiaRequestOrBrowserVisit(c))
			c.Abort()
			return
		}

		c.Next()
	}
}

func decodeBase64URL(encoded string) (string, error) {
	decodedBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func isInertiaRequestOrBrowserVisit(c *gin.Context) bool {
	acceptHeader := c.GetHeader("Accept")
	isBrowserVisit := strings.Contains(acceptHeader, "text/html")
	return inertia.IsInertiaRequest(c.Request) || isBrowserVisit
}
