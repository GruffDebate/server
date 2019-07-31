package gruff

import (
	"reflect"
)

// UserCanUpdate should be called with the list of values that will be updated,
type Restrictor interface {
	UserCanView(ctx *ServerContext) (bool, Error)
	UserCanCreate(ctx *ServerContext) (bool, Error)
	UserCanUpdate(ctx *ServerContext, updates map[string]interface{}) (bool, Error)
	UserCanDelete(ctx *ServerContext) (bool, Error)
}

func IsRestrictor(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Restrictor)(nil)).Elem()
	return t.Implements(modelType)
}
