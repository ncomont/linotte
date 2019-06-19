package console

import (
	"git.ve.home/nicolasc/linotte/models"
	"golang.org/x/crypto/bcrypt"
)

type UsersConsole struct {
	accessor *models.UserAccessor
}

var usersConsole *UsersConsole = nil

func InitUsersConsole(accessor *models.UserAccessor) *UsersConsole {
	if usersConsole == nil {
		usersConsole = &UsersConsole{
			accessor: accessor,
		}
	}
	return usersConsole
}

func (console *UsersConsole) Create(username string, password string) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	return console.accessor.Create(username, string(hash[:]))
}
