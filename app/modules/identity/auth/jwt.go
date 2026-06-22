package auth

import (
	"context"

	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xfig"
)

// Represent the custom application claims
type MyClaims struct {
	Role           enums.UserRole `json:"role"`
	TokenType      string         `json:"token_type"`
	SessionId      string         `json:"session_id"`
	ImpersonatedBy *int64         `json:"impersonated_by,omitempty"` // Original super admin ID who is impersonating
	jwt.RegisteredClaims
}

// Additional helper to get user role from claims
func (c *MyClaims) GetRole() (enums.UserRole, error) {
	return c.Role, nil
}

type JwtService struct {
	config *xfig.Config
}

func NewJwtService(cfg *xfig.Config) *JwtService {
	return &JwtService{config: cfg}
}

// GenerateAccessToken generates a new access token for the given user
func (j *JwtService) GenerateAccessToken(c context.Context, sessionId string, userID int64, impersonatedBy *int64, role enums.UserRole) (string, error) {
	expiry := j.config.GetDuration("auth.jwt.access_token.expiry", 30*time.Minute)
	claims := j.buildClaims(userID, sessionId, impersonatedBy, role, enums.AccessToken, expiry)

	return j.signToken(claims)
}

// GenerateRefreshToken generates a new refresh token for the given user
func (j *JwtService) GenerateRefreshToken(c context.Context, sessionId string, userID int64, impersonatedBy *int64, role enums.UserRole) (string, error) {
	expiry := j.config.GetDuration("auth.jwt.refresh_token.expiry", 168*time.Hour)
	claims := j.buildClaims(userID, sessionId, impersonatedBy, role, enums.RefreshToken, expiry)

	return j.signToken(claims)
}

// ValidateAuthToken validates the given token and returns the claims if valid
func (j *JwtService) ValidateAuthToken(tokenType enums.JwtTokenType, tokenStr string) (*MyClaims, error) {
	expectedIssuer := j.config.GetString("auth.jwt.issuer", "TaskGo")
	expectedAudience := j.config.GetString("auth.jwt.audience", "TaskGoAudience")

	claims := &MyClaims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(expectedIssuer),
		jwt.WithAudience(expectedAudience),
		jwt.WithLeeway(5*time.Second),
	)

	token, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return j.getSecret()
	})

	if err != nil || !token.Valid {
		return nil, xerr.New("unauthenticated: invalid token", enums.XErrUnAuthorizedError, err)
	}

	if !j.isValidTokenType(tokenType, claims.TokenType) {
		return nil, xerr.New("unauthenticated: invalid token type", enums.XErrUnAuthorizedError, nil)
	}

	return claims, nil
}

// signToken signs the token with the secret.
func (j *JwtService) signToken(claims MyClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := j.getSecret()
	if err != nil {
		return "", err
	}
	return token.SignedString(secret)
}

// buildClaims builds the claims used to sign the token.
func (j *JwtService) buildClaims(userID int64, sessionId string, impersonatedBy *int64, role enums.UserRole, tokenType enums.JwtTokenType, expiry time.Duration) MyClaims {
	return MyClaims{
		Role:           role,
		TokenType:      string(tokenType),
		SessionId:      sessionId,
		ImpersonatedBy: impersonatedBy,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			Issuer:    j.config.GetString("auth.jwt.issuer", "TaskGo"),
			Audience:  []string{j.config.GetString("auth.jwt.audience", "TaskGoAudience")},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.FormatInt(userID, 10),
		},
	}
}

// isValidTokenType checks if the token type is valid or not.
func (j *JwtService) isValidTokenType(expected enums.JwtTokenType, actual string) bool {
	return actual == string(expected)
}

// getSecret returns the secret used to sign the token.
func (j *JwtService) getSecret() ([]byte, error) {
	return []byte(utils.GetAppSecret()), nil
}
