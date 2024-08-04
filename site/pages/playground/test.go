package playground

import (
	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/config"
	"github.com/labstack/echo/v4"
)

type TestRoute struct{}

func InitTestRoute(conf config.Config, app *echo.Group) {
	if !conf.Debug {
		return
	}

	h := TestRoute{}
	app.GET("/pg/testpage", common.UseCustomContext(h.HandleTestPage))
}

func (h *TestRoute) HandleTestPage(c *common.Context) error {
	return c.TEMPL(200, testPage(c))
}
