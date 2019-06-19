package converters

import (
	"git.ve.home/nicolasc/linotte/services/job/models"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
)

// ResultModelsToProto converts an array of job result to a ResultsReply
func ResultModelsToProto(in []*models.JobResult) *proto_job.ResultsReply {
	reply := &proto_job.ResultsReply{}

	for _, result := range in {
		reply.Results = append(reply.Results, ResultModelToProto(result))
	}

	return reply
}

// ResultModelToProto converts a job result model to a ResultReply
func ResultModelToProto(in *models.JobResult) *proto_job.ResultReply {
	if in != nil {
		return &proto_job.ResultReply{
			Id:         uint32(in.ID),
			State:      in.State,
			SearchData: in.SearchData,
			Value:      in.Value,
			TaxonId:    uint32(in.TaxonID),
		}
	}

	return nil
}

// ResultProtoToModel converts a result reply to a job result model
func ResultProtoToModel(in *proto_job.ResultReply) *models.JobResult {
	if in != nil {
		return &models.JobResult{
			SearchData: in.SearchData,
			State:      in.State,
			TaxonID:    uint(in.TaxonId),
			Value:      in.Value,
		}
	}

	return nil
}

// ResultProtosToModels converts an array of result replies to a job result model
func ResultProtosToModels(in []*proto_job.ResultReply) []*models.JobResult {
	var results []*models.JobResult

	for _, result := range in {
		results = append(results, ResultProtoToModel(result))
	}

	return results
}
