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
	QueryNameUpsertRecipe QueryName = "upsert_recipe"
	QueryNameInsertTags   QueryName = "insert_tags"
)

// GetQuery - returns query based on name
func (db *DB) GetQuery(name QueryName) (string, map[string]string, error) {
	sq, ok := db.queries[string(name)]
	if !ok {
		return "", nil, fmt.Errorf("schema \"%s\" not not found", name)
	}

	return sq.Query, sq.Tags, nil
}
