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

	app.GET("", common.UseCustomContext(h.Page))
}

type RecipesRequestData struct {
	Tag string `query:"tag"`
}

func (h *RecipesHandler) Page(c *common.Context) error {

	reqData := RecipesRequestData{}
	// echo.QueryParamsBinder(c)
	err := c.Bind(&reqData)
	if err != nil {
		return err
	}

	data := recipesPageData{}

	if reqData.Tag != "" {
		recipes, err := h.recipeRegistry.GetAllByTagName(reqData.Tag)
		if err != nil {
			return err
		}
		data.recipes = recipes

	} else {
		recipes, err := h.recipeRegistry.GetAll()
		if err != nil {
			return err
		}
		data.recipes = recipes
	}

	return c.TEMPL(http.StatusOK, recipesPage(data))
}
