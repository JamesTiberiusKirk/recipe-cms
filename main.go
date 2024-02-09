package main

import (
	"github.com/JamesTiberiusKirk/recipe-cms/config"
	"github.com/JamesTiberiusKirk/recipe-cms/db"
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site"
	"github.com/sirupsen/logrus"
)

func main() {

	conf := config.GetConfig()

	dbc, err := db.Connect(conf.DbURL)
	if err != nil {
		panic(err)
	}

	// recipeRegistry := registry.NewMockRecipeRegistry()
	recipeRegistry := registry.NewRecipe(dbc)
	user := registry.NewUser(dbc)

	s := site.NewSite(conf, recipeRegistry, user)

	logrus.Info(conf.HTTPPort)

	s.Start(conf.HTTPPort)
}
