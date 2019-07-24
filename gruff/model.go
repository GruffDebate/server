package gruff

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/GruffDebate/server/support"
	"github.com/google/uuid"
)

type Model struct {
	Key       string     `json:"_key"`
	CreatedAt time.Time  `json:"start"`
	UpdatedAt time.Time  `json:"mod"`
	DeletedAt *time.Time `json:"end" settable:"false"`
}

func (m *Model) PrepareForCreate(u User) {
	m.Key = uuid.New().String()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return
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

func SetJsonValuesOnStruct(item interface{}, values map[string]interface{}) GruffError {
	for key, value := range values {
		if err := SetByJsonTag(item, key, value); err != nil {
			return err
		}
	}
	return nil
}

func SetKey(item interface{}, key string) GruffError {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return NewBusinessError("Cannot set value on nil item")
		}
		v = v.Elem()
	}

	fv := v.FieldByName("Key")
	if fv.Kind() != reflect.String {
		return NewServerError("Item does not have a Key field")
	}

	fv.SetString(key)
	return nil
}

func SetUserID(item interface{}, id string) error {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return errors.New("Cannot set value on nil item")
		}
		v = reflect.ValueOf(item).Elem()
	}
	t := v.Type()
	if !TypeHasUserField(t) {
		return errors.New("Type does not have a reference to a user")
	}
	userField, _ := UserIDField(t)
	f := v.FieldByName(userField.Name)
	if f.Type().Kind() == reflect.Ptr {
		f.Set(reflect.ValueOf(&id))
	} else {
		f.Set(reflect.ValueOf(id))
	}
	return nil
}

func TypeHasUserField(t reflect.Type) bool {
	userField, _ := UserIDField(t)
	return userField != nil
}

func UserIDField(t reflect.Type) (field *reflect.StructField, dbFieldName string) {
	elemT := t
	if elemT.Kind() == reflect.Ptr {
		elemT = elemT.Elem()
	}

	f, found := elemT.FieldByName("CreatedByID")
	if found {
		field = &f
		dbFieldName = "creator"
	}
	return
}
