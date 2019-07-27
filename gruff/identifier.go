package gruff

import (
	"fmt"
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
	QueryAt     *time.Time `json:"-"`
	CreatedByID string     `json:"creator"`
	UpdatedByID string     `json:"editor,omitempty"`
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

func (i *Identifier) PrepareForCreate(ctx *ServerContext) {
	u := ctx.UserContext
	i.Key = uuid.New().String()
	if i.ID == "" {
		// Brand new item, rather than a new version
		i.ID = uuid.New().String()
	}
	i.CreatedAt = ctx.RequestTime()
	i.UpdatedAt = ctx.RequestTime()
	if i.CreatedByID == "" {
		i.CreatedByID = u.ArangoID()
	}
	i.UpdatedByID = u.ArangoID()
	return
}

func (i *Identifier) PrepareForDelete(ctx *ServerContext) {
	i.DeletedAt = support.TimePtr(ctx.RequestTime())
	return
}

func (i Identifier) QueryDate() *time.Time {
	var queryAt *time.Time
	if i.QueryAt != nil {
		queryAt = i.QueryAt
	} else if i.DeletedAt != nil {
		beforeDelete := i.DeletedAt.Add(-1 * time.Millisecond)
		queryAt = &beforeDelete
	}
	return queryAt
}

func (i Identifier) DateFilter(bindVars map[string]interface{}) string {
	queryAt := i.QueryDate()
	if queryAt != nil {
		bindVars["query_at"] = queryAt
		query := fmt.Sprintf("FILTER obj.start <= @query_at AND (obj.end == null OR obj.end > @query_at)")
		return query
	} else {
		return "FILTER obj.end == null"
	}
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
