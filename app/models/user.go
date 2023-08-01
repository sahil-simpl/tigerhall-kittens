package models

import "time"

type User struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id" validate:"required"`
	Username  string    `json:"username" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

func (user *User) PopulateData(request *CreateUserRequest) {
	user.Email = request.Email
	user.Password = request.Password
	user.Username = request.Username
}
