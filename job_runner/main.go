package main

import (
	"log"

	"fmt"

	"git.ve.home/nicolasc/linotte/job_runner/configuration"
	"git.ve.home/nicolasc/linotte/job_runner/runner"
	"git.ve.home/nicolasc/linotte/libs/helpers"
	"git.ve.home/nicolasc/linotte/libs/qw"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
	"github.com/golang/protobuf/proto"
)

var (
	job     *proto_job.JobReply
	results *qw.Q
	status  *qw.Q
	tasks   *qw.Q
)

func main() {
	settings := config.Get()
	channels := &runner.Channels{
		Progression: make(chan uint32),
		Results:     make(chan []*proto_job.ResultReply),
		Errors:      make(chan error),
	}

	tasks = qw.Create(settings.RabbitEndpoint, settings.RabbitTaskQueueID)
	status = qw.Create(settings.RabbitEndpoint, settings.RabbitStatusQueueID)
	results = qw.Create(settings.RabbitEndpoint, settings.RabbitResultQueueID)

	defer tasks.Close()
	defer status.Close()
	defer results.Close()

	// Wait for queues to be ready
	<-tasks.Start()
	<-status.Start()
	<-results.Start()

	go func() {
		msgs, err := tasks.Consume()
		helpers.HandleError(err)

		for d := range msgs {
			job = &proto_job.JobReply{}
			err := proto.Unmarshal(d.Body, job)
			if err != nil {
				helpers.HandleError(err)
			}

			log.Printf("Received a new job of type: %s. Starting runner.\n", job.Type)
			go handleRunnerMessages(channels, job.Id)
			runner.Get(job.Type).Configure(channels).Run(job)
			d.Ack(false)
			log.Println("ACK sent. Runner free.")
		}
	}()

	log.Printf("Waiting for new job. CTRL+C to exit.")
	<-make(chan bool)
}

func handleRunnerMessages(channels *runner.Channels, id uint32) {
	var last uint32
	for {
		select {
		case r := <-channels.Results:
			report := &proto_job.ReportRequest{
				JobId:   id,
				Results: r,
			}

			report.Status = "PASSED"
			for _, result := range report.Results {
				// TODO: CC status !
				if result.State != "FOUND" {
					report.Status = "WARNING"
					break
				}
			}
			body, _ := proto.Marshal(report)

			results.Write(body)
			fmt.Printf("Results published for job: %d\n", id)
			return
		case e := <-channels.Errors:
			body, _ := proto.Marshal(&qw.Status{
				Error: e.Error(),
				JobId: id,
			})
			status.Write(body)
			fmt.Printf("Error occured on job %d: %v!\n", id, e)
			return
		case p := <-channels.Progression:
			if p > 0 && p != last {
				body, _ := proto.Marshal(&qw.Status{
					JobId:    id,
					Progress: p,
				})
				status.Write(body)
				last = p
			}
		}
	}
}
