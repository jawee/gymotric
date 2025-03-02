package users

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"
	"weight-tracker/internal/repository"
	"weight-tracker/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
)

type loginResponse struct {
	Token string
	UserId string
}

type Service interface {
	CreateAndReturnId(ctx context.Context, arg createUserAndReturnIdRequest) (string, error)
	Login(ctx context.Context, arg loginRequest) (loginResponse, error)
	CreateToken(userId string) (string, error)
}

type usersService struct {
	repo UsersRepository
}

func (u *usersService) CreateToken(userId string) (string, error) {
	signingKey := os.Getenv(utils.EnvJwtSignKey)
	tokenExpiration, err := strconv.Atoi(os.Getenv(utils.EnvJwtExpireMinutes))
	if err != nil {
		slog.Error("Failed to convert JWT_EXPIRATION to int", "error", err)
		return "", err
	}

	mySigningKey := []byte(signingKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * time.Duration(tokenExpiration))),
		Issuer:    "weight-tracker",
		Subject:   userId,
		Audience:  []string{"weight-tracker"},
	})
	return token.SignedString(mySigningKey)
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

func (u *usersService) Login(ctx context.Context, arg loginRequest) (loginResponse, error) {
	signingKey := os.Getenv(utils.EnvJwtSignKey)
	tokenExpiration, err := strconv.Atoi(os.Getenv(utils.EnvJwtExpireMinutes))
	if err != nil {
		slog.Error("Failed to convert JWT_EXPIRATION to int", "error", err)
		return loginResponse{}, err
	}

	mySigningKey := []byte(signingKey)
	user, err := u.repo.GetByUsername(ctx, arg.Username)
	if err != nil {
		return loginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(arg.Password))
	if err != nil {
		return loginResponse{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * time.Duration(tokenExpiration))),
		Issuer:    "weight-tracker",
		Subject:   user.ID,
		Audience:  []string{"weight-tracker"},
	})

	signedToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return loginResponse{}, err
	}

	return loginResponse { Token: signedToken, UserId: user.ID }, nil
}

func NewService(repo UsersRepository) Service {
	return &usersService{repo}
}
