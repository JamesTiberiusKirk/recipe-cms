package recipe

import (
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/labstack/echo/v4"
)

type RecipesHandler struct {
	recipeRegistry registry.IRecipe
}

func InitRecipesHandler(app *echo.Group, rr registry.IRecipe) {
	h := &RecipesHandler{
		recipeRegistry: rr,
	}

	app.GET("", common.UseTemplContext(h.Page))
}

func (h *RecipesHandler) Page(c *common.TemplContext) error {
	recipes, err := h.recipeRegistry.GetAll()
	if err != nil {
		return err
	}

	data := recipesPageData{
		recipes: recipes,
	}

	return c.TEMPL(http.StatusOK, recipesPage(data))
}
