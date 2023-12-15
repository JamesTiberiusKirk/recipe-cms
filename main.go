package main

import (
	"github.com/JamesTiberiusKirk/recipe-cms/registry"
	"github.com/JamesTiberiusKirk/recipe-cms/site"
)

func main() {

	recipeRegistry := registry.NewRecipeRegistry()

	s := site.NewSite(recipeRegistry)
	s.Start("localhost:5000")
}
