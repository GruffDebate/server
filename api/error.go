package api

import (
	"net/http"

	"github.com/GruffDebate/server/gruff"
	"github.com/labstack/echo"
)

func AddPayloadWarning(c echo.Context, payload map[string]interface{}, code int, message string) error {
	payload["warning"] = map[string]interface{}{"code": code, "message": message}
	return nil
}

func AddPayloadError(c echo.Context, payload map[string]interface{}, code int, message string) error {
	return c.JSON(http.StatusBadRequest, map[string]interface{}{"code": code, "message": message})
}

func AddPermissionError(c echo.Context, payload map[string]interface{}, code int, message string) error {
	return c.JSON(http.StatusForbidden, map[string]interface{}{"code": code, "message": message})
}

func AddNotFoundError(c echo.Context, payload map[string]interface{}, code int, message string) error {
	return c.JSON(http.StatusNotFound, map[string]interface{}{"code": code, "message": message})
}

func AddServerError(c echo.Context, payload map[string]interface{}, code int, message string) error {
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{"code": code, "message": message})
}

func AddGruffError(ctx *gruff.ServerContext, c echo.Context, errGruff gruff.GruffError) error {
	var err error

	code := errGruff.Code()
	if errGruff.Subcode() != 0 {
		code = errGruff.Subcode()
	}

	params := errGruff.Data()
	params["code"] = errGruff.Code()
	params["subcode"] = errGruff.Subcode()

	switch errGruff.Code() {
	case 300:
		err = AddPayloadWarning(c, ctx.Payload, code, errGruff.Error())
	case 400:
		err = AddPayloadError(c, ctx.Payload, code, errGruff.Error())
	case 403:
		err = AddPermissionError(c, ctx.Payload, code, errGruff.Error())
	case 404:
		err = AddNotFoundError(c, ctx.Payload, code, errGruff.Error())
	default:
		err = AddServerError(c, ctx.Payload, code, errGruff.Error())
	}

	return err
}
