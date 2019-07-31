package gruff

import (
	"fmt"
	"reflect"
	"time"

	"github.com/GruffDebate/server/support"
	"github.com/google/uuid"
)

type VersionedModel struct {
	Key         string     `json:"_key"`
	ID          string     `json:"id"`
	CreatedAt   time.Time  `json:"start"`
	UpdatedAt   time.Time  `json:"mod"`
	DeletedAt   *time.Time `json:"end"`
	CreatedByID string     `json:"creator"`
	UpdatedByID string     `json:"editor,omitempty"`
	QueryAt     *time.Time `json:"-"`
}

func (vm *VersionedModel) PrepareForCreate(ctx *ServerContext) {
	u := ctx.UserContext
	vm.Key = uuid.New().String()
	if vm.ID == "" {
		// Brand new item, rather than a new version
		vm.ID = uuid.New().String()
	}
	vm.CreatedAt = ctx.RequestTime()
	vm.UpdatedAt = ctx.RequestTime()
	if vm.CreatedByID == "" {
		vm.CreatedByID = u.ArangoID()
	}
	vm.UpdatedByID = u.ArangoID()
	return
}

func (vm *VersionedModel) PrepareForDelete(ctx *ServerContext) {
	vm.DeletedAt = support.TimePtr(ctx.RequestTime())
	return
}

func (vm VersionedModel) QueryDate() *time.Time {
	var queryAt *time.Time
	if vm.QueryAt != nil {
		queryAt = vm.QueryAt
	} else if vm.DeletedAt != nil {
		beforeDelete := vm.DeletedAt.Add(-1 * time.Millisecond)
		queryAt = &beforeDelete
	}
	return queryAt
}

func (vm VersionedModel) DateFilter(bindVars map[string]interface{}) string {
	queryAt := vm.QueryDate()
	if queryAt != nil {
		bindVars["query_at"] = queryAt
		query := fmt.Sprintf("FILTER obj.start <= @query_at AND (obj.end == null OR obj.end > @query_at)")
		return query
	} else {
		return "FILTER obj.end == null"
	}
}

func IsVersionedModel(t reflect.Type) bool {
	_, is := t.FieldByName("VersionedModel")
	return is
}

func GetVersionedModel(item interface{}) (VersionedModel, Error) {
	if !IsVersionedModel(reflect.TypeOf(item)) {
		return VersionedModel{}, NewServerError("Item is not a VersionedModel")
	}

	id := reflect.ValueOf(item).FieldByName("VersionedModel").Interface().(VersionedModel)
	return id, nil
}
