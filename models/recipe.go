package models

type Recipe struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
	Seasonings  []Ingredient `json:"seasonings"`
}

type Ingredient struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
	Unit   Unit   `json:"unit"`
}
