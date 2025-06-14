package users

import (
	"context"
	"database/sql"
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
	InvalidateToken(ctx context.Context, arg repository.CreateExpiredTokenParams) error
	CheckIfTokenExists(ctx context.Context, arg repository.CheckIfTokenExistsParams) (int64, error)
}

type usersRepository struct {
	repo repository.Querier
}

func (u *usersRepository) CheckIfTokenExists(ctx context.Context, arg repository.CheckIfTokenExistsParams) (int64, error) {
	exists, err := u.repo.CheckIfTokenExists(ctx, arg)
	if err != nil {
		return 0, fmt.Errorf("error checking if token exists: %w", err)
	}
	if exists == 0 {
		return 0, sql.ErrNoRows
	}

	return exists, nil
}

func (u *usersRepository) InvalidateToken(ctx context.Context, arg repository.CreateExpiredTokenParams) error {
	rows, err := u.repo.CreateExpiredToken(ctx, arg)
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("token not created")
	}

	return nil
}

func (u *usersRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	user, err := u.repo.GetByEmail(ctx, email)

	if err != nil {
		return User{}, err
	}

	return newUser(user), nil
}

func (u *usersRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := u.repo.EmailExists(ctx, email)
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

func (u *usersRepository) UpdateUser(ctx context.Context, arg repository.UpdateUserParams) error {
	rows, err := u.repo.UpdateUser(ctx, arg)
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (u *usersRepository) GetByUserId(ctx context.Context, userId string) (User, error) {
	user, err := u.repo.GetByUserId(ctx, userId)

	if err != nil {
		return User{}, err
	}

	return newUser(user), nil
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
		ID:         v.ID,
		Username:   v.Username,
		Email:      v.Email,
		Password:   v.Password,
		CreatedOn:  v.CreatedOn,
		UpdatedOn:  v.UpdatedOn,
		IsVerified: v.IsVerified,
	}

	return user
}

func NewRepository(repo repository.Querier) UsersRepository {
	return &usersRepository{repo: repo}
}
