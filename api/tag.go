package api

import (
	"net/http"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
)

func Tags(c echo.Context) error {
	return c.String(http.StatusOK, "Tags")
}

func ListClaimsByTag(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id := c.Param("id")
	if id == "" {
		return AddGruffError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	tag := gruff.Tag{}
	err := db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	claims := []gruff.Claim{}

	db = BasicJoins(ctx, c, db)
	err = db.Model(&tag).Association("Claims").Find(&claims).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if ctx.Payload["ct"] != nil {
		ctx.Payload["results"] = claims
		return c.JSON(http.StatusOK, ctx.Payload)
	}

	return c.JSON(http.StatusOK, claims)
}
