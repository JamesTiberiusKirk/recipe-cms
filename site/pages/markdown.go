package pages

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func InitMarkdownRenderer(app *echo.Group) {
	h := MarkdownRenderer{}

	app.POST("/markdown", common.UseCustomContext(h.Handle))
	app.GET("/markdown", common.UseCustomContext(h.Handle))
}

type MarkdownRenderer struct {
}

func (h *MarkdownRenderer) Handle(c *common.Context) error {
	// logrus.Info("markdown content-type :\n", c.Request().Header.Get("Content-Type"))

	buf := new(bytes.Buffer)
	b := c.Request().Body
	buf.ReadFrom(b)
	respBytes := buf.String()
	respString, _ := url.QueryUnescape((respBytes))
	logrus.Info(url.QueryUnescape(respString))
	return c.TEMPL(http.StatusOK, markdownContainer(respString))
}
