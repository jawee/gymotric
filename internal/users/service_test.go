package users

import (
	"context"
	"testing"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	mock.Mock
}

// CreateAndReturnId implements UsersRepository.
func (r *repoMock) CreateAndReturnId(ctx context.Context, arg repository.CreateUserAndReturnIdParams) (string, error) {
	args := r.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

// GetByUsernameAndPassword implements UsersRepository.
func (r *repoMock) GetByUsernameAndPassword(ctx context.Context, arg repository.GetUserByUsernameAndPasswordParams) (User, error) {
	panic("unimplemented")
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
