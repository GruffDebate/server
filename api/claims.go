package api

import (
	"net/http"
	"reflect"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
)

func GetClaim(c echo.Context) error {
	ctx := ServerContext(c)

	id := c.Param("id")

	// TODO: as of date
	var err gruff.Error
	claim := gruff.Claim{}
	claim.ID = id
	err = claim.Load(ctx)
	if err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, claim)
}

func ListClaims(which string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := ServerContext(c)
		db := ctx.Arango.DB

		var claim gruff.Claim
		var claims []gruff.Claim
		var params gruff.ArangoQueryParameters
		var bindVars map[string]interface{}
		var query string
		switch which {
		case "top":
			query = claim.QueryForTopLevelClaims(params)
		case "new":
			query = gruff.DefaultListQuery(&claim, claim.DefaultQueryParameters())
		default:
			return AddError(ctx, c, gruff.NewNotFoundError("Not found"))
		}

		cursor, err := db.Query(ctx.Context, query, bindVars)
		if err != nil {
			return AddError(ctx, c, gruff.NewServerError(err.Error()))
		}
		defer cursor.Close()
		for cursor.HasMore() {
			claim := gruff.Claim{}
			_, err := cursor.ReadDocument(ctx.Context, &claim)
			if err != nil {
				return AddError(ctx, c, gruff.NewServerError(err.Error()))
			}
			claims = append(claims, claim)
		}

		return c.JSON(http.StatusOK, claims)
	}
}

func AddContext(c echo.Context) error {
	ctx := ServerContext(c)

	parentId := c.Param("parentId")
	id := c.Param("id")

	claim := gruff.Claim{}
	claim.Key = parentId
	if err := claim.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	context, err := gruff.GetArangoObject(ctx, reflect.TypeOf(gruff.Context{}), id)
	if err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.AddContext(ctx, *context.(*gruff.Context)); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, claim)
}

func RemoveContext(c echo.Context) error {
	ctx := ServerContext(c)

	parentId := c.Param("parentId")
	id := c.Param("id")

	claim := gruff.Claim{}
	claim.Key = parentId
	if err := claim.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.RemoveContext(ctx, id); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, claim)
}

/*
func SetScore(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return AddError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	user := ctx.UserContext
	if err != nil {
		return AddError(ctx, c, gruff.NewUnauthorizedError(err.Error()))
	}

	paths := strings.Split(c.Path(), "/")
	scoreType := paths[len(paths)-1]

	var claim bool
	var target, item interface{}

	switch scoreType {
	case "truth":
		claim = true
		target = &gruff.Claim{}
		item = &gruff.ClaimOpinion{UserID: user.ID, ClaimID: id}
	case "strength":
		claim = false
		target = &gruff.Argument{}
		item = &gruff.ArgumentOpinion{UserID: user.ID, ArgumentID: id}
	default:
		return AddError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	err = db.Where("id = ?", id).First(target).Error
	if err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	data := map[string]interface{}{}
	if err := c.Bind(&data); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}
	var score float64
	if val, ok := data["score"]; ok {
		score = val.(float64)
	}

	status := http.StatusCreated
	db = ctx.Database
	db = db.Where("user_id = ?", user.ID)
	if claim {
		db = db.Where("claim_id = ?", id)
	} else {
		db = db.Where("argument_id = ?", id)
	}
	if err := db.First(item).Error; err != nil {
		setScore(item, scoreType, score)
		db = ctx.Database
		err = db.Create(item).Error
		if err != nil {
			return AddError(ctx, c, gruff.NewServerError(err.Error()))
		}
	} else {
		setScore(item, scoreType, score)
		db = ctx.Database
		err = db.Save(item).Error
		if err != nil {
			return AddError(ctx, c, gruff.NewServerError(err.Error()))
		}
		status = http.StatusAccepted
	}

	switch scoreType {
	case "truth":
		target.(*gruff.Claim).UpdateTruth(ServerContext(c))
	case "strength":
		target.(*gruff.Argument).UpdateStrength(ServerContext(c))
	}

	return c.JSON(status, item)
}

func setScore(item interface{}, field string, score float64) {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = reflect.ValueOf(item).Elem()
	}
	f := v.FieldByName(strings.Title(field))
	f.Set(reflect.ValueOf(score))
}
*/
