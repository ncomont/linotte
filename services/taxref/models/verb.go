package models

import "github.com/jinzhu/gorm"

type Verb struct {
	VernacularNames            []VernacularName `gorm:"many2many:verbs_vernacular_names"`
	Name                       string           `gorm:"type:varchar(1024)"`
	Author                     string           `gorm:"type:varchar(1024)"`
	FirstLevelVernacularGroup  string
	SecondLevelVernacularGroup string
	ReferenceTaxon             *Taxon
	Taxon                      *Taxon

	// Database references
	ReferenceTaxonID uint `gorm:"index"`
	TaxonID          uint `sql:"not null;type:int(10) unsigned" gorm:"primary_key"`
}

type VerbAccessor struct {
	*gorm.DB
}

func NewVerbAccessor(db *gorm.DB) *VerbAccessor {
	return &VerbAccessor{db}
}

func (accessor *VerbAccessor) Create(verb *Verb) (*Verb, error) {
	err := accessor.DB.Create(verb).Error
	return verb, err
}

func (accessor *VerbAccessor) GetAll(limit int, offset int) []Verb {
	verbs := []Verb{}
	accessor.DB.Preload("VernacularNames").Limit(limit).Offset(offset).Find(&verbs)
	return verbs
}

func (accessor *VerbAccessor) GetTotalCount() int {
	var count int
	accessor.DB.Table("verbs").Count(&count)
	return count
}
