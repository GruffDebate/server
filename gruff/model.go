package gruff

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/GruffDebate/server/support"
)

// TODO: Update this for ArangoDB
type Model struct {
	ID        uint64     `json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" settable:"false"`
}

type ReplaceMany struct {
	IDS []uint64 `json:"ids"`
}

type ServerContext struct {
	Context     context.Context
	Arango      ArangoContext
	Payload     map[string]interface{}
	Request     map[string]interface{}
	Type        reflect.Type
	ParentType  reflect.Type
	Test        bool
	UserContext User
	AppName     string
	Method      string
	Path        string
	Endpoint    string
	RequestID   string
}

func (ctx ServerContext) Rollback() GruffError {
	return ctx.Arango.Rollback()
}

func ModelToJson(model interface{}) string {
	j, err := json.Marshal(model)
	if err != nil {
		panic(fmt.Sprintf("Error %v encoding JSON for %v", err, model))
	}

	jsonStr := string(j)
	v := reflect.Indirect(reflect.ValueOf(model))
	ot := v.Type()
	t := ot
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		t = t.Elem()
	} else if t.Kind() == reflect.Interface {
		t = v.Elem().Type()
	}
	return jsonStr
}

func ModelToJsonMap(modl interface{}) map[string]interface{} {
	jsonStr := ModelToJson(modl)
	m := JsonToMap(jsonStr)
	return m
}

func JsonToMap(jsonStr string) map[string]interface{} {
	jsonMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		panic(fmt.Sprintf("Error %v unmarshaling JSON for %v", err, jsonStr))
	}

	return jsonMap
}

func JsonToMapArray(jsonStr string) []map[string]interface{} {
	var arr []map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &arr)

	if err != nil {
		panic(fmt.Sprintf("Error %v unmarshaling JSON for %v", err, jsonStr))
	}

	return arr
}

func JsonToModel(jsonStr string, item interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), &item)

	if err == nil {
		v := reflect.Indirect(reflect.ValueOf(item))
		ot := v.Type()
		t := ot
		if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
			t = t.Elem()
		} else if t.Kind() == reflect.Interface {
			t = v.Elem().Type()
		}
	}
	return err
}

func GetFieldByJsonTag(item interface{}, jsonKey string) (field *reflect.StructField, gerr GruffError) {
	data := map[string]interface{}{
		"type": reflect.TypeOf(item),
		"key":  jsonKey,
	}

	if jsonKey == "" || jsonKey == "-" {
		return nil, NewBusinessError("Invalid JSON key", data)
	}

	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, NewBusinessError("Cannot set value on nil item", data)
		}
		v = reflect.ValueOf(item).Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fKey := support.JsonName(f)
		if fKey == jsonKey {
			return &f, nil
		}
	}

	return nil, NewNotFoundError("field not found", data)
}

func SetByJsonTag(item interface{}, jsonKey string, newVal interface{}) GruffError {
	data := map[string]interface{}{
		"type": reflect.TypeOf(item),
		"key":  jsonKey,
		"val":  newVal,
	}

	if jsonKey == "" || jsonKey == "-" {
		return NewBusinessError("Invalid JSON key", data)
	}

	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return NewBusinessError("Cannot set value on nil item", data)
		}
		v = v.Elem()
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag
		fKey := support.JsonName(f)
		vField := v.Field(i)
		if fKey == jsonKey {
			if tag.Get("settable") == "false" {
				return NewPermissionError("field is unsettable", data)
			}
			destType := vField.Type()
			if destType.Kind() == reflect.Ptr {
				destType = destType.Elem()
			}
			support.SetValue(vField.Addr(), destType, newVal)
			return nil
		}
	}

	return NewNotFoundError("field not found", data)
}
