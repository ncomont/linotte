package services

import (
	"log"

	"git.ve.home/nicolasc/linotte/services/job/converters"
	"git.ve.home/nicolasc/linotte/services/job/models"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
)

// ResultService give an access to every possible result operations
type ResultService struct {
	accessor *models.JobResultAccessor
}

// NewResultService creates a new result service from the given accessor
func NewResultService(accessor *models.JobResultAccessor) *ResultService {
	return &ResultService{accessor}
}

// ResultsByState returns every results with the given state
func (service *ResultService) ResultsByState(in *proto_job.ResultRequest) (*proto_job.ResultsReply, error) {
	results, err := service.accessor.GetByState(in.State)
	return converters.ResultModelsToProto(results), err
}

// ResultsByStateAndReportID returns every report's results with the given state
func (service *ResultService) ResultsByStateAndReportID(in *proto_job.ResultRequest) (*proto_job.ResultsReply, error) {
	results, err := service.accessor.GetByStateAndReportID(uint(in.ReportId), in.State)
	return converters.ResultModelsToProto(results), err
}

// ResolveConflict set the given result as FOUND with the given taxon id
func (service *ResultService) ResolveConflict(in *proto_job.ResultRequest) (*proto_job.ResultReply, error) {
	result, err := service.accessor.GetByID(uint(in.Id))

	if err != nil {
		return nil, err
	} else if result.ID == 0 {
		log.Fatal("Conflict not found in SetResultState")
	}

	result.State = "RESOLVED"
	result.TaxonID = uint(in.TaxonId)
	err = service.accessor.Save(result)

	return converters.ResultModelToProto(result), err
}

// UnresolveConflict set the given result as FOUND with the given taxon id
func (service *ResultService) UnresolveConflict(in *proto_job.ResultRequest) (*proto_job.ResultReply, error) {
	result, err := service.accessor.GetByID(uint(in.Id))

	if err != nil {
		return nil, err
	} else if result.ID == 0 {
		log.Fatal("Conflict not found in SetResultState")
	}

	result.State = "NOT_FOUND"
	result.TaxonID = uint(in.TaxonId)
	err = service.accessor.Save(result)

	return converters.ResultModelToProto(result), err
}
