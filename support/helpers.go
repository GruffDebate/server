package support

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func StringPtr(s string) *string {
	return &s
}

func UintPtr(u uint) *uint {
	return &u
}

func IntPtr(i int) *int {
	return &i
}

func AUintToAInt(au []uint) []int {
	ai := make([]int, len(au))
	for i := 0; i < len(au); i++ {
		ai[i] = int(au[i])
	}
	return ai
}

func Atou(a string) (uint, error) {
	i, err := strconv.Atoi(a)
	if err != nil {
		return 0, err
	}
	return uint(i), nil
}

func Utoa(u uint) string {
	return fmt.Sprintf("%d", u)
}

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func CamelToSnake(s string) string {
	if s == "" {
		return ""
	}
	var result string
	var words []string
	var lastPos int
	rs := []rune(s)

	for i := 0; i < len(rs); i++ {
		if i > 0 && unicode.IsUpper(rs[i]) {
			words = append(words, s[lastPos:i])
			lastPos = i
		}
	}

	// append the last word
	if s[lastPos:] != "" {
		words = append(words, s[lastPos:])
	}

	for k, word := range words {
		if k > 0 {
			result += "_"
		}

		result += strings.ToLower(word)
	}

	return result
}

// SnakeToCamel returns a string converted from snake case to uppercase
func SnakeToCamel(s string) string {
	if s == "" {
		return ""
	}
	var result string

	words := strings.Split(s, "_")

	for _, word := range words {
		w := []rune(word)
		w[0] = unicode.ToUpper(w[0])
		result += string(w)
	}

	return result
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func NanoToMicro(n int64) int64 {
	// Conversion to float64 was leading to imperfect results!
	//return int64(Round(float64(n) / 1000.0))
	return (n + 500) / 1000
}

func MicroToNano(ms int64) int64 {
	return ms * 1000
}

func DurationInMs(start, end time.Time) int64 {
	return (end.UnixNano() - start.UnixNano()) / 1000000
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

func RoundToDecimal(f float64, decimals int) float64 {
	factor := math.Pow(10, float64(decimals))
	res := f * factor
	res = Round(res)
	res = res / factor
	return res
}

type Timestamp struct {
	Time time.Time
}

func (ts Timestamp) Timestamp() int64 {
	return NanoToMicro(ts.Time.UnixNano())
}

func (ts *Timestamp) SetTimestamp(val int64) {
	ts.Time = time.Unix(0, MicroToNano(val))
}

func NewTimestamp(val int64) Timestamp {
	t := time.Unix(0, MicroToNano(val))
	return Timestamp{Time: t}
}

type NullableTimestamp struct {
	Time time.Time
}

func (ts NullableTimestamp) Timestamp() int64 {
	return NanoToMicro(ts.Time.UnixNano())
}

func (ts *NullableTimestamp) SetTimestamp(val int64) {
	ts.Time = time.Unix(0, MicroToNano(val))
}

func NewNullableTimestamp(val int64) *NullableTimestamp {
	t := time.Unix(0, MicroToNano(val))
	return &NullableTimestamp{Time: t}
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

func InterfaceToFloat64(i interface{}) (float64, error) {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Int:
		return float64(i.(int)), nil
	case reflect.Uint:
		return float64(i.(uint)), nil
	case reflect.Uint64:
		return float64(i.(uint64)), nil
	case reflect.Int32:
		return float64(i.(int32)), nil
	case reflect.Int64:
		return float64(i.(int64)), nil
	case reflect.Float32:
		return float64(i.(float32)), nil
	case reflect.Float64:
		return i.(float64), nil
	default:
		return 0, fmt.Errorf("not implemented for type %v", v.Kind())
	}
}

func IsTypelessEqual(a, b interface{}) bool {
	aFloat, err := InterfaceToFloat64(a)
	if err != nil {
		return a == b
	}
	bFloat, err := InterfaceToFloat64(b)
	if err != nil {
		return false
	}

	return aFloat == bFloat
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
