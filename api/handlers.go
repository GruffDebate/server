package api

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
)

// TODO: Handle query date
func List(c echo.Context) error {
	// TODO: check authorized
	ctx := ServerContext(c)

	item := reflect.New(ctx.Type).Interface().(gruff.ArangoObject)

	params := item.DefaultQueryParameters()
	params = params.Merge(GetListParametersFromRequest(c))

	userID := ActiveUserID(c, ctx)
	filters := map[string]interface{}{}
	var query string
	if userID != "" {
		filters["creator"] = userID
		query = gruff.DefaultListQueryForUser(item, params)
	} else {
		query = gruff.DefaultListQuery(item, params)
	}

	items, err := gruff.ListArangoObjects(ctx, ctx.Type, query, filters)
	if err != nil {
		return AddGruffError(ctx, c, err)
	}

	ctx.Payload["results"] = items
	return c.JSON(http.StatusOK, ctx.Payload)
}

func Create(c echo.Context) error {
	ctx := ServerContext(c)

	if !gruff.IsCreator(reflect.PtrTo(ctx.Type)) {
		return AddGruffError(ctx, c, gruff.NewServerError("This item doesn't implement the Creator interface"))
	}

	item := reflect.New(ctx.Type).Interface()
	if err := c.Bind(item); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	userID := ActiveUserID(c, ctx)
	if userID != "" {
		gruff.SetUserID(item, userID)
	}

	err := item.(gruff.Creator).Create(ctx)
	if err != nil {
		return AddGruffError(ctx, c, err)
	}

	return c.JSON(http.StatusCreated, item)
}

func Update(c echo.Context) error {
	ctx := ServerContext(c)

	if !gruff.IsUpdater(reflect.PtrTo(ctx.Type)) {
		return AddGruffError(ctx, c, gruff.NewServerError("This item doesn't implement the Updater interface"))
	}

	updates := map[string]interface{}{}
	if err := c.Bind(&updates); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	var key string
	var ok bool
	if key, ok = updates["_key"].(string); !ok {
		return AddGruffError(ctx, c, gruff.NewBusinessError("Key: non zero value required;"))
	}

	item := reflect.New(ctx.Type).Interface()

	// TODO check if end is null

	if err := gruff.SetKey(item, key); err != nil {
		return AddGruffError(ctx, c, err)
	}

	err := item.(gruff.Updater).Update(ctx, updates)
	if err != nil {
		return AddGruffError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, item)
}

// TODO: Handle query date
func Get(c echo.Context) error {
	ctx := ServerContext(c)

	id := c.Param("id")
	if id == "" {
		return AddGruffError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	var result interface{}
	if gruff.IsLoader(reflect.PtrTo(ctx.Type)) {
		item := reflect.New(ctx.Type).Interface()
		loader := item.(gruff.Loader)

		if gruff.IsIdentifier(ctx.Type) {
			ident, err := gruff.GetIdentifier(loader)
			if err != nil {
				return AddGruffError(ctx, c, err)
			}
			// TODO: This is probably NOT going to change the original - this is probably just changing a copy :(
			ident.ID = id
		}

		err := loader.LoadFull(ctx)
		if err != nil {
			return AddGruffError(ctx, c, err)
		}
		result = loader
	} else {
		item, err := gruff.GetArangoObject(ctx, ctx.Type, id)
		if err != nil {
			return AddGruffError(ctx, c, err)
		}
		result = item
	}

	if gruff.IsRestrictor(ctx.Type) {
		r := result.(gruff.Restrictor)
		canView, err := r.UserCanView(ctx)
		if err != nil {
			return AddGruffError(ctx, c, err)
		}
		if !canView {
			return AddGruffError(ctx, c, gruff.NewPermissionError("You do not have permission to view this item"))
		}

	}

	return c.JSON(http.StatusOK, result)
}

// TODO: GetQueryDateFromRequest
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

func ActiveUserID(c echo.Context, ctx *gruff.ServerContext) string {
	userID := ctx.UserContext.ArangoID()
	id := c.Param("userId")
	if id != "" {
		userID = id
	}
	return userID
}
