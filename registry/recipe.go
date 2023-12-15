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
				Description: "test",
				Ingredients: []models.Ingredient{
					{
						Name:   "aaa",
						Amount: 111,
					},
					{
						Name:   "bbb",
						Amount: 222,
					},
					{
						Name:   "ccc",
						Amount: 333,
					},
					{
						Name:   "ddd",
						Amount: 444,
					},
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

func (r *Recipe) Upsert(upsert models.Recipe) error {
	logrus.Infof("upserting %v", upsert)
	if upsert.ID == "" {
		upsert.ID = fmt.Sprint("test-id", len(r.data))
		logrus.Infof("generating id of %s", upsert.ID)
		r.data = append(r.data, upsert)
		return nil
	}

	for i, recipe := range r.data {
		if upsert.ID == recipe.ID {
			r.data[i] = upsert
			return nil
		}
	}

	return nil
}
