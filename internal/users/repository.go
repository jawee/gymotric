package users

import (
	"context"
	"fmt"
	"weight-tracker/internal/repository"
)

type User struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	CreatedOn  string `json:"created_on"`
	UpdatedOn  string `json:"updated_on"`
	Email      any    `json:"email"`
	IsVerified bool   `json:"is_verified"`
}

type UsersRepository interface {
	GetByUsername(ctx context.Context, arg string) (User, error)
	CreateAndReturnId(ctx context.Context, arg repository.CreateUserAndReturnIdParams) (string, error)
	GetByUserId(ctx context.Context, userId string) (User, error)
	UpdateUser(ctx context.Context, arg repository.UpdateUserParams) error
	EmailExists(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (User, error)
}

type usersRepository struct {
	repo repository.Querier
}

func (u *usersRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	user, err := u.repo.GetByEmail(ctx, email)

	if err != nil {
		return User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return newUser(user), nil
}

func (u *usersRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := u.repo.EmailExists(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}

	return exists > 0, nil
}

func (u *usersRepository) UpdateUser(ctx context.Context, arg repository.UpdateUserParams) error {
	rows, err := u.repo.UpdateUser(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (u *usersRepository) GetByUserId(ctx context.Context, userId string) (User, error) {
	user, err := u.repo.GetByUserId(ctx, userId)

	if err != nil {
		return User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return newUser(user), nil
}

func (u *usersRepository) CreateAndReturnId(ctx context.Context, arg repository.CreateUserAndReturnIdParams) (string, error) {
	id, err := u.repo.CreateUserAndReturnId(ctx, arg)

	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (u *usersRepository) GetByUsername(ctx context.Context, username string) (User, error) {
	user, err := u.repo.GetByUsername(ctx, username)

	if err != nil {
		return User{}, fmt.Errorf("failed to get user by username: %w", err)
	}

	return newUser(user), nil
}

func newUser(v repository.User) User {
	user := User{
		ID:        v.ID,
		Username:  v.Username,
		Email:     v.Email,
		Password:  v.Password,
		CreatedOn: v.CreatedOn,
		UpdatedOn: v.UpdatedOn,
		IsVerified: v.IsVerified,
	}

	return user
}

func NewRepository(repo repository.Querier) UsersRepository {
	return &usersRepository{repo: repo}
}
