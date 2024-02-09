package auth

import (
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site/session"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	sessions     *session.Manager
	userRegistry registry.IUser
}

func InitAuthHandler(app *echo.Group, s *session.Manager, ur registry.IUser) {
	h := AuthHandler{
		sessions:     s,
		userRegistry: ur,
	}

	app.GET("/login", common.UseTemplContext(h.LoginPage))
	app.POST("/login", common.UseTemplContext(h.LoginPage))

	app.GET("/logout", common.UseTemplContext(h.Logout))
}

func (h *AuthHandler) LoginPage(c *common.TemplContext) error {
	props := loginPageProps{}

	if c.Request().Method == http.MethodPost {
		props.username = c.FormValue("username")
		password := c.FormValue("password")

		props.loginAttempted = true

		if props.username == "" {
			props.errors.username = "You must enter a valid username"
			return c.TEMPL(http.StatusUnauthorized, loginPage(props))
		}

		if password == "" {
			props.errors.password = "You must enter a password"
			return c.TEMPL(http.StatusUnauthorized, loginPage(props))
		}

		user, err := h.userRegistry.GetOneByUsername(props.username)
		if err != nil {
			return c.TEMPL(http.StatusUnauthorized, loginPage(props))
		}

		// NOTE: IKIK! i dont really care about security, at least not rn
		// I think I'll be overhauling it anyways
		if password != user.Password {
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
