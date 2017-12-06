package support

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type TestStruct struct {
	ID          uint
	Title       string
	Description *string
}

func AssertEqual(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("assertEqual: Expected:\n%v\nbut got:\n%v\n (caller: %s)", expected, actual, CallerInfo())
	}
}

func AssertTrue(t *testing.T, actual bool) {
	if !actual {
		t.Errorf("assertTrue: Assertion is false (caller: %s)", CallerInfo())
	}
}

func AssertFalse(t *testing.T, actual bool) {
	if actual {
		t.Errorf("assertFalse: Assertion is true (caller: %s)", CallerInfo())
	}
}

func CallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func TestAtou(t *testing.T) {
	var nilErr error

	val, err := Atou("1")
	AssertEqual(t, nilErr, err)
	AssertEqual(t, uint(1), val)

	val, err = Atou("0")
	AssertEqual(t, nilErr, err)
	AssertEqual(t, uint(0), val)

	val, err = Atou("1384234")
	AssertEqual(t, nilErr, err)
	AssertEqual(t, uint(1384234), val)

	val, err = Atou("-1")
	AssertEqual(t, nilErr, err)
	AssertEqual(t, uint(18446744073709551615), val)

	val, err = Atou("not a number")
	AssertEqual(t, "strconv.Atoi: parsing \"not a number\": invalid syntax", err.Error())
	AssertEqual(t, uint(0), val)

	val, err = Atou("#234")
	AssertEqual(t, "strconv.Atoi: parsing \"#234\": invalid syntax", err.Error())
	AssertEqual(t, uint(0), val)
}

func TestUtoa(t *testing.T) {
	AssertEqual(t, "1", Utoa(1))
	AssertEqual(t, "0", Utoa(0))
	AssertEqual(t, "1384234", Utoa(1384234))
}

func TestCaptialize(t *testing.T) {
	AssertEqual(t, "Testing", Capitalize("testing"))
	AssertEqual(t, "This is a longer string", Capitalize("this is a longer string"))
	AssertEqual(t, "Already Capitalized", Capitalize("Already Capitalized"))
	AssertEqual(t, "", Capitalize(""))
}

func TestSnakeToCamel(t *testing.T) {
	AssertEqual(t, "Testing", SnakeToCamel("testing"))
	AssertEqual(t, "This is a longer string", SnakeToCamel("this is a longer string"))
	AssertEqual(t, "Already Capitalized", SnakeToCamel("Already Capitalized"))
	AssertEqual(t, "ActualSnake", SnakeToCamel("actual_snake"))
	AssertEqual(t, "ThisIsAMultipleSnake", SnakeToCamel("this_is_a_multiple_snake"))
	AssertEqual(t, "", SnakeToCamel(""))
}

func TestCamelToSnake(t *testing.T) {
	AssertEqual(t, "testing", CamelToSnake("Testing"))
	AssertEqual(t, "testing", CamelToSnake("testing"))
	AssertEqual(t, "actual_camel", CamelToSnake("ActualCamel"))
	AssertEqual(t, "this_is_a_multiple_camel", CamelToSnake("ThisIsAMultipleCamel"))
	AssertEqual(t, "", CamelToSnake(""))
}

func TestRound(t *testing.T) {
	AssertEqual(t, 5.0, Round(4.5))
	AssertEqual(t, 4.0, Round(4.49))
	AssertEqual(t, 5.0, Round(4.7))
	AssertEqual(t, 5.0, Round(5.1))
	AssertEqual(t, 5.0, Round(5.0))
	AssertEqual(t, 0.0, Round(0.49))
	AssertEqual(t, 0.0, Round(-0.49))
	AssertEqual(t, -1.0, Round(-1.0))
	AssertEqual(t, -1.0, Round(-1.5))
	AssertEqual(t, -2.0, Round(-1.51))
}

func TestRoundToDecimal(t *testing.T) {
	AssertEqual(t, 5.0, RoundToDecimal(4.5, 0))
	AssertEqual(t, 4.5, RoundToDecimal(4.5, 1))
	AssertEqual(t, 4.5, RoundToDecimal(4.5, 2))
	AssertEqual(t, 4.5, RoundToDecimal(4.49, 1))
	AssertEqual(t, 4.49, RoundToDecimal(4.49, 2))
	AssertEqual(t, 4.45, RoundToDecimal(4.449, 2))
	AssertEqual(t, 12433.2342, RoundToDecimal(12433.2341698323239823, 4))
}

func TestNanoToMicro(t *testing.T) {
	AssertEqual(t, int64(1452522568079810), NanoToMicro(1452522568079810000))
	AssertEqual(t, int64(1452522568079810), NanoToMicro(1452522568079810260))
	AssertEqual(t, int64(1452522568079811), NanoToMicro(1452522568079810500))
	AssertEqual(t, int64(1452522568079810), NanoToMicro(1452522568079810499))
	AssertEqual(t, int64(1452522568079811), NanoToMicro(1452522568079810599))
}

func TestMicroToNano(t *testing.T) {
	AssertEqual(t, int64(1452522568079810000), MicroToNano(1452522568079810))
	AssertEqual(t, int64(1452522568079811000), MicroToNano(1452522568079811))
}

func TestIsZeroStruct(t *testing.T) {
	AssertEqual(t, false, IsZeroStruct(nil))
	AssertEqual(t, true, IsZeroStruct(TestStruct{}))
	AssertEqual(t, true, IsZeroStruct(TestStruct{ID: 0}))
	AssertEqual(t, true, IsZeroStruct(TestStruct{ID: 0, Title: ""}))
	AssertEqual(t, false, IsZeroStruct(TestStruct{ID: 10}))
	AssertEqual(t, false, IsZeroStruct(TestStruct{Title: "Hello"}))
	s := "Hello, Test"
	AssertEqual(t, false, IsZeroStruct(TestStruct{Description: &s}))
	AssertEqual(t, false, IsZeroStruct(&TestStruct{}))
	var i interface{}
	AssertEqual(t, false, IsZeroStruct(i))
}
