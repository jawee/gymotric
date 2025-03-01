package exercisetypes

import (
	"context"
	"sort"
	"time"
	"weight-tracker/internal/repository"

	"github.com/google/uuid"
)

func NewService(repo ExerciseTypeRepository) Service {
	return &exerciseTypeService{repo}
}

type Service interface {
	GetAll(context context.Context, userId string) ([]ExerciseType, error)
	DeleteById(context context.Context, exerciseTypeId string, userId string) error
	CreateAndReturnId(context context.Context, exerciseType createExerciseTypeRequest, userId string) (string, error)
}

func (s *exerciseTypeService) CreateAndReturnId(context context.Context, exerciseType createExerciseTypeRequest, userId string) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	toCreate := repository.CreateExerciseTypeAndReturnIdParams{
		ID:   uuid.String(),
		Name: exerciseType.Name,
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
		UserID: userId,
	}

	id, err := s.repo.CreateAndReturnId(context, toCreate)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *exerciseTypeService) DeleteById(context context.Context, exerciseTypeId string, userId string) error {
	arg := repository.DeleteExerciseTypeByIdParams{
		ID:     exerciseTypeId,
		UserID: userId,
	}
	return s.repo.DeleteById(context, arg)
}

func (s *exerciseTypeService) GetAll(context context.Context, userId string) ([]ExerciseType, error) {
	exerciseTypes, err := s.repo.GetAll(context, userId)
	if err != nil {
		return []ExerciseType{}, err
	}

	sort.Slice(exerciseTypes, func(i, j int) bool {
		return exerciseTypes[i].Name < exerciseTypes[j].Name
	})

	return exerciseTypes, nil
}

type exerciseTypeService struct {
	repo ExerciseTypeRepository
}
