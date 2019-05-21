package api

import (
	"net/http"
	"strconv"

	"github.com/GruffDebate/server/gruff"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func GetArgument(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	argument := gruff.Argument{}

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
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	pro := []gruff.Argument{}
	db = ctx.Database
	db = db.Preload("Claim")
	db = db.Where("type = ?", gruff.ARGUMENT_FOR)
	db = db.Where("target_argument_id = ?", id)
	db = db.Scopes(gruff.OrderByBestArgument)
	err = db.Find(&pro).Error
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
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
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}
	argument.Con = con

	return c.JSON(http.StatusOK, argument)
}

func CreateArgument(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	arg := gruff.Argument{Claim: &gruff.Claim{}}
	if err := c.Bind(&arg); err != nil {
		return AddGruffError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	arg.CreatedByID = ctx.UserContext.ID

	if arg.ClaimID == uuid.Nil {
		ctxIds := arg.Claim.ContextIDs

		// First create a new Claim for this argument
		claim := gruff.Claim{Title: arg.Title, Description: arg.Description}
		claim.CreatedByID = arg.CreatedByID
		if arg.Claim.Title != "" {
			claim.Title = arg.Claim.Title
			if arg.Title == "" {
				arg.Title = arg.Claim.Title
			}
		}
		if arg.Claim.Description != "" {
			claim.Description = arg.Claim.Description
		}
		valerr := DefaultValidationForCreate(ctx, c, claim)
		if valerr != nil {
			return AddGruffError(ctx, c, gruff.NewServerError(valerr.Error()))
		}
		err := db.Create(&claim).Error
		if err != nil {
			return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
		}
		arg.ClaimID = claim.ID
		arg.Claim = &claim

		for _, ctxID := range ctxIds {
			db.Exec("INSERT INTO claim_contexts (claim_id, context_id) VALUES (?, ?)", claim.ID, ctxID)
		}
	} else {
		arg.Claim = nil
	}

	valerr := DefaultValidationForCreate(ctx, c, arg)
	if valerr != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(valerr.Error()))
	}

	if dberr := db.Set("gorm:save_associations", false).Create(&arg).Error; dberr != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(dberr.Error()))
	}

	return c.JSON(http.StatusCreated, arg)
}

func MoveArgument(c echo.Context) error {
	ctx := ServerContext(c)
	db := ctx.Database

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	newID, err := uuid.Parse(c.Param("newId"))
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	t, err := strconv.Atoi(c.Param("type"))
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	objType, err := strconv.Atoi(c.Param("targetType"))
	if err != nil {
		return AddGruffError(ctx, c, gruff.NewNotFoundError(err.Error()))
	}

	arg := gruff.Argument{}
	if err := db.Where("id = ?", id).First(&arg).Error; err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	if err := (&arg).MoveTo(ServerContext(c), newID, t, objType); err != nil {
		return AddGruffError(ctx, c, gruff.NewServerError(err.Error()))
	}

	ctx.Payload["results"] = arg
	return c.JSON(http.StatusOK, ctx.Payload)
}
