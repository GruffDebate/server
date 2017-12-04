package gruff

import (
	"github.com/google/uuid"
)

type ArgumentOpinion struct {
	Model
	UserID     uint64    `json:"userId"`
	User       *User     `json:"user,omitempty"`
	ArgumentID uuid.UUID `json:"argumentId" sql:"type:uuid"`
	Argument   *Argument `json:"argument,omitempty"`
	Strength   float64   `json:"strength"`
}

func (ao ArgumentOpinion) ValidateForCreate() GruffError {
	return ValidateStruct(ao)
}

func (ao ArgumentOpinion) ValidateForUpdate() GruffError {
	return ao.ValidateForCreate()
}

func (ao ArgumentOpinion) ValidateField(f string) GruffError {
	return ValidateStructField(ao, f)
}
