package users

import (
	"context"
	"weight-tracker/internal/repository"
)

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

type UsersRepository interface {
	GetByUsernameAndPassword(ctx context.Context, arg repository.GetUserByUsernameAndPasswordParams) (User, error)
	CreateAndReturnId(ctx context.Context, arg repository.CreateUserAndReturnIdParams) (string, error)
}

type usersRepository struct {
	repo repository.Querier
}
