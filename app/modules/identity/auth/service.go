package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iMohamedSheta/xerr"
	"github.com/iMohamedSheta/xqb"

	"github.com/imohamedsheta/xapp/app/models"
	"github.com/imohamedsheta/xapp/app/modules/identity/users"
	"github.com/imohamedsheta/xapp/app/shared/enums"
	"github.com/imohamedsheta/xapp/app/shared/utils"
)

// Define the action group for the authentication actions.
type AuthService struct {
	jwtService *JwtService
	userRepo   *users.UserRepository
}

type JwtCommandsContract interface {
	GenerateAccessToken(ctx context.Context, userID int64, role enums.UserRole) (string, error)
	GenerateRefreshToken(ctx context.Context, userID int64, role enums.UserRole) (string, error)
	ValidateAuthToken(tokenType enums.JwtTokenType, token string) (*MyClaims, error)
}

func NewAuthService(jwtService *JwtService, userRepo *users.UserRepository) *AuthService {
	return &AuthService{jwtService: jwtService, userRepo: userRepo}
}

// Authenticate user with given credentials and issue tokens
func (a *AuthService) Authenticate(c context.Context, email string, username string, password string) (*models.User, string, string, error) {
	user, err := a.AuthenticateUserCredentials(c, email, username, password)
	if err != nil {
		return nil, "", "", err
	}

	access_token, refresh_token, err := a.IssueTokens(c, user.Id, nil, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	return user, access_token, refresh_token, nil
}

func (a *AuthService) AuthenticateById(c context.Context, userId int64, impersonatedBy *int64) (*models.User, string, string, error) {
	user, err := a.userRepo.FindById(c, nil, userId)
	if err != nil {
		if errors.Is(err, xqb.ErrNotFound) {
			return nil, "", "", xerr.New("User not found", enums.XErrNotFoundError, err)
		}
		return nil, "", "", err
	}

	access_token, refresh_token, err := a.IssueTokens(c, userId, impersonatedBy, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	return user, access_token, refresh_token, nil
}

func (a *AuthService) AuthenticateByIdAndPassword(c context.Context, userId int64, password string, impersonatedBy *int64) (*models.User, string, string, error) {
	user, err := a.AuthenticateUserByIdAndPassword(c, userId, password)
	if err != nil {
		return nil, "", "", err
	}

	access_token, refresh_token, err := a.IssueTokens(c, userId, impersonatedBy, user.Role)
	if err != nil {
		return nil, "", "", err
	}

	return user, access_token, refresh_token, nil
}

// Issue tokens for the given user id and role access and refresh tokens
func (a *AuthService) IssueTokens(c context.Context, userId int64, impersonatedBy *int64, role enums.UserRole) (accessToken string, refreshToken string, e error) {
	sessionId := uuid.New().String()
	access_token, err := a.jwtService.GenerateAccessToken(c, sessionId, userId, impersonatedBy, role)
	if err != nil {
		return "", "", err
	}
	refresh_token, err := a.jwtService.GenerateRefreshToken(c, sessionId, userId, impersonatedBy, role)
	if err != nil {
		return "", "", err
	}

	return access_token, refresh_token, nil
}

// Check if user can login with given credentials (email or username)
func (a *AuthService) AuthenticateUserCredentials(c context.Context, email string, username string, password string) (*models.User, error) {
	var user *models.User
	var err error

	// Try to find user by email if provided
	if email != "" {
		user, err = a.userRepo.FindUserByEmail(c, email)
		if err != nil {
			if errors.Is(err, xqb.ErrNotFound) {
				return nil, xerr.New("invalid credentials", enums.XErrValidationError, nil).WithDetails(map[string]any{
					"email": "invalid credentials",
				})
			}
			return nil, err
		}
	} else if username != "" {
		// Try to find user by username if provided
		user, err = a.userRepo.FindUserByUsername(c, username)
		if err != nil {
			if errors.Is(err, xqb.ErrNotFound) {
				return nil, xerr.New("invalid credentials", enums.XErrValidationError, nil).WithDetails(map[string]any{
					"username": "invalid credentials",
				})
			}
			return nil, err
		}
	} else {
		return nil, xerr.New("email or username required", enums.XErrValidationError, nil).WithDetails(map[string]any{
			"email": "Please provide either email or username",
		})
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		fieldName := "email"
		if username != "" {
			fieldName = "username"
		}
		return nil, xerr.New("invalid credentials", enums.XErrValidationError, nil).WithDetails(map[string]any{
			fieldName: "invalid credentials",
		})
	}

	return user, nil
}

func (a *AuthService) AuthenticateUserByIdAndPassword(c context.Context, id int64, password string) (*models.User, error) {
	user, err := a.userRepo.FindById(c, nil, id)
	if err != nil {
		if errors.Is(err, xqb.ErrNotFound) {
			return nil, xerr.New("invalid credentials", enums.XErrValidationError, nil).WithDetails(map[string]any{
				"password": "invalid credentials",
			})
		}

		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, xerr.New("invalid credentials", enums.XErrValidationError, nil).WithDetails(map[string]any{
			"password": "invalid credentials",
		})
	}

	return user, nil
}
