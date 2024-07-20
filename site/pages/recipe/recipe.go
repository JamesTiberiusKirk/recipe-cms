package recipe

import (
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site/components"
	"github.com/JamesTiberiusKirk/recipe-cms/site/pages/recipe/edit"
	"github.com/JamesTiberiusKirk/recipe-cms/site/session"
	"github.com/labstack/echo/v4"
)

type RecipeHandler struct {
	recipeRegistry registry.IRecipe
	session        *session.Manager
}

func InitRecipeHandler(app *echo.Group, rr registry.IRecipe, s *session.Manager) {
	h := &RecipeHandler{
		recipeRegistry: rr,
		session:        s,
	}

	// app.Use(common.UseBodyLogger())

	app.GET("/:recipe_id", common.UseCustomContext(h.Page))
	app.GET("/:recipe_id/mfp", common.UseCustomContext(h.MFPPage))
	app.DELETE("/:recipe_id", common.UseCustomContext(h.DeleteRecipe))

	// /edit
	edit.InitEditRecipeHandler(app.Group("/:recipe_id/edit"), rr, s)
}

type RecipeRequestData struct {
	RecipeID string `param:"recipe_id"`
}

func (h *RecipeHandler) MFPPage(c *common.Context) error {
	reqData := RecipeRequestData{}
	err := c.Bind(&reqData)
	if err != nil {
		return err
	}

	if reqData.RecipeID == "" && reqData.RecipeID != "new" {
		return c.TEMPL(http.StatusNotFound, components.NotFound(c))
	}

	recipe, err := h.recipeRegistry.GetOneByID(reqData.RecipeID)
	if err != nil {
		return err
	}

	if recipe == nil {
		c.Logger().Info("recipe not found id: %s", reqData.RecipeID)
		return c.TEMPL(http.StatusNotFound, components.NotFound(c))
	}

	return c.TEMPL(http.StatusOK, recipeMFPView(recipeMFPViewProps{c: c, Recipe: *recipe}))
}

func (h *RecipeHandler) Page(c *common.Context) error {
	reqData := RecipeRequestData{}
	err := c.Bind(&reqData)
	if err != nil {
		return err
	}

	data := recipePageData{
		c:               c,
		Units:           models.DefaultSystemUnits,
		IsAuthenticated: h.session.IsAuthenticated(c, false),
	}

	status := http.StatusOK

	if reqData.RecipeID != "" {
		recipe, err := h.recipeRegistry.GetOneByID(reqData.RecipeID)
		if err != nil {
			return err
		}

		if recipe == nil {
			c.Logger().Info("recipe not found id: %s", reqData.RecipeID)
			return c.TEMPL(http.StatusNotFound, components.NotFound(c))
		}

		data.Recipe = *recipe
	}

	return c.TEMPL(status, recipePage(data))
}

func (h *RecipeHandler) DeleteRecipe(c *common.Context) error {
	if !h.session.IsAuthenticated(c, false) {
		return c.NoContent(http.StatusUnauthorized)
	}

	recipeID := c.Param("recipe_id")

	err := h.recipeRegistry.DeleteOne(recipeID)
	if err != nil {
		c.Logger().Error("Error deleting recipe with id %s, err: %s", recipeID, err.Error())
		return c.NoContent(http.StatusNoContent)
	}

	return c.NoContent(http.StatusOK)
}
