package gruff

import (
	"reflect"
)

// UserCanUpdate should be called with the list of values that will be updated,
type Restrictor interface {
	UserCanView(ctx *ServerContext) (bool, GruffError)
	UserCanCreate(ctx *ServerContext) (bool, GruffError)
	UserCanUpdate(ctx *ServerContext, updates map[string]interface{}) (bool, GruffError)
	UserCanDelete(ctx *ServerContext) (bool, GruffError)
}

func IsRestrictor(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Restrictor)(nil)).Elem()
	return t.Implements(modelType)
}
