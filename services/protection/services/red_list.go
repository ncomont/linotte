package services

import (
	"git.ve.home/nicolasc/linotte/services/job/models"
	proto_protection "git.ve.home/nicolasc/linotte/services/protection/proto"
)

// RedListService give an access to every possible result operations
type RedListService struct {
	accessor *models.JobResultAccessor
}

// NewRedListService creates a new result service from the given accessor
func NewRedListService(accessor *models.JobResultAccessor) *RedListService {
	return &RedListService{accessor}
}

// ProtectionsByTaxonIDs returns every redlist entries related to the given taxon id and loc id
func (service *RedListService) ProtectionsByTaxonIDs(in *proto_protection.ProtectionRequest) (*proto_protection.ProtectionReply, error) {
	return nil, nil
}
