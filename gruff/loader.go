package gruff

import (
	"reflect"
)

type Loader interface {
	Load(*ServerContext) Error
	LoadFull(*ServerContext) Error
}

func IsLoader(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Loader)(nil)).Elem()
	return t.Implements(modelType)
}
