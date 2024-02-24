package auth

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/config"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site/components"
	"github.com/JamesTiberiusKirk/recipe-cms/site/session"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
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

	app.GET("/login", common.UseCustomContext(h.LoginPage))
	app.POST("/login", common.UseCustomContext(h.LoginPage))

	app.GET("/logout", common.UseCustomContext(h.Logout))

	app.GET("/login/shortcode", common.UseCustomContext(h.ShortCode))
	app.GET("/login/:code", common.UseCustomContext(h.ShortLogin))

	app.GET("/login/qr/:code", h.QrImage)

	app.GET("/login/sse", h.LoginSSE)
}

func (h *AuthHandler) LoginPage(c *common.Context) error {

	conf, ok := c.Get("cfg").(config.Config)
	if !ok {
		return fmt.Errorf("could not get config")
	}

	props := loginPageProps{c: c}

	if c.Request().Method == http.MethodPost {

		if conf.Debug {
			devLogin := c.QueryParam("dev_login")
			if devLogin == "true" {
				h.sessions.InitSession(props.username, c)
				hxto := components.HxTriggerOptions{
					ToastSuccess: "Logged in",
				}.ToJson()
				c.Response().Header().Set("HX-Trigger", hxto)

				source := c.QueryParam("source")
				if source != "" {
					c.Response().Header().Set("HX-Location", source)
				}

				return c.TEMPL(http.StatusOK, loginPage(props))
			}
		}
		props.username = c.FormValue("username")
		password := c.FormValue("password")

		if props.username == "" {
			props.errors.username = "You must enter a valid username"
		}

		if password == "" {
			props.errors.password = "You must enter a password"
		}

		if common.HasNonZeroField(props.errors) {
			logrus.Error("errors present ")
			return c.TEMPL(http.StatusUnauthorized, loginPage(props))
		}

		props.loginAttempted = true

		user, err := h.userRegistry.GetOneByUsername(props.username)
		if err != nil {
			logrus.Error("unable to get user from db ", err)
			return c.TEMPL(http.StatusUnauthorized, loginPage(props))
		}

		// NOTE: IKIK! i dont really care about security, at least not rn
		// I think I'll be overhauling it anyways
		if password != user.Password {
			logrus.Error("pass dont match")
			return c.TEMPL(http.StatusUnauthorized, loginPage(props))
		}

		props.success = true

		h.sessions.InitSession(props.username, c)
		hxto := components.HxTriggerOptions{
			ToastSuccess: "Logged in",
		}.ToJson()
		c.Response().Header().Set("HX-Trigger", hxto)

		source := c.QueryParam("source")
		if source != "" {
			c.Response().Header().Set("HX-Location", source)
		}
	}

	return c.TEMPL(http.StatusOK, loginPage(props))
}

func (h *AuthHandler) Logout(c *common.Context) error {
	h.sessions.TerminateSession(c)
	return c.Redirect(http.StatusSeeOther, "/auth/login")
}

// TODO: here we need to setup a SSE for this page so that it either
// refreshes on login or pushes a popup
// Will also need to figure out how to setup a channel so we can figure out when to send the server event
func (h *AuthHandler) ShortCode(c *common.Context) error {
	conf, ok := c.Get("cfg").(config.Config)
	if !ok {
		return fmt.Errorf("could not get config")
	}

	if h.sessions.IsAuthenticated(c, false) {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	short := h.sessions.InitShortCodeSess(c)
	return c.TEMPL(http.StatusOK, loginPageShortCode(loginPageShortCodeProps{code: short, host: conf.Host}))
}

func (h *AuthHandler) ShortLogin(c *common.Context) error {
	if !h.sessions.IsAuthenticated(c, true) {
		return c.Redirect(http.StatusSeeOther, "/auth/login?source="+c.Request().URL.String())
	}

	short := c.Param("code")
	err := h.sessions.AuthShortCodeSess(short, c)
	if err != nil {
		logrus.Errorf("error authenticating short code sess %s", err.Error())
		return err
	}

	source := c.QueryParam("source")
	if source != "" {
		c.Response().Header().Set("HX-Location", source)
	}
	return c.TEMPL(http.StatusOK, shortCodeTempPage(c))
}

func (h *AuthHandler) QrImage(c echo.Context) error {
	conf, ok := c.Get("cfg").(config.Config)
	if !ok {
		return fmt.Errorf("could not get config")
	}

	short := c.Param("code")
	var png []byte
	png, err := qrcode.Encode(conf.Host+"/auth/login/"+short, qrcode.Medium, 256)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "image/png", png)
}

func (h *AuthHandler) LoginSSE(c echo.Context) error {
	ch, err := h.sessions.GetShortCodeChan(c)
	if err != nil {
		logrus.Errorf("could not get code chan %s", err.Error())
		return err
	}

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Content-Type", "text/event-stream")

	for authEvent := range ch {
		if authEvent != "authenticated" {
			continue
		}

		h.sessions.IsAuthenticated(c, true)
		var buf bytes.Buffer
		authenticatedComponent().Render(c.Request().Context(), &buf)

		logrus.Print("sending data")
		// fmt.Fprintf(c.Response().Writer, "event: authEvent\ndata: authed \n\n")
		fmt.Fprintf(c.Response().Writer, "event: auth-event\ndata: %s \n\n", buf.String())
		c.Response().Writer.(http.Flusher).Flush()
		// time.Sleep(2 * time.Second)
	}

	return nil
}
