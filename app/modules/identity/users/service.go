package users

import (
	"context"
	"database/sql"

	"github.com/imohamedsheta/xapp/app/models"
)

type UserService struct {
	userRepo *UserRepository
}

func NewUserService(
	userRepo *UserRepository,
) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (a *UserService) Create(ctx context.Context, user *models.User, tx *sql.Tx) (*models.User, error) {
	return a.userRepo.Create(ctx, user, tx)
}

func (a *UserService) FindByOAuthProvider(ctx context.Context, provider string, providerId string) (*models.User, error) {
	return a.userRepo.FindByOAuthProvider(ctx, provider, providerId)
}

func (a *UserService) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return a.userRepo.FindUserByEmail(ctx, email)
}
