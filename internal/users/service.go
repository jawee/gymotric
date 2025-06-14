package users

import (
	"context"
	"fmt"
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
	Token  string
	UserId string
}

type getMeResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     any    `json:"email"`
	CreatedOn string `json:"created_on"`
	UpdatedOn string `json:"updated_on"`
}

type Service interface {
	CreateAndReturnId(ctx context.Context, arg createUserAndReturnIdRequest) (string, error)
	Login(ctx context.Context, arg loginRequest) (loginResponse, error)
	Logout(ctx context.Context, accessToken string, refreshToken string) error
	CreateToken(userId string) (string, error)
	GetByUserId(ctx context.Context, userId string) (getMeResponse, error)
	ChangePassword(ctx context.Context, request changePasswordRequest, userId string) error
	CreateConfirmationToken(ctx context.Context, userId string, email string) (string, error)
	CreateResetPasswordToken(ctx context.Context, userId string) (string, error)
	ConfirmEmail(ctx context.Context, userId string, email string) error
	GetByEmail(ctx context.Context, email string) (getMeResponse, error)
	ResetPassword(ctx context.Context, userId string, newPassword string) error
	Register(ctx context.Context, arg registrationRequest) (string, error)
	CreateAccountConfirmationToken(ctx context.Context, userId string) (string, error)
	ConfirmAccount(context context.Context, userId string) error
	IsTokenValid(context context.Context, token, tokenType string) bool
}

func (s *usersService) IsTokenValid(context context.Context, cookieTokenStr, tokenType string) bool {
	panic("not implemented") // TODO: Implement this function
}

func (s *usersService) Logout(context context.Context, userToken string, refreshToken string) error {
	accessTokenExpiratation, _ := strconv.Atoi(os.Getenv(utils.EnvJwtExpireMinutes))
	refreshTokenExpiration, _ := strconv.Atoi(os.Getenv(utils.EnvJwtRefreshExpireMinutes))

	if userToken != "" {
		err := s.repo.InvalidateToken(context, repository.CreateExpiredTokenParams{
			Token:     userToken,
			TokenType: "access_token",
			CreatedOn: time.Now().UTC().Format(time.RFC3339),
			RemoveOn:  time.Now().UTC().Add(time.Minute * time.Duration(accessTokenExpiratation)).Format(time.RFC3339),
		})

		if err != nil {
			slog.Error("Failed to invalidate access token", "error", err)
		}
	}

	if refreshToken != "" {
		err := s.repo.InvalidateToken(context, repository.CreateExpiredTokenParams{
			Token:     refreshToken,
			TokenType: "refresh_token",
			CreatedOn: time.Now().UTC().Format(time.RFC3339),
			RemoveOn:  time.Now().UTC().Add(time.Minute * time.Duration(refreshTokenExpiration)).Format(time.RFC3339),
		})
		if err != nil {
			slog.Error("Failed to invalidate refresh token", "error", err)
		}
	}
	return nil
}

func (s *usersService) ConfirmAccount(context context.Context, userId string) error {
	user, err := s.repo.GetByUserId(context, userId)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUser(context, repository.UpdateUserParams{
		ID:         user.ID,
		Email:      user.Email,
		Password:   user.Password,
		IsVerified: true,
		UpdatedOn:  time.Now().UTC().Format(time.RFC3339),
	})

	if err != nil {
		return err
	}

	return nil
}

type usersService struct {
	repo UsersRepository
}

func (u *usersService) Register(ctx context.Context, arg registrationRequest) (string, error) {
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
		Email:     arg.Email,
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
	}
	id, err := u.repo.CreateAndReturnId(ctx, user)
	if err != nil {
		slog.Error("Failed to create user", "error", err)
		return "", err
	}

	return id, nil
}

func (u *usersService) ResetPassword(ctx context.Context, userId string, newPassword string) error {
	user, err := u.repo.GetByUserId(ctx, userId)
	if err != nil {
		return err
	}
	newPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = u.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:         user.ID,
		Email:      user.Email,
		Password:   string(newPasswordBytes),
		UpdatedOn:  time.Now().UTC().Format(time.RFC3339),
		IsVerified: user.IsVerified,
	})

	return err
}

func (s *usersService) GetByEmail(ctx context.Context, email string) (getMeResponse, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return getMeResponse{}, err
	}
	return getMeResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedOn: user.CreatedOn,
		UpdatedOn: user.UpdatedOn,
	}, nil
}

func (s *usersService) ConfirmEmail(ctx context.Context, userId string, email string) error {
	user, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		return err
	}

	emailExists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		return err
	}
	if emailExists {
		return fmt.Errorf("email already exists")
	}

	err = s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:         user.ID,
		Email:      email,
		Password:   user.Password,
		UpdatedOn:  time.Now().UTC().Format(time.RFC3339),
		IsVerified: user.IsVerified,
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *usersService) ChangePassword(ctx context.Context, request changePasswordRequest, userId string) error {
	user, err := s.repo.GetByUserId(ctx, userId)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
	if err != nil {
		return err
	}

	newPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	err = s.repo.UpdateUser(ctx, repository.UpdateUserParams{
		ID:         user.ID,
		Email:      user.Email,
		Password:   string(newPasswordBytes),
		UpdatedOn:  time.Now().UTC().Format(time.RFC3339),
		IsVerified: user.IsVerified,
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *usersService) GetByUserId(ctx context.Context, userId string) (getMeResponse, error) {
	user, err := u.repo.GetByUserId(ctx, userId)
	if err != nil {
		return getMeResponse{}, err
	}

	return getMeResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedOn: user.CreatedOn,
		UpdatedOn: user.UpdatedOn,
	}, nil
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

type emailConfirmationCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type confirmAccountCustomClaims struct {
	jwt.RegisteredClaims
}

func (e emailConfirmationCustomClaims) GetEmail() (string, error) {
	return e.Email, nil
}

func (u *usersService) CreateResetPasswordToken(ctx context.Context, userId string) (string, error) {
	signingKey := os.Getenv(utils.EnvJwtSignKey)

	mySigningKey := []byte(signingKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * time.Duration(utils.ResetPasswordTokenExpireMinutes))),
		Issuer:    "weight-tracker",
		Subject:   userId,
		Audience:  []string{"weight-tracker"},
	})

	return token.SignedString(mySigningKey)
}

func (u *usersService) CreateAccountConfirmationToken(ctx context.Context, userId string) (string, error) {
	signingKey := os.Getenv(utils.EnvJwtSignKey)

	mySigningKey := []byte(signingKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, confirmAccountCustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * time.Duration(utils.AccountConfirmationTokenExpireMinutes))),
			Issuer:    "weight-tracker",
			Subject:   userId,
			Audience:  []string{"weight-tracker"},
		}})

	return token.SignedString(mySigningKey)
}

func (u *usersService) CreateConfirmationToken(ctx context.Context, userId string, email string) (string, error) {
	emailExists, err := u.repo.EmailExists(ctx, email)
	if err != nil {
		return "", err
	}
	if emailExists {
		return "", fmt.Errorf("email already exists")
	}

	signingKey := os.Getenv(utils.EnvJwtSignKey)

	mySigningKey := []byte(signingKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, emailConfirmationCustomClaims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Minute * time.Duration(utils.EmailConfirmationTokenExpireMinutes))),
			Issuer:    "weight-tracker",
			Subject:   userId,
			Audience:  []string{"weight-tracker"},
		}})

	return token.SignedString(mySigningKey)
}

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

	if user.IsVerified == false {
		return loginResponse{}, fmt.Errorf("user is not verified")
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

	return loginResponse{Token: signedToken, UserId: user.ID}, nil
}

func NewService(repo UsersRepository) Service {
	return &usersService{repo}
}
