package gruff

import (
	"reflect"
)

type Deleter interface {
	Delete(*ServerContext) GruffError
}

func IsDeleter(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Deleter)(nil)).Elem()
	return t.Implements(modelType)
}
