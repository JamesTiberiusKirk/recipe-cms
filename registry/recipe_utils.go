package registry

import (
	"fmt"
	"strings"

	"github.com/JamesTiberiusKirk/recipe-cms/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/rustedturnip/goscanql"
)

const (
	recipeTableName     = "recipe"
	ingredientTableName = "ingredient"
	unitTableName       = "unit"
	tagTableName        = "tag"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	recipeTableNameCols = []string{
		"id",
		"recipe_name",
		"intro",
		"description",
		"instructions",
		"length_total",
		"length_hands_on",
		"closing",
		"recipe_version",
		"author_name",
		"images",
	}

	recipeTableColsAS = []string{
		"r.id                AS id",
		"r.recipe_name       AS recipename",
		"r.intro             AS intro",
		"r.description       AS description",
		"r.instructions      AS instructions",
		"r.length_total      AS lengthtotal",
		"r.length_hands_on   AS lengthhandson",
		"r.closing           AS closing",
		"r.recipe_version    AS recipeversion",
		"r.author_name       AS authorname",
		"r.images            AS images",
		"t.tag_name          AS tags",
	}

	ingredientTableCols = []string{
		"recipe_id",
		"array_index",
		"field",
		"ingredient_name",
		"amount",
		"unit_name",
	}

	ingredientTableColsAs = []string{
		"i.recipe_id",
		"i.array_index       AS array_index",
		"i.ingredient_name   AS name",
		"i.amount            AS amount",
		"iu.unit_name        AS unit_name",
		"iu.display_name     AS unit_displayname",
	}

	tagTableCols = []string{
		"recipe_id",
		"tag_name",
	}
)

func makeSelectIngredientSelectSatement(recipeID, field string) sq.SelectBuilder {
	return psql.Select(ingredientTableColsAs...).
		From(ingredientTableName + " i").
		LeftJoin(unitTableName + " as iu ON iu.unit_name = i.unit_name").
		Where(sq.Eq{"i.field": recipeID}).
		Where(sq.Eq{"i.recipe_id": field})
}

func (r *Recipe) getSeasoningsByRecipeID(recipeID string) ([]models.Ingredient, error) {
	rows, err := makeSelectIngredientSelectSatement(recipeID, "SEASONING").
		RunWith(r.dbc.DB).Query()
	if err != nil {
		return nil, fmt.Errorf("error quering db: %w", err)
	}

	seasonings, err := goscanql.RowsToStructs[models.Ingredient](rows)
	if err != nil {
		return nil, fmt.Errorf("error mapping row to structs: %w", err)
	}

	return seasonings, nil
}

func (r *Recipe) getIngredientsByRecipeID(recipeID string) ([]models.Ingredient, error) {
	rows, err := makeSelectIngredientSelectSatement(recipeID, "INGREDIENT").
		RunWith(r.dbc.DB).Query()
	if err != nil {
		return nil, fmt.Errorf("error quering db: %w", err)
	}

	ingredients, err := goscanql.RowsToStructs[models.Ingredient](rows)
	if err != nil {
		return nil, fmt.Errorf("error mapping row to structs: %w", err)
	}

	return ingredients, nil
}

func (r *Recipe) getRecipes(extraWhereParams ...sq.Eq) ([]models.Recipe, error) {
	selectRecipe := psql.Select(recipeTableColsAS...).
		From(recipeTableName + " r").
		LeftJoin(tagTableName + " as t ON t.recipe_id = r.id")

	for _, w := range extraWhereParams {
		selectRecipe = selectRecipe.Where(w)
	}

	rows, err := selectRecipe.RunWith(r.dbc.DB).Query()
	if err != nil {
		return nil, fmt.Errorf("error quering db: %w", err)
	}

	recipes, err := goscanql.RowsToStructs[models.Recipe](rows)
	if err != nil {
		return nil, fmt.Errorf("error mapping row to structs: %w", err)
	}

	for i := range recipes {
		if len(recipes[i].Images) < 1 {
			continue
		}

		trimmed := strings.Trim(strings.Trim(recipes[i].Images[0], "{"), "}")
		recipes[i].Images = strings.Split(trimmed, ",")
	}

	return recipes, nil
}
