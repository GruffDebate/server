package support

import (
	"reflect"
)

func NewInstance(item interface{}) interface{} {
	val := reflect.ValueOf(item)
	tp := reflect.TypeOf(reflect.Indirect(val))
	inst := reflect.New(tp).Interface()
	return inst
}
