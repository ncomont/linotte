package converters

import (
	"git.ve.home/nicolasc/linotte/services/job/models"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
	"github.com/golang/protobuf/ptypes"
)

// ReportModelsToProto converts a job reports array to a ReportsReply
func ReportModelsToProto(reports []*models.JobReport) *proto_job.ReportsReply {
	reply := &proto_job.ReportsReply{}
	for _, report := range reports {
		reply.Reports = append(reply.Reports, ReportModelToProto(report))
	}
	return reply
}

// ReportModelToProto converts a job report model to a ReportReply
func ReportModelToProto(report *models.JobReport) *proto_job.ReportReply {
	if report != nil {
		time, _ := ptypes.TimestampProto(report.CreationDate)
		return &proto_job.ReportReply{
			Id:           uint32(report.ID),
			Message:      report.Message,
			Status:       report.Status,
			CreationDate: time,
			JobId:        uint32(report.JobID),
			// Results:      ResultModelsToProto(report.Results),
		}
	}
	return nil
}

// ReportProtoToModel conserts a ReportRequest to a JobReport model
func ReportProtoToModel(report *proto_job.ReportRequest) *models.JobReport {
	if report != nil {
		return &models.JobReport{
			ID:      uint(report.Id),
			JobID:   uint(report.JobId),
			Results: ResultProtosToModels(report.Results),
			Message: report.Message,
			Status:  report.Status,
		}
	}

	return nil
}

// ReportStatisticsModelToProto converts a ReportStatistics model to a ReportStatisticsReply
func ReportStatisticsModelToProto(stats *models.ReportStatistics) *proto_job.ReportStatisticsReply {
	if stats != nil {
		return &proto_job.ReportStatisticsReply{
			FoundCount:    uint32(stats.FoundCount),
			NotFoundCount: uint32(stats.NotFoundCount),
			ResolvedCount: uint32(stats.ResolvedCount),
		}
	}
	return nil
}
