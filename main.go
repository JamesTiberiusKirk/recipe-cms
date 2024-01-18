package main

import (
	"github.com/JamesTiberiusKirk/recipe-cms/config"
	"github.com/JamesTiberiusKirk/recipe-cms/db"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site"
)

func main() {

	conf := config.GetConfig()

	dbc, err := db.Connect(conf.DbURL)
	if err != nil {
		panic(err)
	}

	recipeRegistry := registry.NewRecipe(dbc)

	s := site.NewSite(recipeRegistry)
	s.Start("localhost:5000")
}
