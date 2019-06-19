package models

import (
	"sort"
	"strings"

	"git.ve.home/nicolasc/linotte/libs/helpers"

	"github.com/jinzhu/gorm"
)

type Taxon struct {
	ID             uint `sql:"not null;type:int(10) unsigned" gorm:"primary_key"`
	ReferenceTaxon *Taxon
	Verb           *Verb `gorm:"ForeignKey:TaxonID"`
	Parent         *Taxon
	Rank           *Rank

	// Database references
	RankID   uint `gorm:"index"`
	ParentID uint `gorm:"index"`

	ReferenceTaxonID uint `gorm:"index"`
}

type TaxonAccessor struct {
	*gorm.DB
}

func NewTaxonAccessor(db *gorm.DB) *TaxonAccessor {
	return &TaxonAccessor{db}
}

func (accessor *TaxonAccessor) Create(taxon *Taxon) (*Taxon, error) {
	err := accessor.DB.Create(taxon).Error
	return taxon, err
}

// GetReferenceByID returns the reference taxon, wathever which Taxref ID is given
func (accessor *TaxonAccessor) GetReferenceByID(id uint) (*Taxon, error) {
	var err error
	taxon := Taxon{}
	err = helpers.FilterGormError(accessor.DB.Preload("Rank").Preload("Verb").First(&taxon, id).Error)

	if err != nil && taxon.ID != taxon.ReferenceTaxonID {
		reference := Taxon{}
		err = helpers.FilterGormError(accessor.DB.Preload("Rank").Preload("Verb").First(&reference, taxon.ReferenceTaxonID).Error)
		return &reference, err
	}

	return &taxon, err
}

func (accessor *TaxonAccessor) GetReferenceByFullName(name string, ignorePunctuation bool) (*Taxon, error) {
	taxon := Taxon{}
	statement := "CONCAT_WS(' ', verbs.name, verbs.author) LIKE ?"

	name = strings.Trim(name, " ")
	if ignorePunctuation {
		name = strings.Replace(name, ",", " ", -1)
		name = strings.Replace(name, ";", " ", -1)
		name = strings.Replace(name, ".", " ", -1)
		statement = "REPLACE(REPLACE(REPLACE(CONCAT_WS(' ', verbs.name, verbs.author), ',', ' '), ';', ' '), '.', ' ') LIKE ?"
	}

	err := accessor.DB.Table("taxons").
		Joins("LEFT JOIN verbs ON taxons.id = verbs.taxon_id").
		Where(statement, name).
		Preload("ReferenceTaxon").
		First(&taxon).
		Error
	err = helpers.FilterGormError(err)

	if taxon.ReferenceTaxonID != 0 {
		return taxon.ReferenceTaxon, err
	}
	return &taxon, err

}

func (accessor *TaxonAccessor) GetReferenceByName(name string) (*Taxon, error) {
	taxon := Taxon{}

	if err := accessor.DB.Table("taxons").
		Joins("LEFT JOIN verbs ON taxons.id = verbs.taxon_id").
		Where("verbs.name LIKE ?", strings.Trim(name, " ")).
		Preload("ReferenceTaxon").
		First(&taxon).
		Error; helpers.FilterGormError(err) != nil {
		return &taxon, err
	}

	if taxon.ReferenceTaxonID != 0 {
		return taxon.ReferenceTaxon, nil
	}
	return &taxon, nil
}

func (accessor *TaxonAccessor) GetReferenceByVernacularName(name string) (*Taxon, error) {
	vernacularNames := []VernacularName{}

	err := accessor.DB.Table("vernacular_names").
		Preload("Verbs").
		Preload("Verbs.ReferenceTaxon").
		Where("vernacular_names.value LIKE ?", strings.Trim(name, " ")).
		Find(&vernacularNames).
		Error
	err = helpers.FilterGormError(err)

	if err != nil || len(vernacularNames) == 0 || len(vernacularNames) > 1 {
		return &Taxon{}, err
	}

	for _, verb := range vernacularNames[0].Verbs {
		if verb.ReferenceTaxonID != vernacularNames[0].Verbs[0].ReferenceTaxonID {
			return &Taxon{}, err
		}
	}

	return vernacularNames[0].Verbs[0].ReferenceTaxon, err
}

func (accessor *TaxonAccessor) GetByID(id uint) (*Taxon, error) {
	taxon := Taxon{}
	err := accessor.DB.
		Preload("Rank").
		Preload("ReferenceTaxon").
		Preload("Verb").
		Preload("Verb.VernacularNames").
		First(&taxon, id).Error
	return &taxon, helpers.FilterGormError(err)
}

func (accessor *TaxonAccessor) GetReferenceAndSynonymsForID(id uint) ([]*Taxon, error) {
	var err error
	var taxons []*Taxon
	reference := Taxon{}

	err = helpers.FilterGormError(accessor.DB.First(&reference, id).Error)
	if err != nil {
		return nil, err
	}

	err = accessor.DB.
		Table("taxons").
		Where("taxons.reference_taxon_id = ?", reference.ReferenceTaxonID).
		Preload("ReferenceTaxon").
		Preload("Rank").
		Preload("Verb").
		Preload("Verb.VernacularNames").
		Find(&taxons).Error
	err = helpers.FilterGormError(err)

	return taxons, err
}

func (accessor *TaxonAccessor) GetClassificationForID(id uint, depth int) (*Taxon, error) {
	taxon, err := accessor.GetByID(id)
	if err != nil {
		return nil, err
	}

	parent := taxon

	for i := 0; i < depth && parent.ParentID != 0; i++ {
		parent.Parent = new(Taxon)
		accessor.DB.Model(&parent).
			Preload("ReferenceTaxon").
			Preload("Rank").
			Preload("Verb").
			Related(parent.Parent, "ParentID")
		parent = parent.Parent
	}

	return taxon, nil
}

func (accessor *TaxonAccessor) GetTotalCount() int {
	count := 0
	accessor.DB.Table("taxons").Count(&count)
	return count
}

func (accessor *TaxonAccessor) GetAll(limit int, offset int) []Taxon {
	taxons := []Taxon{}
	accessor.DB.Preload("Rank").Preload("Verb").Preload("Verb.VernacularNames").
		Limit(limit).Offset(offset).
		Find(&taxons)
	return taxons
}

func (accessor *TaxonAccessor) GetByIds(ids []uint) ([]*Taxon, error) {
	var taxons []*Taxon
	var ordered []*Taxon

	err := accessor.DB.
		Preload("ReferenceTaxon").
		Preload("Rank").
		Preload("Verb").
		Preload("Verb.VernacularNames").
		Where(ids).
		Find(&taxons).Error
	err = helpers.FilterGormError(err)

	if err != nil {
		return ordered, err
	}

	for _, id := range ids {
		i := sort.Search(len(taxons), func(i int) bool {
			return taxons[i].ID >= id
		})
		ordered = append(ordered, taxons[i])
	}

	return ordered, nil
}
