package gruff

import (
	"reflect"
)

type Updater interface {
	Update(*ServerContext, map[string]interface{}) GruffError
}

func IsUpdater(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Updater)(nil)).Elem()
	return t.Implements(modelType)
}
