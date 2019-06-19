package models

import "github.com/jinzhu/gorm"

// RedListEntry is the data model used to represents a red list entry in the database
type RedListEntry struct {
	ID uint `gorm:"primary_key"`

	RedList         *RedList
	RedListCriteria *RedListCriteria

	// Database references
	RedListID         uint `gorm:"index"`
	RedListCriteriaID uint `gorm:"index"`
}

// RedListEntryAccessor hosts every red list entries related methods
type RedListEntryAccessor struct {
	*gorm.DB
}

// NewRedListEntryAccessor helps to instanciate a new red list entry accessor with the given database connection
func NewRedListEntryAccessor(db *gorm.DB) *RedListEntryAccessor {
	return &RedListEntryAccessor{db}
}

// Create add a new red list entry in the database
func (accessor *RedListEntryAccessor) Create(entry *RedListEntry) error {
	return accessor.DB.Create(entry).Error
}

// Save helps to persist modifications on a entry
func (accessor *RedListEntryAccessor) Save(entry *RedListEntry) error {
	return accessor.DB.Save(entry).Error
}

// GetByID returns a single red list entry from its ID
func (accessor *RedListEntryAccessor) GetByID(id uint) (*RedListEntry, error) {
	var entry *RedListEntry
	err := accessor.DB.First(entry, id).Error
	return entry, err
}
