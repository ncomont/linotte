package services

import (
	"git.ve.home/nicolasc/linotte/services/job/converters"
	"git.ve.home/nicolasc/linotte/services/job/models"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
)

// ReportService give an access to every possible report operations
type ReportService struct {
	accessor *models.JobReportAccessor
}

// NewReportService creates a new report service from the given accessor
func NewReportService(accessor *models.JobReportAccessor) *ReportService {
	return &ReportService{accessor}
}

// ReportByID returns a single report from its ID
func (service *ReportService) ReportByID(in *proto_job.ReportRequest) (*proto_job.ReportReply, error) {
	report, err := service.accessor.GetByID(uint(in.Id))
	return converters.ReportModelToProto(report), err
}

// ReportStatisticsByReportID returns a statistics object from the given report ID
func (service *ReportService) ReportStatisticsByReportID(in *proto_job.ReportRequest) (*proto_job.ReportStatisticsReply, error) {
	report, err := service.accessor.StatisticsByID(uint(in.Id))
	return converters.ReportStatisticsModelToProto(report), err
}

// ReportsByJobID returns every reports attached to the given job id
func (service *ReportService) ReportsByJobID(in *proto_job.JobRequest) (*proto_job.ReportsReply, error) {
	reports, err := service.accessor.GetByJobID(uint(in.Id))
	return converters.ReportModelsToProto(reports), err
}

// SaveReport save the given report in database
func (service *ReportService) SaveReport(in *proto_job.ReportRequest) (*proto_job.ReportReply, error) {
	report := converters.ReportProtoToModel(in)
	err := service.accessor.Save(report)
	return converters.ReportModelToProto(report), err
}
