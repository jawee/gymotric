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
	}, nil).Once()

	service := NewService(&repoMock)
	token, err := service.Login(context.Background(), loginRequest{
		Username: "testusername", Password: "test"})

	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	repoMock.AssertExpectations(t)
}
