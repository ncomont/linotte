package converters

import (
	"git.ve.home/nicolasc/linotte/services/taxref/models"
	proto_taxref "git.ve.home/nicolasc/linotte/services/taxref/proto"
)

// TaxonModelsToProto converts an array of taxon models to a TaxonsReply
func TaxonModelsToProto(taxons []*models.Taxon) *proto_taxref.TaxonsReply {
	reply := &proto_taxref.TaxonsReply{}
	for _, taxon := range taxons {
		reply.Taxons = append(reply.Taxons, TaxonModelToProto(taxon))
	}
	return reply
}

// TaxonModelToProto converts a taxon model to a TaxonReply
func TaxonModelToProto(taxon *models.Taxon) *proto_taxref.TaxonReply {
	if taxon != nil {
		var (
			vernacularNames            []string
			name                       string
			author                     string
			firstLevelVernacularGroup  string
			secondLevelVernacularGroup string
		)

		if taxon.Verb != nil {
			for _, v := range taxon.Verb.VernacularNames {
				vernacularNames = append(vernacularNames, v.Value)
			}

			name = taxon.Verb.Name
			author = taxon.Verb.Author
			firstLevelVernacularGroup = taxon.Verb.FirstLevelVernacularGroup
			secondLevelVernacularGroup = taxon.Verb.SecondLevelVernacularGroup
		}

		return &proto_taxref.TaxonReply{
			Id:              uint32(taxon.ID),
			ReferenceId:     uint32(taxon.ReferenceTaxonID),
			VernacularNames: vernacularNames,
			Name:            name,
			Author:          author,
			FirstLevelVernacularGroup:  firstLevelVernacularGroup,
			SecondLevelVernacularGroup: secondLevelVernacularGroup,
			Parent: TaxonModelToProto(taxon.Parent),
			Rank:   RankModelToProto(taxon.Rank),
		}
	}
	return nil
}
