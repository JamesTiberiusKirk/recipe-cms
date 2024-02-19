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

	app.GET("", common.UseTemplContext(h.Handle))
}

func (h *IndexHandler) Handle(e *common.TemplContext) error {
	return e.TEMPL(http.StatusOK, indexPage())
}
