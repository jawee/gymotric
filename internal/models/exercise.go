package models

import "github.com/google/uuid"

type Exercise struct {
	id uuid.UUID
	name string
	weight float64
}
