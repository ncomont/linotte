package services

import (
	"git.ve.home/nicolasc/linotte/services/taxref/converters"
	"git.ve.home/nicolasc/linotte/services/taxref/models"
	proto_taxref "git.ve.home/nicolasc/linotte/services/taxref/proto"
)

// TaxonService give an access to every possible taxon operations
type TaxonService struct {
	accessor *models.TaxonAccessor
}

// NewTaxonService creates a new taxon service from the given accessor
func NewTaxonService(taxonAccessor *models.TaxonAccessor) *TaxonService {
	return &TaxonService{
		accessor: taxonAccessor,
	}
}

// TaxonByID returns a single taxon fro its ID
func (service *TaxonService) TaxonByID(in *proto_taxref.TaxonRequest) (*proto_taxref.TaxonReply, error) {
	taxon, err := service.accessor.GetByID(uint(in.Id))
	return converters.TaxonModelToProto(taxon), err
}

// ReferenceByID returns the reference taxon for the given taxon ID
func (service *TaxonService) ReferenceByID(in *proto_taxref.TaxonRequest) (*proto_taxref.TaxonReply, error) {
	taxon, err := service.accessor.GetReferenceByID(uint(in.Id))
	return converters.TaxonModelToProto(taxon), err
}

// TaxonsByIDs returns a single taxon fro its ID
func (service *TaxonService) TaxonsByIDs(in *proto_taxref.TaxonsRequest) (*proto_taxref.TaxonsReply, error) {
	var ids []uint

	for _, id := range in.Ids {
		ids = append(ids, uint(id))
	}

	taxons, err := service.accessor.GetByIds(ids)
	return converters.TaxonModelsToProto(taxons), err
}

// ReferenceAndSynonymsForTaxonID returns an array of taxons corresponding to the synonyms of the given taxon id
func (service *TaxonService) ReferenceAndSynonymsForTaxonID(in *proto_taxref.TaxonRequest) (*proto_taxref.TaxonsReply, error) {
	taxons, err := service.accessor.GetReferenceAndSynonymsForID(uint(in.Id))
	return converters.TaxonModelsToProto(taxons), err
}

// TaxonClassificationForID returns a taxon with populated Parent property depending on the depth choosen
func (service *TaxonService) TaxonClassificationForID(in *proto_taxref.TaxonClassificationRequest) (*proto_taxref.TaxonReply, error) {
	taxon, err := service.accessor.GetClassificationForID(uint(in.Id), int(in.Depth))
	return converters.TaxonModelToProto(taxon), err
}

// ReferenceByVerb returns a taxon corresponding to the given verb
func (service *TaxonService) ReferenceByVerb(in *proto_taxref.TaxonRequest) (*proto_taxref.TaxonReply, error) {
	var (
		taxon *models.Taxon
		err   error
	)

	if len(in.VernacularName) > 0 {
		taxon, err = service.accessor.GetReferenceByVernacularName(in.VernacularName)
	} else if len(in.Name) > 0 {
		taxon, err = service.accessor.GetReferenceByName(in.FullName)
	} else if len(in.FullName) > 0 {
		taxon, err = service.accessor.GetReferenceByFullName(in.FullName, in.IgnorePunctuation)
	}

	return converters.TaxonModelToProto(taxon), err
}
