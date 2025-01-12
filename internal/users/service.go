package users

import (
	"context"
	"time"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateAndReturnId(ctx context.Context, arg createUserAndReturnIdRequest) (string, error)
}

type usersService struct {
	repo UsersRepository
}

// CreateUserAndReturnId implements Service.
func (u *usersService) CreateAndReturnId(ctx context.Context, arg createUserAndReturnIdRequest) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	pwBytes, err := bcrypt.GenerateFromPassword([]byte(arg.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := repository.CreateUserAndReturnIdParams{
		ID:        uuid.String(),
		Username:  arg.Username,
		Password:  string(pwBytes),
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
	}
	id, err := u.repo.CreateAndReturnId(ctx, user)
	return id, err
}

func NewService(repo UsersRepository) Service {
	return &usersService{repo}
}
