package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkoutDto struct {
	id uuid.UUID
	date time.Time
	name string
	
	exercices []ExerciseDto
}
