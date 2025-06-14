package users

import (
	"context"
	"os"
	"testing"
	"weight-tracker/internal/repository"
	"weight-tracker/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type repoMock struct {
	mock.Mock
}

func (u *repoMock) CheckIfTokenExists(ctx context.Context, arg repository.CheckIfTokenExistsParams) (int64, error) {
	args := u.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (u *repoMock) GetByEmail(ctx context.Context, email string) (User, error) {
	args := u.Called(ctx, email)
	return args.Get(0).(User), args.Error(1)
}

func (u *repoMock) EmailExists(ctx context.Context, email string) (bool, error) {
	args := u.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (u *repoMock) UpdateUser(ctx context.Context, arg repository.UpdateUserParams) error {
	args := u.Called(ctx, arg)
	return args.Error(0)
}

func (r *repoMock) CreateAndReturnId(ctx context.Context, arg repository.CreateUserAndReturnIdParams) (string, error) {
	args := r.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

func (r *repoMock) GetByUsername(ctx context.Context, arg string) (User, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).(User), args.Error(1)
}

func (r *repoMock) GetByUserId(ctx context.Context, arg string) (User, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).(User), args.Error(1)
}

func (r *repoMock) InvalidateToken(ctx context.Context, arg repository.CreateExpiredTokenParams) error {
	args := r.Called(ctx, arg)
	return args.Error(0)
}

func TestCreateAndReturnId(t *testing.T) {
	ctx := context.Background()

	userId, _ := uuid.NewV7()
	repoMock := repoMock{}
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateUserAndReturnIdParams) bool {
		return input.Username == "testusername" && input.ID != "" && input.CreatedOn != "" && input.UpdatedOn != "" && input.Password != ""
	})).Return(userId.String(), nil).Once()

	service := NewService(&repoMock)
	id, err := service.CreateAndReturnId(context.Background(), createUserAndReturnIdRequest{
		Username: "testusername", Password: "test"})

	assert.Nil(t, err)
	assert.Equal(t, userId.String(), id)
	repoMock.AssertExpectations(t)
}

func TestLoginAndReturnToken(t *testing.T) {
	os.Setenv(utils.EnvJwtExpireMinutes, "10")
	ctx := context.Background()

	userId, _ := uuid.NewV7()
	pwBytes, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	repoMock := repoMock{}
	repoMock.On("GetByUsername", ctx, "testusername").Return(User{
		ID:        userId.String(),
		Username:  "testusername",
		Password:  string(pwBytes),
		CreatedOn: "2024-09-05T19:22:00Z",
		UpdatedOn: "2024-09-05T19:22:00Z",
		Email:     "test@test.se",
		IsVerified: true,
	}, nil).Once()

	service := NewService(&repoMock)
	loginResponse, err := service.Login(context.Background(), loginRequest{
		Username: "testusername", Password: "test"})

	assert.Nil(t, err)
	assert.NotEmpty(t, loginResponse)
	assert.NotEmpty(t, loginResponse.Token)
	assert.Equal(t, userId.String(), loginResponse.UserId)
	repoMock.AssertExpectations(t)
}

func TestLoginUserNotFoundReturnsEmptyAndErr(t *testing.T) {
	ctx := context.Background()
	repoMock := repoMock{}
	repoMock.On("GetByUsername", ctx, "testusername").Return(User{}, testError).Once()

	service := NewService(&repoMock)
	token, err := service.Login(context.Background(), loginRequest{
		Username: "testusername", Password: "test"})

	assert.NotNil(t, err)
	assert.Empty(t, token)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestLoginPasswordNotMatchReturnsEmptyAndErr(t *testing.T) {
	os.Setenv(utils.EnvJwtExpireMinutes, "10")
	ctx := context.Background()

	userId, _ := uuid.NewV7()
	pwBytes, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	repoMock := repoMock{}
	repoMock.On("GetByUsername", ctx, "testusername").Return(User{
		ID:        userId.String(),
		Username:  "testusername",
		Password:  string(pwBytes),
		CreatedOn: "2024-09-05T19:22:00Z",
		UpdatedOn: "2024-09-05T19:22:00Z",
		Email:     "test@test.se",
		IsVerified: true,
	}, nil).Once()

	service := NewService(&repoMock)
	token, err := service.Login(context.Background(), loginRequest{
		Username: "testusername", Password: "wrong"})

	assert.NotNil(t, err)
	assert.Empty(t, token)
	assert.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
	repoMock.AssertExpectations(t)
}

func TestCreateToken(t *testing.T) {
	os.Setenv(utils.EnvJwtExpireMinutes, "10")
	os.Setenv(utils.EnvJwtSignKey, "testkey")

	userId, _ := uuid.NewV7()
	repoMock := repoMock{}

	service := NewService(&repoMock)
	token, err := service.CreateToken(userId.String())

	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	repoMock.AssertExpectations(t)
}

func TestGetUserById(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{
		ID:        userId.String(),
		Username:  "testusername",
		Password:  "test",
		CreatedOn: "2024-09-05T19:22:00Z",
		UpdatedOn: "2024-09-05T19:22:00Z",
		Email:     "test@test.se",
	}, nil).Once()
	service := NewService(&repoMock)
	user, err := service.GetByUserId(ctx, userId.String())
	assert.Nil(t, err)
	assert.Equal(t, userId.String(), user.ID)
	assert.Equal(t, "testusername", user.Username)
	assert.Equal(t, "test@test.se", user.Email)
	assert.Equal(t, "2024-09-05T19:22:00Z", user.CreatedOn)
	assert.Equal(t, "2024-09-05T19:22:00Z", user.UpdatedOn)
	repoMock.AssertExpectations(t)
}

func TestGetUserByIdErr(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{}, testError).Once()
	service := NewService(&repoMock)
	user, err := service.GetByUserId(ctx, userId.String())
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	assert.Empty(t, user)
	repoMock.AssertExpectations(t)
}

func TestChangePassword(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	pwBytes, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{
		ID:         userId.String(),
		Username:   "testusername",
		Password:   string(pwBytes),
		CreatedOn:  "2024-09-05T19:22:00Z",
		UpdatedOn:  "2024-09-05T19:22:00Z",
		IsVerified: true,
	}, nil).Once()

	repoMock.On("UpdateUser", ctx, mock.MatchedBy(func(input repository.UpdateUserParams) bool {
		return input.ID == userId.String() && input.Password != "" && input.Email == nil && input.IsVerified == true && input.UpdatedOn != "2024-09-05T19:22:00Z"
	})).Return(nil).Once()

	service := NewService(&repoMock)
	err := service.ChangePassword(ctx, changePasswordRequest{
		NewPassword: "newpassword",
		OldPassword: "test"}, userId.String())

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestChangePasswordUserNotFoundErr(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{}, testError).Once()
	service := NewService(&repoMock)
	err := service.ChangePassword(ctx, changePasswordRequest{
		NewPassword: "newpassword",
		OldPassword: "test"}, userId.String())
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestChangePasswordWrongPasswordErr(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	pwBytes, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{
		ID:        userId.String(),
		Username:  "testusername",
		Password:  string(pwBytes),
		CreatedOn: "2024-09-05T19:22:00Z",
		UpdatedOn: "2024-09-05T19:22:00Z",
	}, nil).Once()

	service := NewService(&repoMock)
	err := service.ChangePassword(ctx, changePasswordRequest{
		NewPassword: "newpassword",
		OldPassword: "wrongpassword"}, userId.String())
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, bcrypt.ErrMismatchedHashAndPassword)
	repoMock.AssertExpectations(t)
}

func TestChangePasswordUpdateUserFailsErr(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	pwBytes, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{
		ID:        userId.String(),
		Username:  "testusername",
		Password:  string(pwBytes),
		CreatedOn: "2024-09-05T19:22:00Z",
		UpdatedOn: "2024-09-05T19:22:00Z",
	}, nil).Once()

	repoMock.On("UpdateUser", ctx, mock.Anything).Return(testError).Once()

	service := NewService(&repoMock)
	err := service.ChangePassword(ctx, changePasswordRequest{
		NewPassword: "newpassword",
		OldPassword: "test",
	}, userId.String())
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestCreateConfirmationToken(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	repoMock := repoMock{}
	repoMock.On("EmailExists", ctx, "test@test.se").Return(false, nil).Once()

	service := NewService(&repoMock)
	token, err := service.CreateConfirmationToken(ctx, userId.String(), "test@test.se")

	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	repoMock.AssertExpectations(t)
}

func TestCreateConfirmationTokenEmailAlreadyExistsErr(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	repoMock := repoMock{}
	repoMock.On("EmailExists", ctx, "test@test.se").Return(true, nil).Once()

	service := NewService(&repoMock)
	token, err := service.CreateConfirmationToken(ctx, userId.String(), "test@test.se")

	assert.NotNil(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "email already exists", err.Error())
	repoMock.AssertExpectations(t)
}

func TestCreateConfirmationTokenRepoErr(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	repoMock := repoMock{}
	repoMock.On("EmailExists", ctx, "test@test.se").Return(false, testError).Once()

	service := NewService(&repoMock)
	token, err := service.CreateConfirmationToken(ctx, userId.String(), "test@test.se")

	assert.NotNil(t, err)
	assert.Empty(t, token)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestCreateResetPasswordToken(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()

	service := NewService(nil)
	token, err := service.CreateResetPasswordToken(ctx, userId.String())

	assert.Nil(t, err)
	assert.NotEmpty(t, token)
}

func TestConfirmEmail(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()

	email := "test@test.se"

	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{
		ID:         userId.String(),
		Username:   "testusername",
		Password:   "test",
		CreatedOn:  "2024-09-05T19:22:00Z",
		UpdatedOn:  "2024-09-05T19:22:00Z",
		Email:      email,
		IsVerified: true,
	}, nil).Once()
	repoMock.On("EmailExists", ctx, email).Return(false, nil).Once()
	repoMock.On("UpdateUser", ctx, mock.MatchedBy(func(input repository.UpdateUserParams) bool {
		return input.ID == userId.String() && input.Email == email && input.IsVerified == true && input.UpdatedOn != "2024-09-05T19:22:00Z"
	})).Return(nil).Once()
	service := NewService(&repoMock)

	err := service.ConfirmEmail(ctx, userId.String(), email)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestConfirmEmailUserNotFound(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()

	email := "test@test.se"

	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{}, testError).Once()

	service := NewService(&repoMock)

	err := service.ConfirmEmail(ctx, userId.String(), email)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestConfirmEmailEmailAlreadyInUse(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()

	email := "test@test.se"

	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{
		ID:        userId.String(),
		Username:  "testusername",
		Password:  "test",
		CreatedOn: "2024-09-05T19:22:00Z",
		UpdatedOn: "2024-09-05T19:22:00Z",
		Email:     email,
	}, nil).Once()
	repoMock.On("EmailExists", ctx, email).Return(true, nil).Once()

	service := NewService(&repoMock)

	err := service.ConfirmEmail(ctx, userId.String(), email)

	assert.NotNil(t, err)
	repoMock.AssertExpectations(t)
}

func TestConfirmEmailEmailExistsErr(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()

	email := "test@test.se"

	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{
		ID:        userId.String(),
		Username:  "testusername",
		Password:  "test",
		CreatedOn: "2024-09-05T19:22:00Z",
		UpdatedOn: "2024-09-05T19:22:00Z",
		Email:     email,
	}, nil).Once()
	repoMock.On("EmailExists", ctx, email).Return(false, testError).Once()
	service := NewService(&repoMock)

	err := service.ConfirmEmail(ctx, userId.String(), email)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestConfirmEmailUpdateErr(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()

	email := "test@test.se"

	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{
		ID:         userId.String(),
		Username:   "testusername",
		Password:   "test",
		CreatedOn:  "2024-09-05T19:22:00Z",
		UpdatedOn:  "2024-09-05T19:22:00Z",
		Email:      email,
		IsVerified: true,
	}, nil).Once()
	repoMock.On("EmailExists", ctx, email).Return(false, nil).Once()
	repoMock.On("UpdateUser", ctx, mock.MatchedBy(func(input repository.UpdateUserParams) bool {
		return input.ID == userId.String() && input.Email == email && input.IsVerified == true && input.UpdatedOn != "2024-09-05T19:22:00Z"
	})).Return(testError).Once()
	service := NewService(&repoMock)

	err := service.ConfirmEmail(ctx, userId.String(), email)

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestResetPassword(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()

	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{
		ID:         userId.String(),
		Username:   "testusername",
		CreatedOn:  "2024-09-05T19:22:00Z",
		UpdatedOn:  "2024-09-05T19:22:00Z",
		IsVerified: true,
	}, nil).Once()
	repoMock.On("UpdateUser", ctx, mock.MatchedBy(func(input repository.UpdateUserParams) bool {
		return input.ID == userId.String() && input.Password != "" && input.Email == nil && input.IsVerified == true && input.UpdatedOn != "2024-09-05T19:22:00Z"
	})).Return(nil).Once()

	service := NewService(&repoMock)
	err := service.ResetPassword(ctx, userId.String(), "newpassword")

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestResetPasswordUserNotFound(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()

	repoMock := repoMock{}
	repoMock.On("GetByUserId", ctx, userId.String()).Return(User{}, testError).Once()

	service := NewService(&repoMock)
	err := service.ResetPassword(ctx, userId.String(), "newpassword")

	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	repoMock.AssertExpectations(t)
}

func TestGetByEmail(t *testing.T) {
	ctx := context.Background()
	userId, _ := uuid.NewV7()
	email := "testemail@example.com"
	repoMock := repoMock{}
	repoMock.On("GetByEmail", ctx, email).Return(User{
		ID:        userId.String(),
		Username:  "testusername",
		CreatedOn: "2024-09-05T19:22:00Z",
		UpdatedOn: "2024-09-05T19:22:00Z",
		Email:     email,
	}, nil).Once()

	service := NewService(&repoMock)
	user, err := service.GetByEmail(ctx, email)
	assert.Nil(t, err)
	assert.Equal(t, userId.String(), user.ID)
	assert.Equal(t, "testusername", user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, "2024-09-05T19:22:00Z", user.CreatedOn)
	assert.Equal(t, "2024-09-05T19:22:00Z", user.UpdatedOn)
	repoMock.AssertExpectations(t)
}

func TestGetByEmailNotFoundReturnsEmptyWithErr(t *testing.T) {
	ctx := context.Background()
	email := "testemail@example.com"
	repoMock := repoMock{}
	repoMock.On("GetByEmail", ctx, email).Return(User{}, testError).Once()

	service := NewService(&repoMock)
	user, err := service.GetByEmail(ctx, email)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, testError)
	assert.Empty(t, user)
	repoMock.AssertExpectations(t)
}
