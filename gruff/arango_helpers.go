package gruff

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/GruffDebate/server/support"
	arango "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type ArangoObject interface {
	CollectionName() string
	ArangoKey() string
	ArangoID() string
	DefaultQueryParameters() ArangoQueryParameters
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

func OpenArangoConnection() (arango.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{os.Getenv("ARANGO_ENDPOINT")},
	})
	if err != nil {
		return nil, err
	}
	conn, err = conn.SetAuthentication(arango.BasicAuthentication(os.Getenv("ARANGO_USER"), os.Getenv("ARANGO_PASS")))
	if err != nil {
		return nil, err
	}
	db, err := arango.NewClient(arango.ClientConfig{
		Connection: conn,
	})

	return db, err
}

func OpenArangoDatabase(client arango.Client) (arango.Database, error) {
	ctx := context.Background()
	db, err := client.Database(ctx, os.Getenv("ARANGO_DB"))
	return db, err
}

func CloseCursor(cursor arango.Cursor) {
	if cursor != nil {
		cursor.Close()
	}
}

// QUERY PARAMETERS
// By convention, simple queries will use "obj" as the collection item reference
// in order to simplify the basic use cases
var DEFAULT_QUERY_PARAMETERS = ArangoQueryParameters{
	Sort:   support.StringPtr("obj.start DESC"),
	Offset: support.IntPtr(0),
	Limit:  support.IntPtr(20),
	Return: support.StringPtr("obj"),
}

type ArangoQueryParameters struct {
	Sort   *string
	Offset *int
	Limit  *int
	Return *string
}

func (aqp ArangoQueryParameters) Merge(params ArangoQueryParameters) ArangoQueryParameters {
	merged := ArangoQueryParameters{
		Sort:   aqp.Sort,
		Offset: aqp.Offset,
		Limit:  aqp.Limit,
		Return: aqp.Return,
	}

	if params.Sort != nil {
		merged.Sort = params.Sort
	}
	if params.Offset != nil {
		merged.Offset = params.Offset
	}
	if params.Limit != nil {
		merged.Limit = params.Limit
	}
	if params.Return != nil {
		merged.Return = params.Return
	}

	return merged
}

func (aqp ArangoQueryParameters) Apply(query string) string {
	var sort, limit, ret string

	queryEnd := " %s %s %s"
	query = fmt.Sprintf("%s%s", query, queryEnd)

	if aqp.Sort != nil {
		sort = fmt.Sprintf("SORT %s", *aqp.Sort)
	}
	if aqp.Limit != nil {
		var offset int
		if aqp.Offset != nil {
			offset = *aqp.Offset
		}
		limit = fmt.Sprintf("LIMIT %d, %d", offset, *aqp.Limit)
	}
	if aqp.Return != nil {
		ret = fmt.Sprintf("RETURN %s", *aqp.Return)
	}

	return fmt.Sprintf(query, sort, limit, ret)
}

// Default Finders

func DefaultListQuery(obj ArangoObject, params ArangoQueryParameters) string {
	query := fmt.Sprintf("FOR obj IN %s FILTER obj.end == null", obj.CollectionName())
	return params.Apply(query)
}

func DefaultListQueryForUser(obj ArangoObject, params ArangoQueryParameters) string {
	query := fmt.Sprintf("FOR obj IN %s FILTER obj.creator == @creator AND obj.end == null", obj.CollectionName())
	return params.Apply(query)
}

func ListArangoObjects(ctx *ServerContext, t reflect.Type, query string, bindVars map[string]interface{}) ([]interface{}, GruffError) {
	db := ctx.Arango.DB

	objs := []interface{}{}

	cursor, err := db.Query(ctx.Context, query, bindVars)
	if err != nil {
		return objs, NewServerError(err.Error())
	}
	defer cursor.Close()
	for cursor.HasMore() {
		obj := reflect.New(t).Interface()
		_, err := cursor.ReadDocument(ctx.Context, &obj)
		if err != nil {
			return objs, NewServerError(err.Error())
		}
		objs = append(objs, obj)
	}

	return objs, nil
}

func GetArangoObject(ctx *ServerContext, t reflect.Type, arangoKey string) (interface{}, GruffError) {
	obj := reflect.New(t).Interface().(ArangoObject)

	col, err := ctx.Arango.CollectionFor(obj)
	if err != nil {
		return nil, err
	}

	if _, dberr := col.ReadDocument(ctx.Context, arangoKey, obj); dberr != nil {
		return nil, NewNotFoundError(dberr.Error())
	}

	return obj, nil
}
