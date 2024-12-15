package models

import (
	"time"

	"github.com/google/uuid"
)

type ExerciseDto struct {
	id uuid.UUID
	name string
	date time.Time

	sets []Set
	exerciseTypeId uuid.UUID
}
