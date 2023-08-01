package services

import (
	"context"
	"tigerhallKittens/app/domains"

	"tigerhallKittens/app/models"
)

//go:generate mockgen -source=./user_service.go -destination=./mock_services/mock_user_service.go -package=mock_services
type UserService interface {
	CreateUser(ctx context.Context, request *models.CreateUserRequest) (*models.User, error)
}

type userService struct {
	userManager domains.UserManager
}

func NewUserService() UserService {
	return &userService{
		userManager: domains.NewUserManager(),
	}
}

func (us userService) CreateUser(ctx context.Context, request *models.CreateUserRequest) (*models.User, error) {
	return us.userManager.CreateUser(ctx, request)
}
