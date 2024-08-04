package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// squirrel helper
type SqlBuilder interface {
	ToSql() (string, []interface{}, error)
}

type DB struct {
	DB *sqlx.DB
}

func Connect(dbUrl string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	return &DB{
		DB: db,
	}, nil
}

type QueryName string

const (
	GetAllRecipesByTagName  QueryName = "get_all_recipes_by_tag_name"
	GetAllRecipes           QueryName = "get_all_recipes"
	GetRecipeByID           QueryName = "get_recipes_by_id"
	UpsertRecipe            QueryName = "upsert_recipe"
	DeleteAllTagsByRecipeID QueryName = "delete_all_tags_by_recipe_id"
	InsertTag               QueryName = "insert_tag"
	DeleteIngredient        QueryName = "delete_ingredient"
	UpsertIngredients       QueryName = "upsert_ingredients"
)
