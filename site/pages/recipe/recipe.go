package recipe

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site/components"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type RecipeHandler struct {
	recipeRegistry *registry.Recipe
}

func InitRecipeHandler(app *echo.Group, rr *registry.Recipe) {
	h := &RecipeHandler{
		recipeRegistry: rr,
	}

	app.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		s, _ := url.QueryUnescape(string(reqBody))
		logrus.Info("REQ BODY ", s)
	}))

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

	logrus.Infof("recipeID: %s\n", reqData.RecipeID)
	logrus.Infof("recipe: %v\n", reqData.Recipe)

	data := recipePageData{
		Units: models.DefaultSystemUnits,
		Edit:  (c.QueryParam("edit") == "true"),
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
		logrus.Infof("content type: %s", c.Request().Header.Get("Content-Type"))

		if c.Request().Header.Get("Content-Type") != "application/json" {
			break
		}

		// TODO: perform db update
		if reqData.Recipe != nil {
			data.Recipe = *reqData.Recipe
			if reqData.RecipeID != "new" {
				data.Recipe.ID = reqData.RecipeID
			}

			logrus.Info("upserting")
			upserted, wasUpserted, err := h.recipeRegistry.Upsert(data.Recipe)
			if err != nil {
				return echo.NewHTTPError(500, "error upserting")
			}

			logrus.Info("UPSERTED", wasUpserted, upserted)
			logrus.Info("ID", reqData.RecipeID)

			if reqData.RecipeID == "new" {
				c.Response().Header().Set("HX-Redirect", fmt.Sprintf("/recipe/%s?edit=true", upserted.ID))
				status = http.StatusCreated
			}
			data.Recipe = upserted
		}

	}

	if len(data.Recipe.Ingredients) <= 0 {
		logrus.Info("ingredients empty so assing an empty")
		data.Recipe.Ingredients = []models.Ingredient{
			{},
		}
	}

	logrus.Info("sending back endit", data.Edit)
	logrus.Info("req data", reqData.Edit)
	logrus.Info("c.QueryParam", c.QueryParam("edit"))

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
