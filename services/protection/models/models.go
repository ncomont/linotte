package models

import (
	"fmt"

	"git.ve.home/nicolasc/linotte/services/job/configuration"

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

// MigrateAll migrates RedList, RedListEntry, and RedListCriteria
func MigrateAll() error {
	return database.AutoMigrate(
		&RedList{},
		&RedListEntry{},
		&RedListCriteria{},
	).Error
}
