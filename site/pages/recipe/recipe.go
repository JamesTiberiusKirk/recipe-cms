package recipe

import (
	"bytes"
	"net/http"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type RecipeHandler struct {
	recipeRegistry *registry.Recipe
}

func InitEditRecipeHandler(app *echo.Group, rr *registry.Recipe) {
	h := &RecipeHandler{
		recipeRegistry: rr,
	}

	app.GET("/:recipe_id", common.UseTemplContext(h.Page))
	app.POST("/:recipe_id", common.UseTemplContext(h.Page))

	app.GET("/ingredient", common.UseTemplContext(h.Ingredient))
}

type RecipeRequestData struct {
	RecipeID string         `param:"recipe_id"`
	Recipe   *models.Recipe `json:"recipe,omitempty"`
}

func (h *RecipeHandler) Page(c *common.TemplContext) error {
	reqData := RecipeRequestData{}
	echo.QueryParamsBinder(c)
	err := c.Bind(&reqData)
	if err != nil {
		return err
	}

	logrus.Infof("recipeID: %s\n", reqData.RecipeID)
	logrus.Infof("recipe: %v\n", reqData.Recipe)

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request().Body)

	logrus.Infof("body: %v\n", buf.String())

	data := recipePageData{
		Units: models.DefaultSystemUnits,
	}

	if reqData.RecipeID != "" && reqData.RecipeID != "new" {
		recipe, err := h.recipeRegistry.GetOneByID(reqData.RecipeID)
		if err != nil {
			return err
		}

		if recipe == nil {
			c.Logger().Info("recipe not found id: %s", reqData.RecipeID)
			return c.Redirect(http.StatusTemporaryRedirect, "/404")
		}

		data.Recipe = *recipe
	}

	if c.Request().Method == http.MethodPost {
		// TODO: perform db update
		if reqData.Recipe != nil {
			data.Recipe = *reqData.Recipe
			h.recipeRegistry.Upsert(*reqData.Recipe)
		}
	}

	if len(data.Recipe.Ingredients) <= 0 {
		data.Recipe.Ingredients = []models.Ingredient{
			{},
		}
	}

	return c.TEMPL(http.StatusOK, recipePage(data))
}

func (h *RecipeHandler) Ingredient(c *common.TemplContext) error {
	return c.TEMPL(http.StatusOK, ingredient(models.Ingredient{}, models.DefaultSystemUnits))
}
