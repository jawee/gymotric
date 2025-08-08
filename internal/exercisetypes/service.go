package exercisetypes

import (
	"context"
	"fmt"
	"sort"
	"strings"
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
	GetLastWeightRepsByExerciseTypeId(context context.Context, exerciseTypeId string, userId string) (MaxLastWeightReps, error)
	GetMaxWeightRepsByExerciseTypeId(context context.Context, exerciseTypeId string, userId string) (MaxLastWeightReps, error)
	UpdateById(context context.Context, exerciseTypeId string, updateExerciseTypeRequest updateExerciseTypeRequest, userId string) error
}

func (s *exerciseTypeService) UpdateById(context context.Context, exerciseTypeId string, updateExerciseTypeRequest updateExerciseTypeRequest, userId string) error {
	arg := repository.UpdateExerciseTypeParams{
		ID: exerciseTypeId,
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
		Name: updateExerciseTypeRequest.Name,
		UserID: userId,
	}

	err := s.repo.UpdateById(context, arg)

	return fmt.Errorf("failed to update exercise type: %w", err) 
}

func (s *exerciseTypeService) GetLastWeightRepsByExerciseTypeId(context context.Context, exerciseTypeId string, userId string) (MaxLastWeightReps, error) {
	arg := repository.GetLastWeightRepsByExerciseTypeIdParams{
		ID:     exerciseTypeId,
		UserID: userId,
	}

	return s.repo.GetLastWeightRepsByExerciseTypeId(context, arg)
}

func (s *exerciseTypeService) GetMaxWeightRepsByExerciseTypeId(context context.Context, exerciseTypeId string, userId string) (MaxLastWeightReps, error) {
	arg := repository.GetMaxWeightRepsByExerciseTypeIdParams{
		ID:     exerciseTypeId,
		UserID: userId,
	}

	return s.repo.GetMaxWeightRepsByExerciseTypeId(context, arg)
}

func (s *exerciseTypeService) CreateAndReturnId(context context.Context, exerciseType createExerciseTypeRequest, userId string) (string, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("failed to generate UUID: %w", err)
	}

	toCreate := repository.CreateExerciseTypeAndReturnIdParams{
		ID:        uuid.String(),
		Name:      strings.TrimSpace(exerciseType.Name),
		CreatedOn: time.Now().UTC().Format(time.RFC3339),
		UpdatedOn: time.Now().UTC().Format(time.RFC3339),
		UserID:    userId,
	}

	id, err := s.repo.CreateAndReturnId(context, toCreate)
	if err != nil {
		return "", fmt.Errorf("failed to create exercise type: %w", err)
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
		return []ExerciseType{}, fmt.Errorf("failed to get all exercise types: %w", err)
	}

	sort.Slice(exerciseTypes, func(i, j int) bool {
		return exerciseTypes[i].Name < exerciseTypes[j].Name
	})

	return exerciseTypes, nil
}

type exerciseTypeService struct {
	repo ExerciseTypeRepository
}
