package services

import (
	"git.ve.home/nicolasc/linotte/services/taxref/converters"
	"git.ve.home/nicolasc/linotte/services/taxref/models"
	proto_taxref "git.ve.home/nicolasc/linotte/services/taxref/proto"
)

// RankService give an access to every possible rank operations
type RankService struct {
	accessor *models.RankAccessor
}

// NewRankService creates a new rank service from the given accessor
func NewRankService(rankAccessor *models.RankAccessor) *RankService {
	return &RankService{rankAccessor}
}

// AllRanks returns every ranks
func (service *RankService) AllRanks() (*proto_taxref.RanksReply, error) {
	ranks, err := service.accessor.GetAll()
	return converters.RankModelsToProto(ranks), err
}
