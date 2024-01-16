package registry

import (
	"fmt"

	"github.com/JamesTiberiusKirk/recipe-cms/models"
	"github.com/sirupsen/logrus"
)

type Recipe struct {
	data []models.Recipe
}

func NewRecipeRegistry() *Recipe {
	return &Recipe{
		data: []models.Recipe{
			{
				ID:          "test-id",
				Name:        "chilly",
				Description: "# test description\n\n- bullet",
				Ingredients: []models.Ingredient{
					{
						Name:   "aaa",
						Amount: 111,
						Unit:   models.DefaultSystemUnits[1],
					},
					{
						Name:   "bbb",
						Amount: 222,
						Unit:   models.DefaultSystemUnits[1],
					},
					{
						Name:   "ccc",
						Amount: 333,
						Unit:   models.DefaultSystemUnits[1],
					},
					{
						Name:   "ddd",
						Amount: 444,
						Unit:   models.DefaultSystemUnits[1],
					},
				},
				Seasonings: []models.Ingredient{
					{
						Name:   "aaa",
						Amount: 111,
						Unit:   models.DefaultSystemUnits[1],
					},
					{
						Name:   "bbb",
						Amount: 222,
						Unit:   models.DefaultSystemUnits[1],
					},
					{
						Name:   "ccc",
						Amount: 333,
						Unit:   models.DefaultSystemUnits[1],
					},
					{
						Name:   "ddd",
						Amount: 444,
						Unit:   models.DefaultSystemUnits[1],
					},
				},
				Instructions:  "put one in another the bake it",
				LengthTotal:   "one billion years",
				LengthHandsOn: "like one second",
				Intro:         "the sukkiest of the recipes",
				Closing:       "never make this",
				Tags:          []string{"meat", "proteine", "keto", "veganunfriendly"},
				Version:       0,
				Author: models.User{
					Name: "Test User",
				},
			},
		},
	}
}

func (r *Recipe) GetAll() ([]models.Recipe, error) {
	return r.data, nil
}

func (r *Recipe) GetOneByID(id string) (*models.Recipe, error) {
	for _, recipe := range r.data {
		if id == recipe.ID {
			return &recipe, nil
		}
	}
	return nil, nil
}

func (r *Recipe) Upsert(upsert models.Recipe) (models.Recipe, bool, error) {
	logrus.Infof("upserting %v", upsert)
	if upsert.ID == "" {
		upsert.ID = fmt.Sprint("test-id", len(r.data))
		logrus.Infof("generating id of %s", upsert.ID)
		r.data = append(r.data, upsert)
		return upsert, false, nil
	}

	for i, recipe := range r.data {
		if upsert.ID == recipe.ID {
			r.data[i] = upsert
			return upsert, true, nil
		}
	}

	return models.Recipe{}, false, nil
}
