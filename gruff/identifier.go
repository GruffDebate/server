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

func (i *Identifier) PrepareForCreate(u User) {
	i.Key = uuid.New().String()
	if i.ID == "" {
		// Brand new item, rather than a new version
		i.ID = uuid.New().String()
	}
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
	i.CreatedByID = u.ArangoID()
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

func GetIdentifier(item interface{}) (Identifier, GruffError) {
	if !IsIdentifier(reflect.TypeOf(item)) {
		return Identifier{}, NewServerError("Item is not an Identifier")
	}

	id := reflect.ValueOf(item).FieldByName("Identifier").Interface().(Identifier)
	return id, nil
}
