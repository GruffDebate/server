package api

import (
	"fmt"
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
	filters := gruff.BindVars{}
	var query string
	if userID != "" && gruff.IsVersionedModel(ctx.Type) {
		filters["creator"] = userID
		query = gruff.DefaultListQueryForUser(item, params)
	} else {
		query = gruff.DefaultListQuery(item, params)
	}

	items := []interface{}{}
	if err := gruff.FindArangoObjects(ctx, query, filters, &items); err != nil {
		return AddError(ctx, c, err)
	}

	ctx.Payload["results"] = items
	return c.JSON(http.StatusOK, ctx.Payload)
}

func Create(c echo.Context) error {
	ctx := ServerContext(c)

	if !gruff.IsArangoObject(reflect.PtrTo(ctx.Type)) {
		return AddError(ctx, c, gruff.NewServerError("This item isn't compatible with this request"))
	}

	item := reflect.New(ctx.Type).Interface()
	if err := c.Bind(item); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	userID := ActiveUserID(c, ctx)
	if userID != "" {
		gruff.SetUserID(item, userID)
	}

	err := item.(gruff.ArangoObject).Create(ctx)
	if err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusCreated, item)
}

func Update(c echo.Context) error {
	ctx := ServerContext(c)

	if !gruff.IsArangoObject(reflect.PtrTo(ctx.Type)) {
		return AddError(ctx, c, gruff.NewServerError(fmt.Sprintf("This item isn't compatible with this request")))
	}

	id := c.Param("id")
	if id == "" {
		return AddError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	item, err := loadItem(c, id)
	if err != nil {
		return AddError(ctx, c, err)
	}

	obj := item.(gruff.ArangoObject)

	updates := gruff.Updates{}
	if err := c.Bind(&updates); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if gruff.IsVersionedModel(ctx.Type) {
		if err := validateKeyParameter(c, obj, updates); err != nil {
			return AddError(ctx, c, err)
		}
	}

	err = obj.Update(ctx, updates)
	if err != nil {
		return AddError(ctx, c, err)
	}

	if gruff.IsLoader(reflect.PtrTo(ctx.Type)) {
		loader := item.(gruff.Loader)
		if err := loader.LoadFull(ctx); err != nil {
			return AddError(ctx, c, err)
		}
	}

	return c.JSON(http.StatusOK, item)
}

func Delete(c echo.Context) error {
	ctx := ServerContext(c)

	if !gruff.IsArangoObject(reflect.PtrTo(ctx.Type)) {
		return AddError(ctx, c, gruff.NewServerError("This item isn't compatible with this request"))
	}

	id := c.Param("id")
	if id == "" {
		return AddError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	item, err := loadItem(c, id)
	if err != nil {
		return AddError(ctx, c, err)
	}

	obj := item.(gruff.ArangoObject)

	err = obj.Delete(ctx)
	if err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, item)
}

// TODO: Handle query date
func Get(c echo.Context) error {
	ctx := ServerContext(c)

	id := c.Param("id")
	if id == "" {
		return AddError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	result, err := loadItem(c, id)
	if err != nil {
		return AddError(ctx, c, err)
	}

	if gruff.IsRestrictor(ctx.Type) {
		r := result.(gruff.Restrictor)
		canView, err := r.UserCanView(ctx)
		if err != nil {
			return AddError(ctx, c, err)
		}
		if !canView {
			return AddError(ctx, c, gruff.NewPermissionError("You do not have permission to view this item"))
		}

	}

	return c.JSON(http.StatusOK, result)
}

func loadItem(c echo.Context, id string) (interface{}, gruff.Error) {
	ctx := ServerContext(c)
	var result interface{}
	if gruff.IsLoader(reflect.PtrTo(ctx.Type)) {
		item := reflect.New(ctx.Type).Interface()
		loader := item.(gruff.Loader)

		if gruff.IsVersionedModel(ctx.Type) {
			gruff.SetID(loader, id)
		}

		err := loader.LoadFull(ctx)
		if err != nil {
			return result, err
		}
		result = loader
	} else {
		if err := gruff.LoadArangoObject(ctx, result, id); err != nil {
			return result, err
		}
	}
	return result, nil
}

func validateKeyParameter(c echo.Context, item gruff.ArangoObject, paramses ...map[string]interface{}) gruff.Error {
	params := map[string]interface{}{}
	if len(paramses) > 0 {
		params = paramses[0]
	} else {
		if err := c.Bind(&params); err != nil {
			return gruff.NewServerError(err.Error())
		}
	}

	var key string
	var ok bool
	if key, ok = params["_key"].(string); !ok {
		return gruff.NewBusinessError("Key: non zero value required;")
	}

	if item.ArangoKey() != key {
		return gruff.NewBusinessError("The key does not match the value for the item you are changing. The item might have already been changed by someone else.")
	}

	return nil
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
