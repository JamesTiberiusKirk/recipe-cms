package recipe

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/JamesTiberiusKirk/recipe-cms/common"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site/session"
	"github.com/labstack/echo/v4"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sirupsen/logrus"
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
	Tag   string `query:"tag"`
	Query string `query:"query"`
}

func (h *RecipesHandler) Page(c *common.Context) error {
	reqData := RecipesRequestData{}
	err := c.Bind(&reqData)
	if err != nil {
		return err
	}

	data := recipesPageData{c: c}

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

	// NOTE: bro this shit is actually deep fried
	// The complexity of this is probs retarded lol
	if reqData.Query != "" {
		c := make([]fmt.Stringer, len(data.recipes))
		for i, v := range data.recipes {
			c[i] = v
		}

		matches := fuzzy.RankFindStringer(reqData.Query, c)
		sort.Sort(matches)
		recipes := make([]models.Recipe, len(matches))
		for i, ri := range matches {
			r, ok := ri.Target.(models.Recipe)
			if !ok {
				logrus.Error("YEAH NOT OK")
			}

			recipes[i] = r
		}

		data.recipes = recipes
	}

	sess, ok := c.Get("session").(*session.Manager)
	if ok {
		data.isAuthenticated = sess.IsAuthenticated(c, false)
	}

	return c.TEMPL(http.StatusOK, recipesPage(data))
}
