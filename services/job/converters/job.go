package converters

import (
	"git.ve.home/nicolasc/linotte/services/job/models"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
	"github.com/golang/protobuf/ptypes"
)

// JobModelsToProto converts a job models array to a JobsReply
func JobModelsToProto(jobs []models.Job) *proto_job.JobsReply {
	reply := &proto_job.JobsReply{}
	for _, job := range jobs {
		reply.Jobs = append(reply.Jobs, JobModelToProto(&job))
	}
	return reply
}

// JobModelToProto converts a job model to a JobReply
func JobModelToProto(job *models.Job) *proto_job.JobReply {
	if job != nil {
		time, _ := ptypes.TimestampProto(job.LastUpdate)
		return &proto_job.JobReply{
			Id:         uint32(job.ID),
			Name:       job.Name,
			File:       job.File,
			LastUpdate: time,
			Status:     job.Status,
			Type:       job.Type,
			Data:       job.Data,
			Reports:    ReportModelsToProto(job.Reports).Reports,
		}
	}
	return nil
}
