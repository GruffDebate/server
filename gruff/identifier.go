package gruff

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"time"

	"github.com/google/uuid"
)

type Identifier struct {
	ID          uuid.UUID  `json:"uuid" sql:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time  `json:"-" sql:"DEFAULT:current_timestamp"`
	UpdatedAt   time.Time  `json:"-" sql:"DEFAULT:current_timestamp"`
	DeletedAt   *time.Time `json:"-" settable:"false"`
	CreatedByID uint64     `json:"createdById"`
	CreatedBy   *User      `json:"createdBy"`
}

func (i Identifier) ValidateForCreate() GruffError {
	return ValidateStruct(i)
}

func (i Identifier) ValidateForUpdate() GruffError {
	return i.ValidateForCreate()
}

func (i Identifier) ValidateField(f string) GruffError {
	return ValidateStructField(i, f)
}

func SetCreatedByID(item interface{}, id uint64) error {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return errors.New("Cannot set value on nil item")
		}
		v = reflect.ValueOf(item).Elem()
	}
	f := v.FieldByName("Identifier")
	f = f.FieldByName("CreatedByID")
	f.Set(reflect.ValueOf(id))
	return nil
}

func (i *Identifier) GenerateUUID() uuid.UUID {
	i.ID = uuid.New()
	return i.ID
}

func IsIdentifier(t reflect.Type) bool {
	_, is := t.FieldByName("Identifier")
	return is
}

// NullableUUID wrapper to fix nullable UUID. See https://github.com/golang/go/issues/8415
type NullableUUID struct {
	UUID uuid.UUID
}

// Value implements Sql/Value so it can be converted to DB value
func (u *NullableUUID) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}

	return u.UUID.Value()
}

// Setting from DB via Gorm
func (u *NullableUUID) Scan(value interface{}) error {
	if u == nil {
		return nil
	}

	if value != nil {
		ub := value.([]byte)
		up, err := uuid.ParseBytes(ub)
		if err != nil {
			return err
		}
		u.UUID = up
	}

	return nil
}

// MarshalText helps convert to value for JSON
func (u *NullableUUID) MarshalText() ([]byte, error) {
	if u == nil {
		return nil, nil
	}
	return u.UUID.MarshalText()
}

// UnmarshalText helps convert to value for JSON
func (u *NullableUUID) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	parsed, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}

	u.UUID = parsed
	return nil
}

func NUUID(id uuid.UUID) *NullableUUID {
	return &NullableUUID{UUID: id}
}

func NewNUUID() *NullableUUID {
	return NUUID(uuid.New())
}
