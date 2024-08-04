package registry

import (
	"fmt"

	"github.com/JamesTiberiusKirk/recipe-cms/db"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type IRecipe interface {
	GetAll() ([]models.Recipe, error)
	GetAllByTagName(string) ([]models.Recipe, error)
	GetOneByID(id string) (*models.Recipe, error)
	Upsert(upsert models.Recipe) (models.Recipe, bool, error)
	DeleteOne(id string) error
}

type Recipe struct {
	dbc *db.DB
}

func NewRecipe(dbc *db.DB) *Recipe {
	return &Recipe{
		dbc: dbc,
	}
}

type ingredientChanResponse struct {
	Ingredients []models.Ingredient
	Err         error
	Index       int
}

func (r *Recipe) GetAll() ([]models.Recipe, error) {
	recipes, err := r.getRecipes()
	if err != nil {
		logrus.Errorf("error getting all recipes: %s", err.Error())
		return nil, fmt.Errorf("error getting all recipes: %w", err)
	}

	errChan := make(chan error, len(recipes)*2)
	defer close(errChan)

	for i, recipe := range recipes {
		go func(i int, recipeID string) {
			ingredients, err := r.getIngredientsByRecipeID(recipeID)
			errChan <- err
			recipes[i].Ingredients = ingredients
		}(i, recipe.ID)
		go func(i int, recipeID string) {
			seasonings, err := r.getSeasoningsByRecipeID(recipeID)
			errChan <- err
			recipes[i].Seasonings = seasonings
		}(i, recipe.ID)
	}

	for i := 0; i < len(recipes)*2; i++ {
		cErr := <-errChan

		if cErr != nil {
			err = cErr
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error getting all recipes: %w", err)
	}

	return recipes, nil
}

func (r *Recipe) GetAllByTagName(tag string) ([]models.Recipe, error) {
	recipes, err := r.getRecipes(sq.Eq{"t.tag_name": tag})
	if err != nil {
		logrus.Errorf("error getting all recipes: %s", err.Error())
		return nil, fmt.Errorf("error getting all recipes: %w", err)
	}

	if len(recipes) < 1 {
		return nil, nil
	}

	errChan := make(chan error, len(recipes)*2)
	defer close(errChan)

	for i, recipe := range recipes {
		go func(i int, recipeID string) {
			ingredients, err := r.getIngredientsByRecipeID(recipeID)
			errChan <- err
			recipes[i].Ingredients = ingredients
		}(i, recipe.ID)
		go func(i int, recipeID string) {
			seasonings, err := r.getSeasoningsByRecipeID(recipeID)
			errChan <- err
			recipes[i].Seasonings = seasonings
		}(i, recipe.ID)
	}

	for i := 0; i < len(recipes)*2; i++ {
		cErr := <-errChan

		if cErr != nil {
			err = cErr
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error getting all recipes: %w", err)
	}

	return recipes, nil
}

func (r *Recipe) GetOneByID(id string) (*models.Recipe, error) {
	recipes, err := r.getRecipes(sq.Eq{"r.id": id})
	if err != nil {
		logrus.Errorf("error getting all recipes: %s", err.Error())
		return nil, fmt.Errorf("error getting all recipes: %w", err)
	}

	if len(recipes) > 1 {
		return nil, fmt.Errorf("more than one recipe found with id %s", id)
	}

	if len(recipes) < 1 {
		return nil, nil
	}

	errChan := make(chan error, 2)
	defer close(errChan)

	recipe := recipes[0]

	go func(recipeID string) {
		ingredients, err := r.getIngredientsByRecipeID(recipeID)
		errChan <- err
		recipe.Ingredients = ingredients
	}(recipe.ID)
	go func(recipeID string) {
		seasonings, err := r.getSeasoningsByRecipeID(recipeID)
		errChan <- err
		recipe.Seasonings = seasonings
	}(recipe.ID)

	for i := 0; i < 2; i++ {
		cErr := <-errChan

		if cErr != nil {
			err = cErr
		}
	}
	if err != nil {
		return nil, fmt.Errorf("error getting all recipes: %w", err)
	}

	return &recipe, nil
}

func (r *Recipe) Upsert(upsert models.Recipe) (models.Recipe, bool, error) {
	// NOTE: so that COALESCE would work with the array
	var images any
	if len(upsert.Images) > 0 {
		images = pq.Array(upsert.Images)
	}

	transactions := []db.SqlBuilder{
		// upsert recipe
		psql.Insert(recipeTableName).Columns(recipeTableNameCols...).Values(
			upsert.ID,
			upsert.Name,
			upsert.Intro,
			upsert.Description,
			upsert.Instructions,
			upsert.LengthTotal,
			upsert.LengthHandsOn,
			upsert.Closing,
			upsert.RecipeVersion,
			upsert.AuthorName,
			images,
		).Suffix(`ON CONFLICT (id) DO UPDATE SET
			recipe_name =       COALESCE(NULLIF(EXCLUDED.recipe_name, ''),     recipe.recipe_name),
			intro =             COALESCE(NULLIF(EXCLUDED.intro, ''),           recipe.intro),
			description =       COALESCE(NULLIF(EXCLUDED.description, ''),     recipe.description),
			instructions =      COALESCE(NULLIF(EXCLUDED.instructions, ''),    recipe.instructions),
			length_total =      COALESCE(NULLIF(EXCLUDED.length_total, ''),    recipe.length_total),
			length_hands_on =   COALESCE(NULLIF(EXCLUDED.length_hands_on, ''), recipe.length_hands_on),
			closing =           COALESCE(NULLIF(EXCLUDED.closing, ''),         recipe.closing),
			recipe_version =    COALESCE(NULLIF(EXCLUDED.recipe_version, 0),   recipe.recipe_version),
			author_name =       COALESCE(NULLIF(EXCLUDED.author_name, ''),     recipe.author_name),
			images =            COALESCE(EXCLUDED.images, recipe.images); -- doing the check for this in go
		`),
	}

	// manually upserting tags
	if len(upsert.Tags) > 0 {
		transactions = append(transactions,
			psql.Delete(tagTableName).Where(sq.Eq{"recipe_id": upsert.ID}),
		)

		stmt := psql.Insert(tagTableName).Columns(tagTableCols...)

		for _, t := range upsert.Tags {
			stmt = stmt.Values(upsert.ID, t)
		}

		transactions = append(transactions, stmt)
	}

	if len(upsert.Ingredients) > 0 {
		transactions = append(transactions,
			psql.Delete(ingredientTableName).Where(sq.And{sq.Eq{"recipe_id": upsert.ID}, sq.Eq{"field": "INGREDIENT"}}),
		)

		stmt := psql.Insert(ingredientTableName).Columns(ingredientTableCols...)

		for i, ing := range upsert.Ingredients {
			stmt = stmt.Values(upsert.ID, i, "INGREDIENT", ing.Name, ing.Amount, ing.Unit.Name)
		}

		transactions = append(transactions, stmt)
	}

	if len(upsert.Seasonings) > 0 {
		transactions = append(transactions,
			psql.Delete(ingredientTableName).Where(sq.And{sq.Eq{"recipe_id": upsert.ID}, sq.Eq{"field": "SEASONING"}}),
		)

		stmt := psql.Insert(ingredientTableName).Columns(ingredientTableCols...)

		for i, ing := range upsert.Seasonings {
			stmt = stmt.Values(upsert.ID, i, "SEASONING", ing.Name, ing.Amount, ing.Unit.Name)
		}

		transactions = append(transactions, stmt)
	}

	trx := r.dbc.DB.MustBegin()
	// NOTE: This does not rollback if trx is done
	defer trx.Rollback()

	for i, t := range transactions {
		q, args, err := t.ToSql()
		if err != nil {
			logrus.Errorf("error getting query %d, %s: %s", i, q, err.Error())
			return models.Recipe{}, false, fmt.Errorf("error getting query: %w", err)
		}

		if args != nil {
			_, err = trx.Exec(q, args...)
		} else {
			_, err = trx.Exec(q)
		}
		if err != nil {
			logrus.Errorf("error creating transaction %s", err.Error())
			logrus.Errorf("QUERY: %s", q)
			fmt.Printf("ARGS: len(%d) %#v\n", len(args), args)
			return models.Recipe{}, false, fmt.Errorf("error creating transaction: %w", err)
		}
	}

	err := trx.Commit()
	if err != nil {
		logrus.Errorf("error committing transaction: %s", err.Error())
		return models.Recipe{}, false, fmt.Errorf("error committing transaction: %w", err)
	}

	return upsert, true, nil
}

func (r *Recipe) DeleteOne(id string) error {
	usersq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete("recipe").Where(sq.Eq{"id": id})

	_, err := usersq.RunWith(r.dbc.DB).Exec()
	if err != nil {
		return err
	}

	return nil
}
