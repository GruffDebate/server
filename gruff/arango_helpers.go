package gruff

import (
	"context"

	arango "github.com/arangodb/go-driver"
)

type ArangoObject interface {
	CollectionName() string
	ArangoKey() string
	ArangoID() string
}

type ArangoContext struct {
	Context     context.Context
	DB          arango.Database
	Collections map[string]arango.Collection
}

func (ctx ArangoContext) Rollback() GruffError {
	return NewServerError("Not implemented yet")
}

func (ctx ArangoContext) Collection(name string) (arango.Collection, GruffError) {
	if ctx.Collections == nil {
		ctx.Collections = make(map[string]arango.Collection)
	}
	if col, ok := ctx.Collections[name]; ok {
		return col, nil
	}

	col, err := ctx.DB.Collection(ctx.Context, name)
	if err != nil {
		return col, NewServerError(err.Error())
	}

	ctx.Collections[name] = col
	return col, nil
}

func (ctx ArangoContext) CollectionFor(item ArangoObject) (arango.Collection, GruffError) {
	return ctx.Collection(item.CollectionName())
}
