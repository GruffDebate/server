package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/GruffDebate/server/gruff"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func GetClaim(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	claim := gruff.Claim{}

	db = db.Preload("Links")
	db = db.Preload("Contexts")
	db = db.Preload("Values")
	db = db.Preload("Tags")
	db = db.Where("id = ?", id)
	err = db.First(&claim).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	proArgs := []gruff.Argument{}
	db = ctx.Database
	db = db.Preload("Claim.Contexts")
	db = db.Where("type = ?", gruff.ARGUMENT_FOR)
	db = db.Where("target_claim_id = ?", id)
	db = db.Scopes(gruff.OrderByBestArgument)
	err = db.Find(&proArgs).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}
	claim.Pro = proArgs

	conArgs := []gruff.Argument{}
	db = ctx.Database
	db = db.Preload("Claim.Contexts")
	db = db.Where("type = ?", gruff.ARGUMENT_AGAINST)
	db = db.Where("target_claim_id = ?", id)
	db = db.Scopes(gruff.OrderByBestArgument)
	err = db.Find(&conArgs).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}
	claim.Con = conArgs

	return c.JSON(http.StatusOK, claim)
}

func ListTopClaims(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	claims := []gruff.Claim{}

	db = BasicJoins(ctx, c, db)
	db = db.Where("0 = (SELECT COUNT(*) FROM arguments WHERE claim_id = claims.id)")
	db = BasicPaging(ctx, c, db)

	err := db.Find(&claims).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if ctx.Payload["ct"] != nil {
		ctx.Payload["results"] = claims
		return c.JSON(http.StatusOK, ctx.Payload)
	}

	return c.JSON(http.StatusOK, claims)
}

func SetScore(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	user := ctx.UserContext
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewUnauthorizedError(err.Error()))
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
		return AddGruffError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	err = db.Where("id = ?", id).First(target).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	data := map[string]interface{}{}
	if err := c.Bind(&data); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
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
			return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
		}
	} else {
		setScore(item, scoreType, score)
		db = ctx.Database
		err = db.Save(item).Error
		if err != nil {
			return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
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
