package gruff

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type GruffError interface {
	Code() int
	Subcode() int
	Error() string
	Location() string
	Data() map[string]interface{}
	IsWarning() bool
}

const ERROR_CODE_WARNING int = 300
const ERROR_CODE_BUSINESS_ERROR int = 400
const ERROR_CODE_UNAUTHORIZED_ERROR int = 401
const ERROR_CODE_PERMISSION_ERROR int = 403
const ERROR_CODE_NOT_FOUND int = 404
const ERROR_CODE_SERVER_ERROR int = 500

const ERROR_SUBCODE_UNDEFINED_IGNORE int = -1000
const ERROR_SUBCODE_UNDEFINED int = -1999

const ERROR_SUBCODE_USERNAME_TAKEN int = -2002
const ERROR_SUBCODE_USERNAME_LENGTH int = -2003
const ERROR_SUBCODE_USERNAME_FORMAT int = -2004
const ERROR_SUBCODE_EMAIL_TAKEN int = -2005
const ERROR_SUBCODE_EMAIL_FORMAT int = -2006
const ERROR_SUBCODE_PASSWORD_LENGTH int = -2007
const ERROR_SUBCODE_PASSWORD_FORMAT int = -2008
const ERROR_SUBCODE_CREDENTIALS_INVALID int = -2009

type CoreError struct {
	ErrCode     int
	ErrSubcode  int
	Message     string
	ErrLocation string
	ErrData     map[string]interface{}
}

func (err CoreError) Code() int {
	return err.ErrCode
}

func (err CoreError) Subcode() int {
	return err.ErrSubcode
}

func (err CoreError) Error() string {
	return err.Message
}

func (err CoreError) Location() string {
	return err.ErrLocation
}

func (err CoreError) Data() map[string]interface{} {
	return err.ErrData
}

func (err CoreError) IsWarning() bool {
	return err.ErrCode == ERROR_CODE_WARNING
}

func NewGruffError(code, subcode int, location, msg string, data map[string]interface{}) GruffError {
	if subcode == 0 {
		if code == ERROR_CODE_BUSINESS_ERROR {
			subcode = ERROR_SUBCODE_UNDEFINED
		} else {
			subcode = ERROR_SUBCODE_UNDEFINED_IGNORE
		}
	}
	return CoreError{
		ErrCode:     code,
		ErrSubcode:  subcode,
		Message:     msg,
		ErrLocation: location,
		ErrData:     data,
	}
}

func NewWarning(msg string, opts ...interface{}) GruffError {
	return newElipsisError(ERROR_CODE_WARNING, msg, opts...)
}

func NewUnauthorizedError(msg string, opts ...interface{}) GruffError {
	return newElipsisError(ERROR_CODE_UNAUTHORIZED_ERROR, msg, opts...)
}

func NewBusinessError(msg string, opts ...interface{}) GruffError {
	return newElipsisError(ERROR_CODE_BUSINESS_ERROR, msg, opts...)
}

func NewPermissionError(msg string, opts ...interface{}) GruffError {
	return newElipsisError(ERROR_CODE_PERMISSION_ERROR, msg, opts...)
}

func NewNotFoundError(msg string, opts ...interface{}) GruffError {
	return newElipsisError(ERROR_CODE_NOT_FOUND, msg, opts...)
}

func NewServerError(msg string, opts ...interface{}) GruffError {
	if strings.Contains(msg, "uix_users_email") {
		msg = "Email is already in use"
		opts = []interface{}{ERROR_SUBCODE_EMAIL_TAKEN}
	}

	return newElipsisError(ERROR_CODE_SERVER_ERROR, msg, opts...)
}

func newElipsisError(code int, msg string, opts ...interface{}) GruffError {
	subcode := 0
	data := map[string]interface{}{}

	for _, opt := range opts {
		switch reflect.TypeOf(opt).Kind() {
		case reflect.Int:
			subcode = opt.(int)
		case reflect.Map:
			data = opt.(map[string]interface{})

		}
	}

	return NewGruffError(code, subcode, ParentCallerInfo(), msg, data)
}

func ParentCallerInfo() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}
