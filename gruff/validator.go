package gruff

import (
	"fmt"
	"reflect"

	"github.com/asaskevich/govalidator"
)

type Validator interface {
	ValidateForCreate() Error
	ValidateForUpdate(map[string]interface{}) Error
	ValidateForDelete() Error
	ValidateField(string) Error
}

func IsValidator(t reflect.Type) bool {
	modelType := reflect.TypeOf((*Validator)(nil)).Elem()
	return t.Implements(modelType)
}

func ValidateStruct(item interface{}) Error {
	_, err := govalidator.ValidateStruct(item)
	if err != nil {
		return NewBusinessError(err.Error())
	}
	return nil
}

func ValidateStructField(item interface{}, f string) Error {
	_, err := govalidator.ValidateStruct(item)
	errStr := govalidator.ErrorByField(err, f)

	if errStr != "" {
		return NewBusinessError(fmt.Sprintf("%s: %s;", f, errStr))
	}

	return nil
}

func ValidateStructFields(item interface{}, fs []string) Error {
	_, err := govalidator.ValidateStruct(item)
	if err == nil {
		return nil
	}

	result := ""
	for _, f := range fs {
		errStr := govalidator.ErrorByField(err, f)
		if errStr != "" {
			result = fmt.Sprintf("%s%s: %s;", result, f, errStr)
		}
	}

	if result == "" {
		return nil
	}

	return NewBusinessError(result)
}

func ValidateRequiredFields(item interface{}, fields []string) Error {
	itemVal := reflect.ValueOf(item)
	errStr := ""

	for _, fName := range fields {
		f := itemVal.FieldByName(fName)
		if IsEmptyValue(f) {
			errStr = fmt.Sprintf("%s%s: non zero value required;", errStr, fName)
		}
	}

	if errStr != "" {
		return NewBusinessError(errStr)
	}

	return nil
}

// TODO: I think there's a reflect.Zero(Type) method that should simplify this. First, make sure tests are in place.
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
