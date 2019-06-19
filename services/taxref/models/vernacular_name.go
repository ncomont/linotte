package models

import "github.com/jinzhu/gorm"

type VernacularName struct {
	ID    uint   `gorm:"primary_key"`
	Value string `gorm:"type:varchar(2048)"`
	Verbs []Verb `gorm:"many2many:verbs_vernacular_names"`
}

type VernacularNameAccessor struct {
	*gorm.DB
}

func NewVernacularNameAccessor(db *gorm.DB) *VernacularNameAccessor {
	return &VernacularNameAccessor{db}
}

func (accessor *VernacularNameAccessor) Create(vernacularName *VernacularName) (*VernacularName, error) {
	err := accessor.DB.Create(vernacularName).Error
	return vernacularName, err
}
