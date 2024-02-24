package pages

import (
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/site/components"
)

func HandleNotFound(c *common.Context) error {
	return c.TEMPL(http.StatusNotFound, components.NotFound(c))
}
