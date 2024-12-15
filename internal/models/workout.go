package models

import (
	"time"

	"github.com/google/uuid"
)

type Workout struct {
	id uuid.UUID
	date time.Time
	name string
	
	exercices []Exercise
}
