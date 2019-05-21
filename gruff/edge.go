package gruff

import (
	"reflect"
	"time"

	"github.com/GruffDebate/server/support"
	"github.com/google/uuid"
)

type Edge struct {
	Key         string     `json:"_key"`
	CreatedAt   time.Time  `json:"start"`
	DeletedAt   *time.Time `json:"end"`
	CreatedByID string     `json:"creator"`
}

func (e Edge) ValidateForCreate() GruffError {
	return ValidateStruct(e)
}

func (e Edge) ValidateField(f string) GruffError {
	return ValidateStructField(e, f)
}

func (e *Edge) PrepareForCreate() {
	e.Key = uuid.New().String()
	e.CreatedAt = time.Now()
	return
}

func (e *Edge) PrepareForDelete() {
	e.DeletedAt = support.TimePtr(time.Now())
	return
}

func IsEdge(t reflect.Type) bool {
	_, is := t.FieldByName("Edge")
	return is
}
