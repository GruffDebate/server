package api

import (
	"fmt"
	"net/http"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
)

func ListClaims(which string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := ServerContext(c)

		var claim gruff.Claim
		var claims []gruff.Claim
		var bindVars gruff.BindVars
		var query string

		params := claim.DefaultQueryParameters()
		params = params.Merge(GetListParametersFromRequest(c))

		switch which {
		case "top":
			query = claim.QueryForTopLevelClaims(params)
		case "new":
			query = gruff.DefaultListQuery(&claim, claim.DefaultQueryParameters())
		default:
			return AddError(ctx, c, gruff.NewNotFoundError(fmt.Sprintf("Not found")))
		}

		err := gruff.FindArangoObjects(ctx, query, bindVars, &claims)
		if err != nil {
			return AddError(ctx, c, err)
		}
		return c.JSON(http.StatusOK, claims)
	}
}

func AddContext(c echo.Context) error {
	ctx := ServerContext(c)

	parentId := c.Param("parentId")
	id := c.Param("id")

	claim := gruff.Claim{}
	claim.ID = parentId
	if err := claim.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	if err := validateKeyParameter(c, &claim); err != nil {
		return AddError(ctx, c, err)
	}

	context := gruff.Context{}
	if err := gruff.LoadArangoObject(ctx, &context, id); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.AddContext(ctx, context); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.LoadFull(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, claim)
}

func RemoveContext(c echo.Context) error {
	ctx := ServerContext(c)

	parentId := c.Param("parentId")
	id := c.Param("id")

	claim := gruff.Claim{}
	claim.ID = parentId
	if err := claim.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.RemoveContext(ctx, id); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.LoadFull(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, claim)
}

func ConvertClaimToMultiPremise(c echo.Context) error {
	ctx := ServerContext(c)

	id := c.Param("id")

	claim := gruff.Claim{}
	claim.ID = id
	if err := claim.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	if err := validateKeyParameter(c, &claim); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.ConvertToMultiPremise(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.LoadFull(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, claim)
}

func AddPremise(c echo.Context) error {
	ctx := ServerContext(c)

	parentId := c.Param("parentId")
	id := c.Param("id")

	claim := gruff.Claim{}
	claim.ID = parentId
	if err := claim.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	if err := validateKeyParameter(c, &claim); err != nil {
		return AddError(ctx, c, err)
	}

	premise := gruff.Claim{}
	premise.ID = id
	if err := premise.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.AddPremise(ctx, &premise); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.LoadFull(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	return c.JSON(http.StatusOK, claim)
}

func RemovePremise(c echo.Context) error {
	ctx := ServerContext(c)

	parentId := c.Param("parentId")
	id := c.Param("id")

	claim := gruff.Claim{}
	claim.ID = parentId
	if err := claim.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.RemovePremise(ctx, id); err != nil {
		return AddError(ctx, c, err)
	}

	if err := claim.LoadFull(ctx); err != nil {
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
