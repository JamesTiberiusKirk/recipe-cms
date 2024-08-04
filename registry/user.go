package registry

import (
	"github.com/JamesTiberiusKirk/recipe-cms/db"
	"github.com/JamesTiberiusKirk/recipe-cms/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/rustedturnip/goscanql"
)

type IUser interface {
	GetOneByUsername(username string) (models.User, error)
}

type User struct {
	dbc *db.DB
}

func NewUser(dbc *db.DB) *User {
	return &User{
		dbc: dbc,
	}
}

// GetOneByUsername retrieves a user by username, case-insensitive
func (u *User) GetOneByUsername(username string) (models.User, error) {
	usersq := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("username, password").From("author").Where("LOWER(username)=LOWER(?)", username)

	rows, err := usersq.RunWith(u.dbc.DB).Query()
	if err != nil {
		return models.User{}, err
	}

	user, err := goscanql.RowsToStruct[models.User](rows)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
