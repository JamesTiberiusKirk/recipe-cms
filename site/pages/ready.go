package pages

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Ready(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
