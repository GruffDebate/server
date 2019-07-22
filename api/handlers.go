package api

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
)

func List(c echo.Context) error {
	ctx := ServerContext(c)

	item := reflect.New(ctx.Type).Interface().(gruff.ArangoObject)

	params := item.DefaultQueryParameters()
	params = params.Merge(GetListParametersFromRequest(c))

	items, err := gruff.ListArangoObjects(ctx, ctx.Type, gruff.DefaultListQuery(item, params), map[string]interface{}{})
	if err != nil {
		return AddGruffError(ctx, c, err)
	}

	if ctx.Payload["ct"] != nil {
		ctx.Payload["results"] = items
		return c.JSON(http.StatusOK, ctx.Payload)
	}

	return c.JSON(http.StatusOK, items)
}

func Create(c echo.Context) error {
	ctx := ServerContext(c)

	if !gruff.IsCreator(ctx.Type) {
		return AddGruffError(ctx, c, gruff.NewServerError("This item doesn't implement the Creator interface"))
	}

	item := reflect.New(ctx.Type).Interface()
	if err := c.Bind(item); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	err := item.(gruff.Creator).Create(ctx)
	if err != nil {
		return AddGruffError(ctx, c, err)
	}

	return c.JSON(http.StatusCreated, item)
}

func Get(c echo.Context) error {
	ctx := ServerContext(c)

	id := c.Param("id")
	if id == "" {
		return AddGruffError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	if !gruff.IsLoader(ctx.Type) {
		return AddGruffError(ctx, c, gruff.NewServerError("This item doesn't implement the Loader interface"))
	}

	item := reflect.New(ctx.Type).Interface().(gruff.Loader)

	if gruff.IsIdentifier(ctx.Type) {
		ident, err := gruff.GetIdentifier(item)
		if err != nil {
			return AddGruffError(ctx, c, err)
		}
		// TODO: This is probably NOT going to change the original - this is probably just changing a copy :(
		ident.ID = id
	}

	err := item.Load(ctx)
	if err != nil {
		return AddGruffError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, item)
}

func GetListParametersFromRequest(c echo.Context) gruff.ArangoQueryParameters {
	params := gruff.ArangoQueryParameters{}

	start, _ := strconv.Atoi(c.QueryParam("start"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if start >= 0 {
		params.Offset = &start

	}

	if limit > 0 {
		params.Limit = &limit
	}

	return params
}

func itemsOrEmptySlice(t reflect.Type, items interface{}) interface{} {
	if reflect.ValueOf(items).IsNil() {
		items = reflect.MakeSlice(reflect.SliceOf(t), 0, 0)
	}
	return items
}
