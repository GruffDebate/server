package gruff

import (
	"context"
	"reflect"
	"time"

	"github.com/GruffDebate/server/support"
)

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
	RequestAt   *time.Time
}

// TODO: Test
func (ctx ServerContext) UserLoggedIn() bool {
	return ctx.UserContext.ArangoKey() != ""
}

func (ctx *ServerContext) RequestTime() time.Time {
	if ctx.RequestAt == nil {
		ctx.RequestAt = support.TimePtr(time.Now())
	}
	return *ctx.RequestAt
}

func (ctx ServerContext) Rollback() Error {
	return ctx.Arango.Rollback()
}
