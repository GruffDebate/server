package gruff

import (
	"reflect"
)

type Scorer interface {
	Score(*ServerContext) (float32, Error)
	UpdateScore(*ServerContext) Error
}

func IsScorer(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Scorer)(nil)).Elem()
	return t.Implements(modelType)
}
