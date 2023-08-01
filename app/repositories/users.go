package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tigerhallKittens/app/lib/db"
	"tigerhallKittens/app/models"
)

// UsersRepository interface is to interact with various user level db operations.
//
//go:generate mockgen -source=./users.go -destination=./mock_repositories/mock_users.go -package=mock_repositories
type UsersRepository interface {
	CreateUser(ctx context.Context, request *models.User) (*models.User, error)
}

// usersRepository implements UsersRepository interface
type usersRepository struct {
	DB *gorm.DB
}

// NewUsersRepository instantiates and returns UsersRepository Instance with the DB connection.
func NewUsersRepository() UsersRepository {
	return &usersRepository{DB: db.Get()}
}

func (repo *usersRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if user.UserID == "" {
		user.UserID = uuid.New().String()
	}

	queryErr := repo.DB.WithContext(ctx).Create(user).Error
	if queryErr != nil {
		return nil, queryErr
	}

	return user, nil
}
