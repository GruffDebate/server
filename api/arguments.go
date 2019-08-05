package api

import (
	"net/http"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
)

func MoveArgument(c echo.Context) error {
	ctx := ServerContext(c)

	id := c.Param("id")
	if id == "" {
		return AddError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	targetId := c.Param("targetId")
	if targetId == "" {
		return AddError(ctx, c, gruff.NewNotFoundError("Not Found"))
	}

	params := map[string]interface{}{}
	if err := c.Bind(&params); err != nil {
		return AddError(ctx, c, gruff.NewServerError(err.Error()))
	}

	var pro, ok bool
	if pro, ok = params["pro"].(bool); !ok {
		return AddError(ctx, c, gruff.NewBusinessError("Pro: non zero value required;"))
	}

	arg := gruff.Argument{}
	arg.ID = id
	if err := arg.Load(ctx); err != nil {
		return AddError(ctx, c, err)
	}

	if err := validateKeyParameter(c, &arg); err != nil {
		return AddError(ctx, c, err)
	}

	var target gruff.ArangoObject
	if c.Param("type") == "claims" {
		claim := gruff.Claim{}
		claim.ID = targetId
		if err := claim.Load(ctx); err != nil {
			return AddError(ctx, c, err)
		}
		target = &claim
	} else {
		targ := gruff.Argument{}
		targ.ID = targetId
		if err := targ.Load(ctx); err != nil {
			return AddError(ctx, c, err)
		}
		target = &targ
	}

	if err := arg.MoveTo(ctx, target, pro); err != nil {
		return AddError(ctx, c, err)
	}

	ctx.Payload["results"] = arg
	return c.JSON(http.StatusOK, ctx.Payload)
}
