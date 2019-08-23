package support

import (
	"reflect"
	"strings"
)

func NewInstance(item interface{}) interface{} {
	val := reflect.ValueOf(item)
	tp := reflect.TypeOf(reflect.Indirect(val))
	inst := reflect.New(tp).Interface()
	return inst
}

func JsonName(f reflect.StructField) string {
	tag := f.Tag
	jsonTag := tag.Get("json")
	if jsonTag == "" {
		return ""
	}

	vals := strings.Split(jsonTag, ",")
	return vals[0]
}

func IsZeroStruct(item interface{}) bool {
	zero := false

	if item != nil {
		v := reflect.ValueOf(item)
		t := v.Type()
		k := t.Kind()

		if k == reflect.Struct {
			z := reflect.Zero(t)
			if reflect.DeepEqual(v.Interface(), z.Interface()) {
				zero = true
			}
		}
	}

	return zero
}

func SetValue(v reflect.Value, destType reflect.Type, newVal interface{}) {
	vSet := v.Elem()

	if vSet.Kind() == reflect.Ptr {
		if newVal == nil {
			vSet.Set(reflect.Zero(vSet.Type()))
			return
		}

		floatVal, err := InterfaceToFloat64(newVal)
		if err == nil {
			switch destType.Kind() {
			case reflect.Int:
				n := int(floatVal)
				vSet.Set(reflect.ValueOf(&n))
			case reflect.Int32:
				n := int32(floatVal)
				vSet.Set(reflect.ValueOf(&n))
			case reflect.Int64:
				n := int64(floatVal)
				vSet.Set(reflect.ValueOf(&n))
			case reflect.Uint:
				n := uint(floatVal)
				vSet.Set(reflect.ValueOf(&n))
			case reflect.Uint8:
				n := uint8(floatVal)
				vSet.Set(reflect.ValueOf(&n))
			case reflect.Uint16:
				n := uint16(floatVal)
				vSet.Set(reflect.ValueOf(&n))
			case reflect.Uint32:
				n := uint32(floatVal)
				vSet.Set(reflect.ValueOf(&n))
			case reflect.Uint64:
				n := uint64(floatVal)
				vSet.Set(reflect.ValueOf(&n))
			case reflect.Float32:
				n := float32(floatVal)
				vSet.Set(reflect.ValueOf(&n))
			case reflect.Float64:
				vSet.Set(reflect.ValueOf(&newVal))
			case reflect.Struct:
				if destType == reflect.TypeOf(NullableTimestamp{}) {
					ts := NewNullableTimestamp(int64(floatVal))
					vSet.Set(reflect.ValueOf(ts))

				} else if destType == reflect.TypeOf(Timestamp{}) {
					ts := NewTimestamp(int64(floatVal))
					vSet.Set(reflect.ValueOf(ts))

				}
			default:
				vSet.Set(reflect.ValueOf(&floatVal))
			}
		} else {
			if destType.Kind() == reflect.String {
				strVal := newVal.(string)
				vSet.Set(reflect.ValueOf(&strVal))
			} else {
				vSet.Set(reflect.ValueOf(&newVal))
			}
		}
	} else {
		floatVal, err := InterfaceToFloat64(newVal)
		if err == nil {
			switch destType.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64:
				vSet.SetInt(int64(floatVal))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				vSet.SetUint(uint64(floatVal))
			case reflect.Float32, reflect.Float64:
				vSet.SetFloat(floatVal)
			default:
				vSet.Set(reflect.ValueOf(floatVal))
			}
		} else {
			vSet.Set(reflect.ValueOf(newVal))
		}
	}
}

func SetZeroValue(v reflect.Value) {
	v.Set(reflect.Zero(v.Type()))
}
