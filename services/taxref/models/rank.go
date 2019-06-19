package models

import "github.com/jinzhu/gorm"

type Rank struct {
	ID    uint   `gorm:"primary_key"`
	Key   string `gorm:"unique"`
	Name  string
	Order int
}

type RankAccessor struct {
	*gorm.DB
}

func NewRankAccessor(db *gorm.DB) *RankAccessor {
	return &RankAccessor{db}
}

func (accessor *RankAccessor) Create(rank *Rank) (*Rank, error) {
	err := accessor.DB.Create(rank).Error
	return rank, err
}

func (accessor *RankAccessor) GetAll() ([]*Rank, error) {
	var ranks []*Rank
	err := accessor.DB.Find(&ranks).Error
	return ranks, err
}
