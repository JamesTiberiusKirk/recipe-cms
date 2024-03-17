package models

import "strings"

type User struct {
	Username string `json:"name"     goscanql:"name"`
	Password string `json:"password" goscanql:"password"`
}

type Recipe struct {
	ID            string       `json:"id"              goscanql:"id"`
	Name          string       `json:"name"            goscanql:"recipename"`
	Intro         string       `json:"intro"           goscanql:"intro"`
	Description   string       `json:"description"     goscanql:"description"`
	Ingredients   []Ingredient `json:"ingredients"     goscanql:"ingredients"`
	Seasonings    []Ingredient `json:"seasonings"      goscanql:"seasonings"`
	Instructions  string       `json:"instructions"    goscanql:"instructions"`
	LengthTotal   string       `json:"length_total"    goscanql:"lengthtotal"`
	LengthHandsOn string       `json:"length_hands_on" goscanql:"lengthhandson"`
	Closing       string       `json:"closing"         goscanql:"closing"`
	Tags          []string     `json:"tags"            goscanql:"tags"`
	RecipeVersion int          `json:"recipe_version"  goscanql:"recipeversion"`
	AuthorName    string       `json:"author_name"     goscanql:"authorname"`
	Images        []string     `json:"images"          goscanql:"images"`
}

func (r Recipe) String() string {
	tags := ""
	for _, t := range r.Tags {
		tags += t + " "
	}

	return strings.ToUpper(r.Name + " " +
		r.Intro + " " +
		r.Description + " " +
		r.Instructions + " " +
		r.Closing + " " +
		tags + " " +
		r.AuthorName)
}

type Ingredient struct {
	Name   string `json:"name"   goscanql:"name"`
	Amount int    `json:"amount" goscanql:"amount"`
	Unit   Unit   `json:"unit"   goscanql:"unit"`
}

var (
	DefaultSystemUnits = []Unit{
		{
			DisplayName: "kg",
			Name:        "kg",
		},
		{
			DisplayName: "g",
			Name:        "g",
		},
		{
			DisplayName: "l",
			Name:        "l",
		},
		{
			DisplayName: "ml",
			Name:        "ml",
		},
		{
			DisplayName: "unit",
			Name:        "unit",
		},
	}
)

type Unit struct {
	DisplayName string `json:"display_name" goscanql:"name"`
	Name        string `json:"unit_name"    goscanql:"unitname"`
}

// TODO: need to make a converter here
