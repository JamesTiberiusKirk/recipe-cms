package registry

import (
	"github.com/JamesTiberiusKirk/recipe-cms/db"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
)

type IRecipe interface {
	GetAll() ([]models.Recipe, error)
	GetOneByID(id string) (*models.Recipe, error)
	Upsert(upsert models.Recipe) (models.Recipe, bool, error)
}

type Recipe struct {
	dbc *db.DB
}

func NewRecipe(dbc *db.DB) *Recipe {
	return &Recipe{
		dbc: dbc,
	}
}

func (r *Recipe) GetAll() ([]models.Recipe, error) {
	panic("not implemented") // TODO: Implement
}

func (r *Recipe) GetOneByID(id string) (*models.Recipe, error) {
	panic("not implemented") // TODO: Implement
}

func (r *Recipe) Upsert(upsert models.Recipe) (models.Recipe, bool, error) {
	panic("not implemented") // TODO: Implement
}
