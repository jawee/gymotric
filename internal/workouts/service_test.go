package workouts

import (
	"context"
	"testing"
	"time"
	"weight-tracker/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	mock.Mock
}

func (r *repoMock) CompleteById(ctx context.Context, arg repository.CompleteWorkoutByIdParams) (int64, error) {
	args := r.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (r *repoMock) CreateAndReturnId(ctx context.Context, arg repository.CreateWorkoutAndReturnIdParams) (string, error) {
	args := r.Called(ctx, arg)
	return args.String(0), args.Error(1)
}

func (r *repoMock) GetAll(ctx context.Context) ([]Workout, error) {
	args := r.Called(ctx)
	return args.Get(0).([]Workout), args.Error(1)
}

func (r *repoMock) GetById(ctx context.Context, id string) (Workout, error) {
	args := r.Called(ctx, id)
	return args.Get(0).(Workout), args.Error(1)
}

func TestGetAll(t *testing.T) {
	expected := []Workout{
		{ID: "a", Name: "A", CreatedOn: time.Now().UTC().Format(time.RFC3339), CompletedOn: time.Now().UTC().Format(time.RFC3339), UpdatedOn: time.Now().UTC().Format(time.RFC3339)},
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetAll", ctx).Return([]Workout{
		{ID: "a", Name: "A", CreatedOn: time.Now().UTC().Format(time.RFC3339), CompletedOn: time.Now().UTC().Format(time.RFC3339), UpdatedOn: time.Now().UTC().Format(time.RFC3339)},
	}, nil).Once()

	service := NewService(&repoMock)

	result, err := service.GetAll(ctx)

	assert.Nil(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}

func TestGetById(t *testing.T) {
	expected := Workout{
		ID:          "a",
		Name:        "A",
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn:   time.Now().UTC().Format(time.RFC3339),
	}

	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("GetById", ctx, "a").Return(Workout{
		ID:          "a",
		Name:        "A",
		CreatedOn:   time.Now().UTC().Format(time.RFC3339),
		CompletedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn:   time.Now().UTC().Format(time.RFC3339),
	}, nil).Once()

	service := NewService(&repoMock)

	result, err := service.GetById(ctx, "a")
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	repoMock.AssertExpectations(t)
}

func TestCreateAndReturnId(t *testing.T) {
	ctx := context.Background()

	request := createWorkoutRequest{
		Name: "A",
	}

	repoMock := repoMock{}
	repoMock.On("CreateAndReturnId", ctx, mock.MatchedBy(func (input repository.CreateWorkoutAndReturnIdParams) bool {
		return input.Name == "A" && input.ID != "" && input.CreatedOn != "" && input.UpdatedOn != ""
	})).Return("a", nil).Once()

	service := NewService(&repoMock)

	result, err := service.CreateAndReturnId(ctx, request)

	assert.Nil(t, err)
	assert.Equal(t, "a", result)
	repoMock.AssertExpectations(t)
}

func TestCompleteById(t *testing.T) {
	ctx := context.Background()

	repoMock := repoMock{}
	repoMock.On("CompleteById", ctx, mock.MatchedBy(func (input repository.CompleteWorkoutByIdParams) bool {
		return input.ID == "a"
	})).Return(int64(1), nil).Once()

	service := NewService(&repoMock)

	err := service.CompleteById(ctx, "a")

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}
