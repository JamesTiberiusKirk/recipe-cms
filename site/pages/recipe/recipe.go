package recipe

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/config"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site/components"
	"github.com/JamesTiberiusKirk/recipe-cms/site/session"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
	app.POST("/:recipe_id", common.UseCustomContext(h.Page))
	app.DELETE("/:recipe_id", common.UseCustomContext(h.DeleteRecipe))

	app.GET("/ingredient", common.UseCustomContext(h.Ingredient))

	app.POST("/image", common.UseCustomContext(h.Image))
	app.DELETE("/image/:image", common.UseCustomContext(h.ImageDelete))
}

type RecipeRequestData struct {
	RecipeID string         `param:"recipe_id"`
	Recipe   *models.Recipe `json:"recipe,omitempty"`
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
		Edit:            (c.QueryParam("edit") == "true" || reqData.RecipeID == "new"),
		IsAuthenticated: h.session.IsAuthenticated(c, false),
	}

	if data.Edit && !data.IsAuthenticated {
		return c.Redirect(http.StatusSeeOther, "/auth/login?source="+url.QueryEscape(c.Request().URL.String()))
		// c.Response().Header().Set("HX-Redirect", "/auth/login?source="+url.QueryEscape(c.Request().URL.String()))
		// return c.NoContent(http.StatusUnauthorized)
	}

	status := http.StatusOK

	if reqData.RecipeID != "" && reqData.RecipeID != "new" {
		recipe, err := h.recipeRegistry.GetOneByID(reqData.RecipeID)
		if err != nil {
			return err
		}

		if recipe == nil {
			c.Logger().Info("recipe not found id: %s", reqData.RecipeID)
			return c.TEMPL(http.StatusNotFound, components.NotFound(c))
		}

		fmt.Printf("ingredients: %v\n", recipe.Ingredients)

		data.Recipe = *recipe
	}

	switch c.Request().Method {
	case http.MethodPost:
		if c.Request().Header.Get("Content-Type") != "application/json" {
			break
		}

		if reqData.Recipe == nil {
			break
		}

		reqData.Recipe.AuthorName = data.Recipe.AuthorName
		reqData.Recipe.ID = data.Recipe.ID
		data.Recipe = *reqData.Recipe

		if reqData.RecipeID == "new" {
			// NOTE: the only data that isnt coming back from the browser is author
			user, err := h.session.GetUser(c)
			if err != nil {
				return err
			}

			data.Recipe.AuthorName = user
			data.Recipe.ID = uuid.NewString()

		}

		upserted, _, err := h.recipeRegistry.Upsert(data.Recipe)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "error upserting")
		}

		if reqData.RecipeID == "new" {
			logrus.Print("EDIT REDIRECT")
			c.Response().Header().Set("HX-Redirect", "/recipe/"+upserted.ID+"?edit=true")
			status = http.StatusCreated
		}
		data.Recipe = upserted

		if !data.Edit {
			logrus.Print("NO EDIT REDIRECT")
			c.Response().Header().Set("HX-Redirect", "/recipe/"+upserted.ID)
		}
	case http.MethodGet:
		if reqData.RecipeID == "new" && c.QueryParam("edit") != "true" {
			return c.Redirect(http.StatusSeeOther, "/recipe/new?edit=true")
		}
	}

	// if c.QueryParam("edit") != "true" {
	// 	if len(data.Recipe.Ingredients) <= 0 {
	// 		data.Recipe.Ingredients = []models.Ingredient{{}}
	// 	}
	//
	// 	if len(data.Recipe.Seasonings) <= 0 {
	// 		data.Recipe.Seasonings = []models.Ingredient{{}}
	// 	}
	// }

	if reqData.RecipeID == "new" {
		data.Recipe.ID = "new"
	}

	return c.TEMPL(status, recipePage(data))
}

func (h *RecipeHandler) Ingredient(c *common.Context) error {
	t := c.QueryParam("type")
	if t == "" {
		t = "ingredient"
	}

	return c.TEMPL(http.StatusOK, components.Ingredient(components.IngredientProps{
		ID:             t,
		FormName:       []string{"recipe", fmt.Sprintf("%ss", t), ""},
		Ingredient:     models.Ingredient{},
		AvailableUnits: models.DefaultSystemUnits,
	}))
}

func (h *RecipeHandler) Image(c *common.Context) error {
	cfg, ok := c.Get("cfg").(config.Config)
	if !ok {
		return fmt.Errorf("could not get config from context")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["files"]

	recipeID := form.Value["recipe_id"][0]

	recipePath := fmt.Sprintf("%s/%s", cfg.Volume, recipeID)
	if _, err := os.Stat(recipePath); os.IsNotExist(err) {
		err := os.Mkdir(recipePath, 0755)
		if err != nil {
			logrus.Errorf("error making recipe directory %s", err.Error())
			return fmt.Errorf("error making recipe directory %w", err)
		}
	}
	if err != nil {
		logrus.Errorf("error getting recipe folder stats stats %s", err.Error())
		return fmt.Errorf("error getting recipe folder stats stats %w", err)
	}

	fileNames := []string{}

	for _, file := range files {
		src, err := file.Open()
		if err != nil {
			logrus.Errorf("error opening file: %s", err.Error())
			return err
		}
		defer src.Close()

		fileSplit := strings.Split(file.Filename, ".")
		fileName := fmt.Sprintf("%s/%s.%s", recipeID, uuid.New(), fileSplit[len(fileSplit)-1])
		filePath := fmt.Sprintf("%s/%s", cfg.Volume, fileName)

		dst, err := os.Create(filePath)
		if err != nil {
			logrus.Errorf("error creating file: %s", err.Error())
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			logrus.Errorf("error copying file: %s", err.Error())
			return err
		}

		fileNames = append(fileNames, "/images/"+fileName)
	}

	if c.QueryParam("type") == "add" {
		recipe, err := h.recipeRegistry.GetOneByID(recipeID)
		if err != nil {
			logrus.Errorf("error getting existing recipe %s", err.Error())
			return fmt.Errorf("error getting existing recipe %w", err)
		}

		fileNames = append(recipe.Images, fileNames...)
	}

	_, _, err = h.recipeRegistry.Upsert(models.Recipe{
		ID:     recipeID,
		Images: fileNames,
	})
	if err != nil {
		logrus.Infof("error upserting images: %s", err.Error())
		return fmt.Errorf("error upserting images: %w", err)
	}

	// TODO: need to figure out how to manage an id of new
	// TODO: delete pictures which are getting replaces

	return c.TEMPL(http.StatusOK, imageForm(imageFormProps{
		RecipeID: recipeID,
		Images:   fileNames,
	}))
}

func (h *RecipeHandler) ImageDelete(c *common.Context) error {
	cfg, ok := c.Get("cfg").(config.Config)
	if !ok {
		return fmt.Errorf("could not get config from context")
	}

	imageName := c.Param("image")

	imagePath := fmt.Sprintf("%s/%s", cfg.Volume, imageName)

	// TODO: figure out how to remove from db

	logrus.Infof("image path: %s", imagePath)

	_, err := os.Stat(imagePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist")
	}

	err = os.Remove(imagePath)
	if err != nil {
		return fmt.Errorf("error removing file: %w", err)
	}

	// TODO: send back the entire images form

	return nil
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
