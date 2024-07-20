package common

import (
	"net/http"
	"net/url"

	"github.com/JamesTiberiusKirk/recipe-cms/site/session"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(session *session.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !session.IsAuthenticated(c, false) {
				return c.Redirect(http.StatusSeeOther, "/auth/login?source="+url.QueryEscape(c.Request().URL.String()))
			}
			return next(c)
		}
	}
}
