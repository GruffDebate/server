package gruff

import (
	"reflect"
	"time"

	"github.com/GruffDebate/server/support"
	"github.com/google/uuid"
)

type Identifier struct {
	Key         string     `json:"_key"`
	ID          string     `json:"id"`
	CreatedAt   time.Time  `json:"start"`
	UpdatedAt   time.Time  `json:"mod"`
	DeletedAt   *time.Time `json:"end"`
	CreatedByID string     `json:"creator"`
}

func (i Identifier) ValidateForCreate() GruffError {
	return ValidateStruct(i)
}

func (i Identifier) ValidateForUpdate() GruffError {
	return i.ValidateForCreate()
}

func (i Identifier) ValidateField(f string) GruffError {
	return ValidateStructField(i, f)
}

func (i *Identifier) PrepareForCreate() {
	i.Key = uuid.New().String()
	if i.ID == "" {
		// Brand new item, rather than a new version
		i.ID = uuid.New().String()
	}
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
	return
}

func (i *Identifier) PrepareForDelete() {
	i.DeletedAt = support.TimePtr(time.Now())
	return
}

func IsIdentifier(t reflect.Type) bool {
	_, is := t.FieldByName("Identifier")
	return is
}
