package models

import "github.com/jinzhu/gorm"

// RedList is the data model used to represents a red list in the database
type RedList struct {
	ID uint `gorm:"primary_key"`

	RedListEntry []*RedListEntry
}

// RedListAccessor hosts every red list related methods
type RedListAccessor struct {
	*gorm.DB
}

// NewRedListAccessor helps to instanciate a new red list accessor with the given database connection
func NewRedListAccessor(db *gorm.DB) *RedListAccessor {
	return &RedListAccessor{db}
}

// Create add a new red list in the database
func (accessor *RedListAccessor) Create(redlist *RedList) error {
	return accessor.DB.Create(redlist).Error
}

// Save helps to persist modifications on a red list
func (accessor *RedListAccessor) Save(redlist *RedList) error {
	return accessor.DB.Save(redlist).Error
}

// GetByID returns a single red list from its ID
func (accessor *RedListAccessor) GetByID(id uint) (*RedList, error) {
	var redlist *RedList
	err := accessor.DB.First(redlist, id).Error
	return redlist, err
}
