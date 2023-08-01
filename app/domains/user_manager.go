package domains

import (
	"context"
	"tigerhallKittens/app/repositories"

	"tigerhallKittens/app/models"
)

// UserManager is the manager used to interact with various user related methods.
//
//go:generate mockgen -source=./user_manager.go -destination=./mock_domains/mock_user_manager.go -package=mock_domains
type UserManager interface {
	CreateUser(ctx context.Context, request *models.CreateUserRequest) (*models.User, error)
}

type userManager struct {
	repo repositories.UsersRepository
}

// NewUserManager creates a new user manager
func NewUserManager() UserManager {
	return &userManager{
		repo: repositories.NewUsersRepository(),
	}
}

// CreateUser creates a new user.
func (um *userManager) CreateUser(ctx context.Context, request *models.CreateUserRequest) (*models.User, error) {
	var createUserReq = new(models.User)
	createUserReq.PopulateData(request)

	user, repoErr := um.repo.CreateUser(ctx, createUserReq)
	if repoErr != nil {
		return nil, repoErr
	}
	return user, nil
}
