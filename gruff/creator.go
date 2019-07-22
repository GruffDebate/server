package gruff

import (
	"reflect"
)

type Creator interface {
	Create(*ServerContext) GruffError
}

func IsCreator(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Creator)(nil)).Elem()
	return t.Implements(modelType)
}
