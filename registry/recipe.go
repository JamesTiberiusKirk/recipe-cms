package registry

import (
	"fmt"
	"strings"

	"github.com/JamesTiberiusKirk/recipe-cms/db"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
	"github.com/lib/pq"
	"github.com/rustedturnip/goscanql"
	"github.com/sirupsen/logrus"
)

type IRecipe interface {
	GetAll() ([]models.Recipe, error)
	GetAllByTagName(string) ([]models.Recipe, error)
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
	query, _, err := r.dbc.GetQuery(db.GetAllRecipes)
	if err != nil {
		return nil, fmt.Errorf("error getting qeury: %w", err)
	}

	rows, err := r.dbc.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error quering db: %w", err)
	}

	recipes, err := goscanql.RowsToStructs[models.Recipe](rows)
	if err != nil {
		return nil, fmt.Errorf("error mapping row to structs: %w", err)
	}

	// TODO: fix this, should be done in sql
	for i := range recipes {
		if len(recipes[i].Images) < 1 {
			continue
		}

		trimmed := strings.Trim(strings.Trim(recipes[i].Images[0], "{"), "}")
		recipes[i].Images = strings.Split(trimmed, ",")
	}

	return recipes, nil
}

func (r *Recipe) GetAllByTagName(tag string) ([]models.Recipe, error) {
	query, _, err := r.dbc.GetQuery(db.GetAllRecipesByTagName)
	if err != nil {
		return nil, fmt.Errorf("error getting qeury: %w", err)
	}

	rows, err := r.dbc.DB.Query(query, tag)
	if err != nil {
		return nil, fmt.Errorf("error quering db: %w", err)
	}

	recipes, err := goscanql.RowsToStructs[models.Recipe](rows)
	if err != nil {
		return nil, fmt.Errorf("error mapping row to structs: %w", err)
	}

	// TODO: fix this, should be done in sql
	for i := range recipes {
		if len(recipes[i].Images) < 1 {
			continue
		}

		trimmed := strings.Trim(strings.Trim(recipes[i].Images[0], "{"), "}")
		recipes[i].Images = strings.Split(trimmed, ",")
	}

	return recipes, nil
}

func (r *Recipe) GetOneByID(id string) (*models.Recipe, error) {
	query, _, err := r.dbc.GetQuery(db.GetRecipeByID)
	if err != nil {
		return nil, fmt.Errorf("error getting qeury: %w", err)
	}

	rows, err := r.dbc.DB.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error quering db: %w", err)
	}

	recipe, err := goscanql.RowsToStruct[*models.Recipe](rows)
	if err != nil {
		return nil, fmt.Errorf("error mapping row to struct: %w", err)
	}

	// TODO: fix this, should be done in sql
	if len(recipe.Images) > 0 {
		trimmed := strings.Trim(strings.Trim(recipe.Images[0], "{"), "}")
		recipe.Images = strings.Split(trimmed, ",")
	}

	return recipe, nil
}

func (r *Recipe) Upsert(upsert models.Recipe) (models.Recipe, bool, error) {
	tagsInsertArgs := []map[string]any{}
	for _, t := range upsert.Tags {
		tagsInsertArgs = append(tagsInsertArgs, map[string]any{
			"recipeid": upsert.ID,
			"tagname":  t,
		})
	}

	ingredientInsertArgs := []map[string]any{}
	for i, ing := range upsert.Ingredients {
		ingredientInsertArgs = append(ingredientInsertArgs, map[string]any{
			"recipeid":       upsert.ID,
			"arrayindex":     i,
			"field":          "INGREDIENT",
			"ingredientname": ing.Name,
			"amount":         ing.Amount,
			"unitname":       ing.Unit.Name,
		})
	}

	seasoningInsertArgs := []map[string]any{}
	for i, ing := range upsert.Seasonings {
		seasoningInsertArgs = append(seasoningInsertArgs, map[string]any{
			"recipeid":       upsert.ID,
			"arrayindex":     i,
			"field":          "SEASONING",
			"ingredientname": ing.Name,
			"amount":         ing.Amount,
			"unitname":       ing.Unit.Name,
		})
	}

	transactions := []db.Transaction{
		{
			QueryName: db.UpsertRecipe,
			Args: map[string]any{
				"id":            upsert.ID,
				"recipename":    upsert.Name,
				"intro":         upsert.Intro,
				"description":   upsert.Description,
				"instructions":  upsert.Instructions,
				"lengthtotal":   upsert.LengthTotal,
				"lengthhandson": upsert.LengthHandsOn,
				"closing":       upsert.Closing,
				"recipeversion": upsert.RecipeVersion,
				"authorname":    upsert.AuthorName,
				"images":        pq.Array(upsert.Images),
			},
		},
		{
			QueryName: db.DeleteAllTagsByRecipeID,
			Args: map[string]any{
				"recipeid": upsert.ID,
			},
		},
		{
			QueryName: db.InsertTag,
			Args:      tagsInsertArgs,
		},
		{
			QueryName: db.DeleteIngredient,
			Args: map[string]any{
				"recipeid": upsert.ID,
				"field":    "INGREDIENT",
			},
		},
		{
			QueryName: db.UpsertIngredients,
			Args:      ingredientInsertArgs,
		},
		{
			QueryName: db.DeleteIngredient,
			Args: map[string]any{
				"recipeid": upsert.ID,
				"field":    "SEASONING",
			},
		},
		{
			QueryName: db.UpsertIngredients,
			Args:      seasoningInsertArgs,
		},
	}

	trx := r.dbc.DB.MustBegin()
	defer trx.Rollback()

	for _, t := range transactions {
		q, _, err := r.dbc.GetQuery(t.QueryName)
		if err != nil {
			logrus.Errorf("error getting query %s: %s", t.QueryName, err.Error())
			return models.Recipe{}, false, fmt.Errorf("error getting query: %w", err)
		}

		fmt.Println(q)

		switch args := t.Args.(type) {
		case []map[string]any:
			stmt, err := trx.PrepareNamed(q)
			if err != nil {
				logrus.Errorf("error prepairing transaction %s: %s", t.QueryName, err.Error())
				return models.Recipe{}, false, fmt.Errorf("error prepairing transaction: %w", err)
			}
			defer stmt.Close()

			for _, a := range args {
				_, err := stmt.Exec(a)
				if err != nil {
					logrus.Errorf("error creating prepaired transaction %s: %s", t.QueryName, err.Error())
					return models.Recipe{}, false, fmt.Errorf("error creating transaction: %w", err)
				}
			}
		default:
			_, err = trx.NamedExec(q, args)
			if err != nil {
				logrus.Errorf("error creating transaction %s: %s", t.QueryName, err.Error())
				return models.Recipe{}, false, fmt.Errorf("error creating transaction: %w", err)
			}
		}
	}

	err := trx.Commit()
	if err != nil {
		logrus.Errorf("error committing transaction: %s", err.Error())
		return models.Recipe{}, false, fmt.Errorf("error committing transaction: %w", err)
	}

	return upsert, false, nil
}
