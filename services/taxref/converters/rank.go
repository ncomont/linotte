package converters

import (
	"git.ve.home/nicolasc/linotte/services/taxref/models"
	proto_taxref "git.ve.home/nicolasc/linotte/services/taxref/proto"
)

// RankModelsToProto converts an array of ranks to a RanksReply
func RankModelsToProto(ranks []*models.Rank) *proto_taxref.RanksReply {
	return &proto_taxref.RanksReply{
		Ranks: RankModelsToProtos(ranks),
	}
}

// RankModelsToProtos converts a ranks array to a RankReply array
func RankModelsToProtos(ranks []*models.Rank) []*proto_taxref.RankReply {
	var replies []*proto_taxref.RankReply
	for _, rank := range ranks {
		replies = append(replies, RankModelToProto(rank))
	}
	return replies
}

// RankModelToProto converts a rank mofrl to a RankReply
func RankModelToProto(rank *models.Rank) *proto_taxref.RankReply {
	if rank != nil {
		return &proto_taxref.RankReply{
			Id:    uint32(rank.ID),
			Key:   rank.Key,
			Name:  rank.Name,
			Order: uint32(rank.Order),
		}
	}
	return nil
}
