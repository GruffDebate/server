package api

import (
	"github.com/labstack/echo"
	"net/http"
)

func Home(c echo.Context) error {
	hello := map[string]interface{}{"status": "GRUFF API"}
	return c.JSON(http.StatusOK, hello)
}
