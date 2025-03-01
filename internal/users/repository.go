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
		GetByUsername(ctx context.Context, arg string) (User, error)
		CreateAndReturnId(ctx context.Context, arg repository.CreateUserAndReturnIdParams) (string, error)
}

type usersRepository struct {
		repo repository.Querier
}

func (u *usersRepository) CreateAndReturnId(ctx context.Context, arg repository.CreateUserAndReturnIdParams) (string, error) {
		id, err := u.repo.CreateUserAndReturnId(ctx, arg)

		if err != nil {
				return "", err
		}

		return id, nil
}

func (u *usersRepository) GetByUsername(ctx context.Context, username string) (User, error) {
		user, err := u.repo.GetByUsername(ctx, username)

		if err != nil {
				return User{}, err
		}

		return newUser(user), nil
}

func newUser(v repository.User) User {
		user := User{
				ID:        v.ID,
				Username:  v.Username,
				Password:  v.Password,
				CreatedOn: v.CreatedOn,
				UpdatedOn: v.UpdatedOn,
		}

		return user
}

func NewRepository(repo repository.Querier) UsersRepository {
		return &usersRepository{repo: repo}
}
