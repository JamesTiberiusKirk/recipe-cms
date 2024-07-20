package edit

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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

type EditRecipeHandler struct {
	recipeRegistry registry.IRecipe
	session        *session.Manager
}

func InitEditRecipeHandler(app *echo.Group, rr registry.IRecipe, s *session.Manager) {
	h := &EditRecipeHandler{
		recipeRegistry: rr,
		session:        s,
	}

	// app.Use(common.UseBodyLogger())
	app.Use(common.AuthMiddleware(s))

	app.GET("", common.UseCustomContext(h.Page))
	app.POST("", common.UseCustomContext(h.Page))

	app.POST("/image", common.UseCustomContext(h.Image))
	app.DELETE("/image/:image", common.UseCustomContext(h.ImageDelete))
}

type RecipeRequestData struct {
	RecipeID string         `param:"recipe_id"`
	Recipe   *models.Recipe `json:"recipe,omitempty"`
}

func (h *EditRecipeHandler) Page(c *common.Context) error {
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

	if reqData.RecipeID == "" {
		return c.TEMPL(http.StatusNotFound, components.NotFound(c))
	}

	if reqData.RecipeID != "new" {
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

	switch c.Request().Method {
	case http.MethodPost:

		action := c.QueryParam("action")
		if action != "" {
			err := h.manageAction(c, &reqData, &data, action)
			if err != nil && !errors.Is(err, errTerminated) {
				c.Logger().Error("error managing action: %s", err.Error())
				return echo.NewHTTPError(http.StatusInternalServerError, "error managing action")
			}
			return nil
		}

		err := h.manageUpsert(c, &reqData, &data)
		if err != nil && !errors.Is(err, errTerminated) {
			c.Logger().Error("error managing upsert: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "error upserting")
		}

		if c.QueryParam("done") == "true" {
			return c.Redirect(http.StatusSeeOther, "/recipe/"+data.Recipe.ID)
		}

	case http.MethodGet:

	}

	return c.TEMPL(status, recipePage(data))
}

var (
	errTerminated = errors.New("http call terminated")
)

const (
	ActionAddIngredient = "add_ingredient"
	ActionAddSeasoning  = "add_seasoning"
)

func (h *EditRecipeHandler) manageAction(c *common.Context, reqData *RecipeRequestData, data *recipePageData, action string) error {
	switch action {
	case ActionAddIngredient:
		data.Recipe.Ingredients = append(reqData.Recipe.Ingredients, models.Ingredient{})
	case ActionAddSeasoning:
		data.Recipe.Seasonings = append(reqData.Recipe.Seasonings, models.Ingredient{})
	default:
		_ = c.String(http.StatusBadRequest, "invalid action")
		return errTerminated
	}

	err := c.TEMPL(http.StatusOK, recipePage(*data))
	if err != nil {
		return fmt.Errorf("error rendering page: %w", err)
	}
	return errTerminated
}

func (h *EditRecipeHandler) manageUpsert(c *common.Context, reqData *RecipeRequestData, data *recipePageData) error {
	if c.Request().Header.Get("Content-Type") != "application/json" {
		return nil
	}

	if reqData.Recipe == nil {
		return nil
	}

	reqData.Recipe.AuthorName = data.Recipe.AuthorName
	reqData.Recipe.ID = data.Recipe.ID
	data.Recipe = *reqData.Recipe

	if reqData.RecipeID == "new" {
		// NOTE: the only data that isnt coming back from the browser is author
		user, err := h.session.GetUser(c)
		if err != nil {
			return fmt.Errorf("error getting user session: %s", err)
		}

		data.Recipe.AuthorName = user
		data.Recipe.ID = uuid.NewString()

	}

	upserted, _, err := h.recipeRegistry.Upsert(data.Recipe)
	if err != nil {
		return fmt.Errorf("error upserting: %w", err)
	}

	if reqData.RecipeID == "new" {
		c.Logger().Print("EDIT REDIRECT")
		c.Response().Header().Set("HX-Redirect", "/recipe/"+upserted.ID+"/edit")
		_ = c.TEMPL(http.StatusCreated, recipePage(*data))

		return errTerminated
	}
	data.Recipe = upserted

	return nil
}

func (h *EditRecipeHandler) Image(c *common.Context) error {
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

func (h *EditRecipeHandler) ImageDelete(c *common.Context) error {
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

func (h *EditRecipeHandler) DeleteRecipe(c *common.Context) error {
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
