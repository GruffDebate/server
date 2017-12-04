package api

import (
	"github.com/labstack/echo"
	"net/http"
)

func Contexts(c echo.Context) error {
	return c.String(http.StatusOK, "Contexts")
}
