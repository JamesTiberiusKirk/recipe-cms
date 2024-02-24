package pages

import (
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct {
}

func InitIndexHandler(app *echo.Group) {
	h := &IndexHandler{}

	app.GET("", common.UseCustomContext(h.Handle))
}

func (h *IndexHandler) Handle(c *common.Context) error {
	return c.TEMPL(http.StatusOK, indexPage(c))
}
