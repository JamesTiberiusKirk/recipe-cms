package main

import (
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site"
)

func main() {

	// conf := config.GetConfig()

	// dbc, err := db.Connect(conf.DbURL)
	// if err != nil {
	// 	panic(err)
	// }

	// recipeRegistry := registry.NewRecipe(dbc)

	test := "some data"
	print(test)

	recipeRegistry := registry.NewMockRecipeRegistry()
	s := site.NewSite(recipeRegistry)
	s.Start("localhost:5000")
}
