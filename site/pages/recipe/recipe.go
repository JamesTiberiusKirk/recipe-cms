package recipe

import (
	"fmt"
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site/components"
	"github.com/labstack/echo/v4"
)

type RecipeHandler struct {
	recipeRegistry *registry.Recipe
}

func InitRecipeHandler(app *echo.Group, rr *registry.Recipe) {
	h := &RecipeHandler{
		recipeRegistry: rr,
	}

	app.Use(common.UseBodyLogger())

	app.GET("/:recipe_id", common.UseTemplContext(h.Page))
	app.POST("/:recipe_id", common.UseTemplContext(h.Page))

	app.GET("/ingredient", common.UseTemplContext(h.Ingredient))
}

type RecipeRequestData struct {
	RecipeID string         `param:"recipe_id"`
	Edit     bool           `query:"edit"`
	Recipe   *models.Recipe `json:"recipe,omitempty"`
}

func (h *RecipeHandler) Page(c *common.TemplContext) error {
	reqData := RecipeRequestData{}
	echo.QueryParamsBinder(c)
	err := c.Bind(&reqData)
	if err != nil {
		return err
	}

	data := recipePageData{
		Units: models.DefaultSystemUnits,
		Edit:  (c.QueryParam("edit") == "true" || reqData.RecipeID == "new"),
	}

	status := http.StatusOK

	if reqData.RecipeID != "" && reqData.RecipeID != "new" {
		recipe, err := h.recipeRegistry.GetOneByID(reqData.RecipeID)
		if err != nil {
			return err
		}

		if recipe == nil {
			c.Logger().Info("recipe not found id: %s", reqData.RecipeID)
			return c.TEMPL(http.StatusNotFound, components.NotFound())
		}

		data.Recipe = *recipe
	}

	switch c.Request().Method {
	case http.MethodPost:
		if c.Request().Header.Get("Content-Type") != "application/json" {
			break
		}

		if reqData.Recipe != nil {
			data.Recipe = *reqData.Recipe
			if reqData.RecipeID != "new" {
				data.Recipe.ID = reqData.RecipeID
			}

			upserted, _, err := h.recipeRegistry.Upsert(data.Recipe)
			if err != nil {
				return echo.NewHTTPError(500, "error upserting")
			}

			if reqData.RecipeID == "new" {
				c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/recipe/%s?edit=true", upserted.ID))
				status = http.StatusCreated
			}
			data.Recipe = upserted
		}
	}

	if len(data.Recipe.Ingredients) <= 0 {
		data.Recipe.Ingredients = []models.Ingredient{{}}
	}

	if len(data.Recipe.Seasonings) <= 0 {
		data.Recipe.Seasonings = []models.Ingredient{{}}
	}

	if reqData.RecipeID == "new" {
		data.Recipe.ID = "new"
	}

	return c.TEMPL(status, recipePage(data))
}

func (h *RecipeHandler) Ingredient(c *common.TemplContext) error {
	t := c.QueryParam("type")
	if t == "" {
		t = "ingredient"
	}

	return c.TEMPL(http.StatusOK, components.Ingredient(components.IngredientProps{
		ID:             t,
		FormName:       []string{"recipe", fmt.Sprintf("%ss", t), ""},
		Ingredient:     models.Ingredient{},
		AvailableUnits: models.DefaultSystemUnits}))
}
