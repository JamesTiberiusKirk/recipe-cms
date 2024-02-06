package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql"
	_ "github.com/lib/pq"
)

type DB struct {
	DB      *sqlx.DB
	queries goyesql.Queries
}

type Transaction struct {
	QueryName QueryName
	Args      any
}

func Connect(dbUrl string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	queries := goyesql.MustParseFile("./sql/queries.sql")

	return &DB{
		DB:      db,
		queries: queries,
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

// GetQuery - returns query based on name
func (db *DB) GetQuery(name QueryName) (string, map[string]string, error) {
	sq, ok := db.queries[string(name)]
	if !ok {
		return "", nil, fmt.Errorf("schema \"%s\" not not found", name)
	}

	return sq.Query, sq.Tags, nil
}
