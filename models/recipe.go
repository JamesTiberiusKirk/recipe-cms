package models

type User struct {
	Name string `json:"name"`
}

type Recipe struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Intro         string       `json:"intro"`
	Description   string       `json:"description"`
	Ingredients   []Ingredient `json:"ingredients"`
	Seasonings    []Ingredient `json:"seasonings"`
	Instructions  string       `json:"instructions"`
	LengthTotal   string       `json:"length_total"`
	LengthHandsOn string       `json:"length_hands_on"`
	Closing       string       `json:"closing"`
	Tags          []string     `json:"tags"`
	Version       int          `json:"version"`
	Author        User         `json:"author"`
}

type Ingredient struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
	Unit   Unit   `json:"unit"`
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
	DisplayName string `json:"display_name"`
	Name        string `json:"name"`
}

// TODO: need to make a converter here
