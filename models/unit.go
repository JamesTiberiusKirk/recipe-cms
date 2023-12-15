package models

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
