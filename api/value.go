package api

import (
	"github.com/labstack/echo"
	"net/http"
)

func Values(c echo.Context) error {
	return c.String(http.StatusOK, "Values")
}
