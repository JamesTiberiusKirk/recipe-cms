package pages

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func InitMarkdownRenderer(app *echo.Group) {
	h := MarkdownRenderer{}

	app.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		s, _ := url.QueryUnescape(string(reqBody))
		logrus.Info("REQ BODY ", s)
		// logrus.Info("RES BODY ", string(resBody))
	}))

	app.POST("/markdown", common.UseTemplContext(h.Handle))
	app.GET("/markdown", common.UseTemplContext(h.Handle))
}

type MarkdownRenderer struct {
}

func (h *MarkdownRenderer) Handle(c *common.TemplContext) error {

	logrus.Info("markdown content-type :\n", c.Request().Header.Get("Content-Type"))

	buf := new(bytes.Buffer)
	b := c.Request().Body
	buf.ReadFrom(b)
	respBytes := buf.String()
	respString, _ := url.QueryUnescape((respBytes))
	logrus.Info(url.QueryUnescape(respString))
	return c.TEMPL(http.StatusOK, markdownContainer(respString))
}
