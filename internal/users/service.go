package users

import (
	"context"
	"os"
	"time"
	"weight-tracker/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateAndReturnId(ctx context.Context, arg createUserAndReturnIdRequest) (string, error)
	Login(ctx context.Context, arg loginRequest) (string, error)
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

func (u *usersService) Login(ctx context.Context, arg loginRequest) (string, error) {
	signingKey := os.Getenv("JWT_SIGN_KEY")
	mySigningKey := []byte(signingKey)
	user, err := u.repo.GetByUsername(ctx, arg.Username)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(arg.Password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * 15)),
		Issuer:    "weight-tracker",
		Subject:   user.ID,
		Audience:  []string{"weight-tracker"},
	})
	return token.SignedString(mySigningKey)
}

func NewService(repo UsersRepository) Service {
	return &usersService{repo}
}
