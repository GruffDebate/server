package api

import (
	"net/http"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
)

func GetArgument(c echo.Context) error {
	ctx := ServerContext(c)

	id := c.Param("id")

	// TODO: as of date
	var err gruff.Error
	argument := gruff.Argument{}
	argument.ID = id
	err = argument.Load(ctx)
	if err != nil {
		return AddError(ctx, c, err)
	}

	/*
		db = db.Preload("Claim.Links")
		db = db.Preload("Claim.Contexts")
		db = db.Preload("Claim.Values")
		db = db.Preload("Claim.Tags")
		db = db.Preload("TargetClaim.Links")
		db = db.Preload("TargetClaim.Contexts")
		db = db.Preload("TargetClaim.Values")
		db = db.Preload("TargetClaim.Tags")
		db = db.Preload("TargetArgument.Claim")
		err = db.Where("id = ?", id).First(&argument).Error
		if err != nil {
			return AddError(ctx, c, gruff.NewServerError(err.Error()))
		}

		pro := []gruff.Argument{}
		db = ctx.Database
		db = db.Preload("Claim")
		db = db.Where("type = ?", gruff.ARGUMENT_FOR)
		db = db.Where("target_argument_id = ?", id)
		db = db.Scopes(gruff.OrderByBestArgument)
		err = db.Find(&pro).Error
		if err != nil {
			return AddError(ctx, c, gruff.NewServerError(err.Error()))
		}
		argument.Pro = pro

		con := []gruff.Argument{}
		db = ctx.Database
		db = db.Preload("Claim")
		db = db.Where("type = ?", gruff.ARGUMENT_AGAINST)
		db = db.Where("target_argument_id = ?", id)
		db = db.Scopes(gruff.OrderByBestArgument)
		err = db.Find(&con).Error
		if err != nil {
			return AddError(ctx, c, gruff.NewServerError(err.Error()))
		}
		argument.Con = con
	*/

	return c.JSON(http.StatusOK, argument)
}

func CreateArgument(c echo.Context) error {
	ctx := ServerContext(c)

	arg := gruff.Argument{Claim: &gruff.Claim{}}
	if err := c.Bind(&arg); err != nil {
		return AddError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	arg.CreatedByID = ctx.UserContext.ArangoID()

	return c.JSON(http.StatusCreated, arg)
}

/*
func MoveArgument(c echo.Context) error {
	ctx := ServerContext(c)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return AddError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	newID, err := uuid.Parse(c.Param("newId"))
	if err != nil {
		return AddError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	t, err := strconv.Atoi(c.Param("type"))
	if err != nil {
		return AddError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	objType, err := strconv.Atoi(c.Param("targetType"))
	if err != nil {
		return AddError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	arg := gruff.Argument{}
	if err := db.Where("id = ?", id).First(&arg).Error; err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if err := (&arg).MoveTo(ServerContext(c), newID, t, objType); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	ctx.Payload["results"] = arg
	return c.JSON(http.StatusOK, ctx.Payload)
}
*/
