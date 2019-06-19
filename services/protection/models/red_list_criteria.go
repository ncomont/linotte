package models

import "github.com/jinzhu/gorm"

// RedListCriteria is the data model used to represents a red list criteria in the database
type RedListCriteria struct {
	ID uint `gorm:"primary_key"`
}

// RedListCriteriaAccessor hosts every red list criteria related methods
type RedListCriteriaAccessor struct {
	*gorm.DB
}

// NewRedListCriteriaAccessor helps to instanciate a new red list criteria accessor with the given database connection
func NewRedListCriteriaAccessor(db *gorm.DB) *RedListCriteriaAccessor {
	return &RedListCriteriaAccessor{db}
}

// Create add a new red list criteria in the database
func (accessor *RedListCriteriaAccessor) Create(criteria *RedListCriteria) error {
	return accessor.DB.Create(criteria).Error
}

// Save helps to persist modifications on a criteria
func (accessor *RedListCriteriaAccessor) Save(criteria *RedListCriteria) error {
	return accessor.DB.Save(criteria).Error
}

// GetByID returns a single red list criteria from its ID
func (accessor *RedListCriteriaAccessor) GetByID(id uint) (*RedListCriteria, error) {
	var criteria *RedListCriteria
	err := accessor.DB.First(criteria, id).Error
	return criteria, err
}
