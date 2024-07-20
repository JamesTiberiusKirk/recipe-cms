package models

import "strings"

type User struct {
	Username string `json:"name"     sql:"name"`
	Password string `json:"password" sql:"password"`
}

type Recipe struct {
	ID            string       `json:"id"              sql:"id"`
	Name          string       `json:"name"            sql:"recipename"`
	Intro         string       `json:"intro"           sql:"intro"`
	Description   string       `json:"description"     sql:"description"`
	Ingredients   []Ingredient `json:"ingredients"     sql:"ingredients"`
	Seasonings    []Ingredient `json:"seasonings"      sql:"seasonings"`
	Instructions  string       `json:"instructions"    sql:"instructions"`
	LengthTotal   string       `json:"length_total"    sql:"lengthtotal"`
	LengthHandsOn string       `json:"length_hands_on" sql:"lengthhandson"`
	Closing       string       `json:"closing"         sql:"closing"`
	Tags          []string     `json:"tags"            sql:"tags"`
	RecipeVersion int          `json:"recipe_version"  sql:"recipeversion"`
	AuthorName    string       `json:"author_name"     sql:"authorname"`
	Images        []string     `json:"images"          sql:"images"`
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
	ArrayIndex int    `json:"array_index" sql:"array_index"`
	Name       string `json:"name" sql:"name"`
	Amount     int    `json:"amount" sql:"amount"`
	Unit       Unit   `json:"unit" sql:"unit"`
}

var (
	DefaultSystemUnits = []Unit{
		{
			DisplayName: "g",
			Name:        "g",
		},
		{
			DisplayName: "ml",
			Name:        "ml",
		},
		{
			DisplayName: "Can",
			Name:        "can",
		},
		{
			DisplayName: "Clove",
			Name:        "clove",
		},
		{
			DisplayName: "Unit",
			Name:        "unit",
		},
		{
			DisplayName: "Large",
			Name:        "large",
		},
		{
			DisplayName: "Medium",
			Name:        "medium",
		},
		{
			DisplayName: "Small",
			Name:        "small",
		},
	}
)

type Unit struct {
	DisplayName string `json:"display_name" sql:"name"`
	Name        string `json:"unit_name"    sql:"unitname"`
}

// TODO: need to make a converter here
