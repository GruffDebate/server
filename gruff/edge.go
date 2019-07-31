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
	From        string     `json:"_from,omitempty"`
	To          string     `json:"_to,omitempty"`
}

func (e Edge) ValidateForCreate() Error {
	return ValidateStruct(e)
}

func (e Edge) ValidateField(f string) Error {
	return ValidateStructField(e, f)
}

func (e *Edge) PrepareForCreate(ctx *ServerContext) {
	u := ctx.UserContext
	e.Key = uuid.New().String()
	e.CreatedAt = ctx.RequestTime()
	e.CreatedByID = u.ArangoID()
	return
}

func (e *Edge) PrepareForDelete(ctx *ServerContext) {
	e.DeletedAt = support.TimePtr(ctx.RequestTime())
	return
}

func IsEdge(t reflect.Type) bool {
	_, is := t.FieldByName("Edge")
	return is
}
