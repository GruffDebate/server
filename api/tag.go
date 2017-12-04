package api

import (
	"github.com/labstack/echo"
	"net/http"
)

func Tags(c echo.Context) error {
	return c.String(http.StatusOK, "Tags")
}
