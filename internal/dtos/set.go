package models

import "github.com/google/uuid"

type SetDto struct {
	id uuid.UUID
	weight float64
	repetitions int

	exerciseId uuid.UUID
}
