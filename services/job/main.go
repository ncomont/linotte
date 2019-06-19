package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"git.ve.home/nicolasc/linotte/libs/qw"
	"git.ve.home/nicolasc/linotte/services/job/configuration"
	"git.ve.home/nicolasc/linotte/services/job/models"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
	"git.ve.home/nicolasc/linotte/services/job/services"
)

var (
	database *models.Connection
	stack    chan *proto_job.JobReply
	stacked  chan services.TaskRequestStatus
	instance *server
	tasks    *qw.Q
	results  *qw.Q
	status   *qw.Q
)

type server struct {
	jobService    *services.JobService
	reportService *services.ReportService
	resultService *services.ResultService
}

// AllJobs returns every jobs
func (s *server) AllJobs(ctx context.Context, in *proto_job.JobsRequest) (*proto_job.JobsReply, error) {
	log.Print("AllJobs")
	return s.jobService.AllJobs(in)
}

// JobByID returns a single job
func (s *server) JobByID(ctx context.Context, in *proto_job.JobRequest) (*proto_job.JobReply, error) {
	log.Print("JobByID")
	return s.jobService.JobByID(in)
}

// ReportByID returns a single report
func (s *server) ReportByID(ctx context.Context, in *proto_job.ReportRequest) (*proto_job.ReportReply, error) {
	log.Print("ReportsByID")
	return s.reportService.ReportByID(in)
}

// ReportsByJobID returns every reports attached to the given job
func (s *server) ReportsByJobID(ctx context.Context, in *proto_job.JobRequest) (*proto_job.ReportsReply, error) {
	log.Print("ReportsByJobID")
	return s.reportService.ReportsByJobID(in)
}

// ReportsByJobID returns every reports attached to the given job
func (s *server) ReportStatisticsByReportID(ctx context.Context, in *proto_job.ReportRequest) (*proto_job.ReportStatisticsReply, error) {
	log.Print("ReportStatisticsByReportID")
	return s.reportService.ReportStatisticsByReportID(in)
}

// SaveReport saves the given report in database
func (s *server) SaveReport(ctx context.Context, in *proto_job.ReportRequest) (*proto_job.ReportReply, error) {
	log.Print("SaveReport")
	return s.reportService.SaveReport(in)
}

// ResolveConflict sets the "RESOLVED" state to the given result
func (s *server) ResolveConflict(ctx context.Context, in *proto_job.ResultRequest) (*proto_job.ResultReply, error) {
	log.Print("ResolveConflict")
	return s.resultService.ResolveConflict(in)
}

// ResultsByState returns every results filtered by state
func (s *server) ResultsByState(ctx context.Context, in *proto_job.ResultRequest) (*proto_job.ResultsReply, error) {
	log.Print("ResultsByState")
	return s.resultService.ResultsByState(in)
}

// ResultsByStateAndReportID returns every results attached to the given report filtered by state
func (s *server) ResultsByStateAndReportID(ctx context.Context, in *proto_job.ResultRequest) (*proto_job.ResultsReply, error) {
	log.Print("ResultsByStateAndReportID")
	return s.resultService.ResultsByStateAndReportID(in)
}

// UpdateJob update the job's status and puts it in the stack if necessary
func (s *server) UpdateJob(ctx context.Context, in *proto_job.JobRequest) (*proto_job.JobReply, error) {
	log.Print("UpdateJob")
	if in.Status != "" {
		return s.jobService.UpdateJobStatus(in)
	}
	return &proto_job.JobReply{}, errors.New("UpdateJob can only update job's status")
}

func initialize() (*server, error) {
	var err error

	stack = make(chan *proto_job.JobReply)
	stacked = make(chan services.TaskRequestStatus)

	if database, err = models.Connect(config.Get()); err == nil {
		if err = models.MigrateAll(); err == nil {
			instance = &server{
				jobService:    services.NewJobService(models.NewJobAccessor(database.DB), stack, stacked),
				reportService: services.NewReportService(models.NewJobReportAccessor(database.DB)),
				resultService: services.NewResultService(models.NewJobResultAccessor(database.DB)),
			}
		}
	}

	return instance, err
}

func main() {
	var err error

	settings := config.Get()
	instance, err = initialize()
	defer database.Close()
	if err != nil {
		log.Fatalf("failed to initialize: %v", err)
	}

	tasks = qw.Create(settings.RabbitEndpoint, settings.RabbitTaskQueueID)
	defer tasks.Close()
	go tasksListener()

	results = qw.Create(settings.RabbitEndpoint, settings.RabbitResultQueueID)
	defer results.Close()
	go resultsListener()

	status = qw.Create(settings.RabbitEndpoint, settings.RabbitStatusQueueID)
	defer status.Close()
	go statusListener()

	lis, err := net.Listen("tcp", config.Get().ServiceEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	proto_job.RegisterJobServer(server, instance)
	reflection.Register(server)

	log.Println("Starting Job service")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func tasksListener() {
	<-tasks.Start() // Wait for tasks queue to be ready

	for {
		select {
		case job := <-stack:
			log.Println("Adding new job to stack")
			body, _ := proto.Marshal(job)
			err := tasks.Write(body)

			stacked <- services.TaskRequestStatus{
				Success: err == nil,
				Err:     err,
			}
		}
	}
}

func resultsListener() {
	var err error
	<-results.Start() // Wait for results queue to be ready

	msgs, err := results.Consume()
	if err != nil {
		log.Fatalf("Unable to consume results queue: %v", err)
	}

	for msg := range msgs {
		r := &proto_job.ReportRequest{}

		if err = proto.Unmarshal(msg.Body, r); err != nil {
			log.Fatalf("failed to parse report: %v", err)
		}
		if _, err = instance.reportService.SaveReport(r); err != nil {
			log.Printf("Error saving report: %v\n", err)
		}
		if _, err = instance.jobService.UpdateJobStatus(&proto_job.JobRequest{
			Id:     r.JobId,
			Status: "IDLE",
		}); err != nil {
			log.Printf("Error saving job status: %v\n", err)
		}

		msg.Ack(false)
	}
}

func statusListener() {
	<-status.Start() // Wait for status queue to be ready

	msgs, err := status.Consume()
	if err != nil {
		log.Fatalf("Unable to consume results queue: %v", err)
	}

	for msg := range msgs {
		s := &qw.Status{}

		if err = proto.Unmarshal(msg.Body, s); err != nil {
			log.Fatalf("failed to parse status: %v", err)
		}

		if s.Error != "" {
			log.Println("error occured ! from queue")
			if _, err = instance.reportService.SaveReport(&proto_job.ReportRequest{
				JobId:   s.JobId,
				Status:  "ERROR",
				Message: s.Error,
			}); err != nil {
				log.Printf("Error saving report: %v\n", err)
			}

			if _, err = instance.jobService.UpdateJobStatus(&proto_job.JobRequest{
				Id:     s.JobId,
				Status: "IDLE",
			}); err != nil {
				log.Printf("Error saving job status: %v\n", err)
			}

			log.Println("Error report saved")
		} else {
			log.Printf("Progression received: %d\n", s.Progress)

			if _, err = instance.jobService.UpdateJobStatus(&proto_job.JobRequest{
				Id:     s.JobId,
				Status: fmt.Sprintf("%d%%", s.Progress),
			}); err != nil {
				log.Printf("Error saving job status: %v\n", err)
			}
		}

		msg.Ack(false)
	}
}
