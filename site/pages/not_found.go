package pages

import (
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
)

func HandleNotFound(c *common.TemplContext) error {
	return c.TEMPL(http.StatusNotFound, notFound())
}
