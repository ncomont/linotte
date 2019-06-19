package models

import (
	"fmt"

	"git.ve.home/nicolasc/linotte/services/taxref/configuration"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Connection struct {
	*gorm.DB
}

type Accessor interface{}

var database Connection

// Connect establish the database connection after it has ben initialized
func Connect(settings *config.Configuration) (*Connection, error) {
	connnectionString := fmt.Sprintf(
		"%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		settings.DatabaseUser,
		settings.DatabasePassword,
		settings.DatabaseHost,
		settings.DatabasePort,
		settings.DatabaseName,
	)

	db, err := gorm.Open("mysql", connnectionString)

	if err != nil {
		return nil, fmt.Errorf("Connection failed: %v", err)
	}

	database.DB = db
	database.LogMode(settings.DatabaseVerboseMode)

	return &database, nil
}

// MigrateAll migrates every tables used by the taxref service
func MigrateAll() error {
	return database.AutoMigrate(
		&Rank{},
		&Taxon{},
		&Verb{},
		&VernacularName{},
	).Error
}
