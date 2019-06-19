package services

import (
	"errors"
	"log"

	"git.ve.home/nicolasc/linotte/services/job/converters"
	"git.ve.home/nicolasc/linotte/services/job/models"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
)

// JobService give an access to every possible job operations
type JobService struct {
	accessor *models.JobAccessor
	stack    chan *proto_job.JobReply
	stacked  chan TaskRequestStatus
}

// TODO: move in stack helper
type TaskRequestStatus struct {
	Success bool
	Err     error
}

// NewJobService creates a new job service from the given accessor
func NewJobService(accessor *models.JobAccessor, stack chan *proto_job.JobReply, stacked chan TaskRequestStatus) *JobService {
	return &JobService{accessor, stack, stacked}
}

// JobByID returns a single job from its ID
func (service *JobService) JobByID(in *proto_job.JobRequest) (*proto_job.JobReply, error) {
	job, err := service.accessor.GetByID(uint(in.Id))
	return converters.JobModelToProto(job), err
}

// AllJobs returns every jobs
func (service *JobService) AllJobs(in *proto_job.JobsRequest) (*proto_job.JobsReply, error) {
	jobs, err := service.accessor.GetAll()
	return converters.JobModelsToProto(jobs), err
}

// UpdateJobStatus updates the given job's status and put it in the stack if necessary
func (service *JobService) UpdateJobStatus(in *proto_job.JobRequest) (*proto_job.JobReply, error) {
	job, err := service.accessor.GetByID(uint(in.Id))
	if err != nil {
		return nil, err
	}

	if job.ID > 0 {
		if job.Status != in.Status && job.Status != models.JobStatusStacked && in.Status == models.JobStatusStacked {
			service.stack <- converters.JobModelToProto(job)
			if status := <-service.stacked; status.Success {
				job.Status = models.JobStatusStacked
				log.Println("Stacked, saving new status")
				if err = service.accessor.Save(job); err != nil {
					return nil, err
				}
			} else {
				log.Println("Error stacking job")
				return nil, status.Err
			}
		} else if job.Status != in.Status {
			job.Status = in.Status
			if err = service.accessor.Save(job); err != nil {
				return nil, err
			}
		}

		return converters.JobModelToProto(job), err
	}

	return nil, errors.New("Job cannot be found")
}
