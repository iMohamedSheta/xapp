package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"
	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/modules/identity/tenants"
	"github.com/imohamedsheta/xapp/app/modules/identity/users"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/events"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/config"
	"github.com/imohamedsheta/xsocial"
)

var (
	ErrUserBlocked           = errors.New("user_blocked")
	ErrUserEmailAlreadyTaken = errors.New("email_already_taken")
)

// AuthAction orchestrates authentication use cases
type AuthAction struct {
	authService   *AuthService
	tenantService *tenants.TenantService
	userService   *users.UserService
}

func NewAuthAction(
	authService *AuthService,
	tenantService *tenants.TenantService,
	userService *users.UserService,
) *AuthAction {
	return &AuthAction{
		authService:   authService,
		tenantService: tenantService,
		userService:   userService,
	}
}

// // Login orchestrates authentication and token generation
func (a *AuthAction) Login(c context.Context, email, username, password string) (user *models.User, accessToken, refreshToken string, err error) {
	user, accessToken, refreshToken, err = a.authService.Authenticate(c, email, username, password)
	if err != nil {
		return nil, "", "", err
	}

	if user.IsBlockedFromLogin() {
		return nil, "", "", xerr.New("User is blocked from login", enums.XErrUnAuthorizedError, nil).
			WithPublicMessage("انت محظور من تسجيل الدخول يرجى التواصل مع الدعم")
	}

	if err := x.EventBus().Publish(c, events.EventUserLoggedIn, &events.UserLoggedInPayload{
		UserId:        user.Id,
		AuditableType: enums.AuditableTypeUser,
		Summary:       fmt.Sprintf("تم تسجيل الدخول كمستخدم: %s", user.Name),
	}); err != nil {
		x.Logger(config.EventBusLog).Error("publish user.logged_in failed: " + err.Error())
	}

	return user, accessToken, refreshToken, nil
}

// Register orchestrates user creation and token issuance
func (a *AuthAction) Register(ctx context.Context, req *RegisterRequest) (user *models.User, accessToken string, refreshToken string, err error) {
	err = xqb.Transaction(func(tx *sql.Tx) error {

		tenant, err := a.tenantService.Create(ctx, &models.Tenant{
			Name:   req.TenantName,
			Status: enums.TenantStatusActive,
		}, tx)
		if err != nil {
			return err
		}

		user, err = a.userService.Create(ctx, &models.User{
			Username: sql.NullString{String: req.Username, Valid: true},
			Email:    sql.NullString{String: req.Email, Valid: true},
			Password: req.Password,
			Name:     req.Name,
			Role:     enums.RoleSuperAdmin,
			Status:   enums.UserStatusActive,
			TenantId: tenant.Id,
		}, tx)
		if err != nil {
			return err
		}

		user.Tenant = tenant

		return nil
	})

	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err = a.authService.IssueTokens(ctx, user.Id, nil, enums.UserRole(user.Role))
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// LoginById orchestrates authentication by user ID and token generation
func (a *AuthAction) LoginById(c context.Context, userId int64, impersonatedBy *int64) (user *models.User, accessToken, refreshToken string, err error) {
	user, accessToken, refreshToken, err = a.authService.AuthenticateById(c, userId, impersonatedBy)
	if err != nil {
		return nil, "", "", err
	}

	if user.Role == enums.RoleSuperManager || user.Role == enums.RoleManager {
		return nil, "", "", xerr.New("forrbidden can't login by id to manager user", enums.XErrForbiddenError, nil)
	}

	if err := x.EventBus().Publish(c, events.EventUserLoggedIn, &events.UserLoggedInPayload{
		UserId:          user.Id,
		AuditableType:   enums.AuditableTypeUser,
		ImpersionatedBy: impersonatedBy,
		Summary:         fmt.Sprintf("تم تسجيل الدخول كمستخدم: %s", user.Name),
	}); err != nil {
		x.Logger(config.EventBusLog).Error("publish user.logged_in failed: " + err.Error())
	}

	return user, accessToken, refreshToken, nil
}

// OAuthRegister register user with oauth provider
func (a *AuthAction) OAuthRegisterOrLogin(ctx context.Context, provider string, socialiteUser *xsocial.User) (user *models.User, accessToken string, refreshToken string, err error) {
	user, err = a.userService.FindByOAuthProvider(ctx, provider, fmt.Sprintf("%v", socialiteUser.ID))
	if err != nil && !errors.Is(err, xqb.ErrNotFound) {
		return nil, "", "", err
	}

	if user != nil {
		if user.IsBlockedFromLogin() {
			return nil, "", "", ErrUserBlocked
		}

		return a.LoginById(ctx, user.Id, nil)
	}

	existingUser, err := a.userService.FindUserByEmail(ctx, socialiteUser.Email)
	if err != nil && !errors.Is(err, xqb.ErrNotFound) {
		return nil, "", "", err
	}

	if existingUser != nil {
		return nil, "", "", ErrUserEmailAlreadyTaken
	}

	err = xqb.Transaction(func(tx *sql.Tx) error {

		nickname := socialiteUser.Nickname

		if nickname == "" {
			nickname = strings.Split(socialiteUser.Name, " ")[0]
		}

		tenant, err := a.tenantService.Create(ctx, &models.Tenant{
			Name:   nickname,
			Status: enums.TenantStatusActive,
		}, tx)
		if err != nil {
			return err
		}

		user, err = a.userService.Create(ctx, &models.User{
			Email:      sql.NullString{String: socialiteUser.Email, Valid: true},
			Password:   "",
			Name:       socialiteUser.Name,
			Role:       enums.RoleSuperAdmin,
			Status:     enums.UserStatusActive,
			Provider:   sql.NullString{String: provider, Valid: true},
			ProviderId: sql.NullString{String: fmt.Sprintf("%v", socialiteUser.ID), Valid: true},
			EmailVerifiedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			TenantId: tenant.Id,
		}, tx)
		if err != nil {
			xqb.Dump(err)
			return err
		}

		user.Tenant = tenant

		return nil
	})

	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err = a.authService.IssueTokens(ctx, user.Id, nil, enums.UserRole(user.Role))
	if err != nil {
		return nil, "", "", err
	}

	if err := x.EventBus().Publish(ctx, events.EventUserRegister, &events.UserRegisterPayload{
		UserId:        user.Id,
		AuditableType: enums.AuditableTypeUser,
		Summary:       fmt.Sprintf("تم تسجيل مستخدم جديد %s من خلال oauth provider: %s", socialiteUser.Name, provider),
	}); err != nil {
		x.Logger(config.EventBusLog).Error("publish user.register failed: " + err.Error())
	}

	return user, accessToken, refreshToken, nil
}

func (a *AuthAction) AuthenticateByIdAndPassword(c context.Context, userId int64, password string, impersonatedBy *int64) (user *models.User, accessToken, refreshToken string, err error) {
	user, accessToken, refreshToken, err = a.authService.AuthenticateByIdAndPassword(c, userId, password, impersonatedBy)
	if err != nil {
		return nil, "", "", err
	}
	return user, accessToken, refreshToken, nil
}
