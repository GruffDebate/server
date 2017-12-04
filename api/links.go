package api

import (
	"github.com/labstack/echo"
	"net/http"
)

func Links(c echo.Context) error {
	return c.String(http.StatusOK, "Links")
}
