package models

import (
	"fmt"

	"git.ve.home/nicolasc/linotte/services/user/configuration"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type databaseSettings struct {
	User     string
	Password string
	Database string
	Address  string
	Port     string
	Verbose  bool
	valid    bool
}

// Connection represents the database connextion
type Connection struct {
	*gorm.DB
}

// Accessor hosts every database access methods
type Accessor interface{}

var database Connection

// Connect establish the database connection after it has ben initialized
func Connect(settings *config.Configuration) (*Connection, error) {
	str := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.DatabaseUser,
		settings.DatabasePassword,
		settings.DatabaseHost,
		settings.DatabasePort,
		settings.DatabaseName,
	)

	db, err := gorm.Open("mysql", str)

	if err != nil {
		return nil, fmt.Errorf("Connection failed: %v", err)
	}

	database.DB = db
	database.LogMode(settings.DatabaseVerboseMode)

	return &database, nil
}

// MigrateAll migrates User
func MigrateAll() error {
	return database.AutoMigrate(
		&User{},
	).Error
}
