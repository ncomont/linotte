package console

import "git.ve.home/nicolasc/linotte/models"

type Console struct {
	Users *UsersConsole
}

var instance *Console = nil

func Initialize() (*Console, error) {
	var (
		err        error
		connection *models.Connection
	)

	if instance == nil {
		models.InitDatabase()
		if connection, err = models.Connect(); err == nil {
			if err = models.MigrateAll(); err == nil {
				instance = &Console{
					Users: InitUsersConsole(models.NewUserAccessor(connection.DB)),
				}
			}
		}
	}

	return instance, err
}
