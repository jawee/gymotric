package models

import "github.com/google/uuid"

type Set struct {
	id uuid.UUID
	weight float64
	repetitions int

	exerciseId uuid.UUID
}
