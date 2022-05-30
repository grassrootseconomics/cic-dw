package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func handleProtectedResource(c echo.Context) error {
	return c.String(http.StatusOK, "unlocked")
}
