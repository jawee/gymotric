package sets

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

func (r *repoMock) CreateAndReturnId(ctx context.Context, arg repository.CreateSetAndReturnIdParams) (string, error) {
	args := r.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

func (r *repoMock) DeleteById(ctx context.Context, id string) (int64, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repoMock) GetAll(ctx context.Context) ([]Set, error) {
	args := r.Called(ctx)
	return args.Get(0).([]Set), args.Error(1)
}

func (r *repoMock) GetById(ctx context.Context, id string) (Set, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(Set), args.Error(1)
}

func (r *repoMock) GetByExerciseId(ctx context.Context, exerciseID string) ([]Set, error) {
	args := r.Called(ctx, exerciseID)
	return args.Get(0).([]Set), args.Error(1)
}

func TestGetAll(t *testing.T) {
	expected := []Set{
		{ID: "a", Repetitions: 1, Weight: 10.0, ExerciseID: "exerciseA"},
		{ID: "b", Repetitions: 1, Weight: 10.0, ExerciseID: "exerciseA"},
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetByExerciseId", ctx, "exerciseA").Return([]Set{
		{ID: "a", Repetitions: 1, Weight: 10.0, ExerciseID: "exerciseA"},
		{ID: "b", Repetitions: 1, Weight: 10.0, ExerciseID: "exerciseA"},
	}, nil).Once()

	service := NewService(&repoMock)

	result, err := service.GetByExerciseId(ctx, "exerciseA")

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}

func TestCreateAndReturnId(t *testing.T) {
	ctx := context.Background()

	setId, _ := uuid.NewV7()
	exerciseId, _ := uuid.NewV7()

	repoMock := repoMock{}
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func(input repository.CreateSetAndReturnIdParams) bool {
		return input.Weight == 10.5 && input.Repetitions == 1 && input.ExerciseID == exerciseId.String()
	})).Return(setId.String(), nil).Once()

	service := NewService(&repoMock)
	id, err := service.CreateAndReturnId(context.Background(), createSetRequest{
		Repetitions: 1,
		Weight: 10.5,

	}, exerciseId.String())

	assert.Nil(t, err)
	assert.Equal(t, setId.String(), id)
	repoMock.AssertExpectations(t)
}

func TestDeleteById(t *testing.T) {
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("DeleteById", ctx, "a").Return(int64(1), nil).Once()

	service := NewService(&repoMock)
	err := service.DeleteById(ctx, "a")

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}
