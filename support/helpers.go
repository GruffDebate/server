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
