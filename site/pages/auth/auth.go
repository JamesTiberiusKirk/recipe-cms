package auth

import (
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/site/session"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	sessions *session.Manager
}

func InitAuthHandler(app *echo.Group, s *session.Manager) {
	h := AuthHandler{
		sessions: s,
	}

	app.GET("/login", common.UseTemplContext(h.LoginPage))
	app.POST("/login", common.UseTemplContext(h.LoginPage))

	app.GET("/logout", common.UseTemplContext(h.Logout))
}

func (h *AuthHandler) LoginPage(c *common.TemplContext) error {
	props := loginPageProps{}

	if c.Request().Method == http.MethodPost {
		// TODO: actually tie this in with the db lol

		props.username = c.FormValue("username")
		password := c.FormValue("password")

		props.loginAttempted = true

		if password != "test" {
			logrus.Info("unauthorised", password)
			return c.TEMPL(http.StatusUnauthorized, loginPage(props))
		}

		props.success = true

		h.sessions.InitSession(props.username, c)
		return c.TEMPL(http.StatusUnauthorized, loginPage(props))
	}

	return c.TEMPL(http.StatusOK, loginPage(props))
}

func (h *AuthHandler) Logout(c *common.TemplContext) error {
	h.sessions.TerminateSession(c)
	return c.Redirect(http.StatusSeeOther, "/auth/login")
}
