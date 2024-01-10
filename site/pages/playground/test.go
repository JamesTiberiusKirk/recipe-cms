package playground

import (
	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/labstack/echo/v4"
)

type TestRoute struct{}

func InitTestRoute(app *echo.Group) {
	h := TestRoute{}

	app.GET("/testpage", common.UseTemplContext(h.HandleTestPage))
}

func (h *TestRoute) HandleTestPage(c *common.TemplContext) error {
	return c.TEMPL(200, testPage())
}
