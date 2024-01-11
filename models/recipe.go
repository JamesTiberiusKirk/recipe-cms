package models

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
