package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
)

type JobReportController struct {
	job proto_job.JobClient
}

func NewJobReportController(job proto_job.JobClient) *JobReportController {
	return &JobReportController{job}
}

func (controller *JobReportController) ByID(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		id    int
		reply *proto_job.ReportReply
	)

	if id, err = strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		log.Fatal("Illegal id for ReportsByJobID")
	}

	reply, _ = controller.job.ReportByID(context.Background(), &proto_job.ReportRequest{Id: uint32(id)})
	json.NewEncoder(w).Encode(reply)
}

func (controller *JobReportController) ByJobID(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		id    int
		reply *proto_job.ReportsReply
	)

	if id, err = strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		log.Fatal("Illegal id for ReportsByJobID")
	}

	reply, _ = controller.job.ReportsByJobID(context.Background(), &proto_job.JobRequest{Id: uint32(id)})
	json.NewEncoder(w).Encode(reply)
}

func (controller *JobReportController) StatisticsForID(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		id    int
		reply *proto_job.ReportStatisticsReply
	)

	if id, err = strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		log.Fatal("Illegal id for StatisticsForID")
	}

	reply, _ = controller.job.ReportStatisticsByReportID(context.Background(), &proto_job.ReportRequest{Id: uint32(id)})
	json.NewEncoder(w).Encode(reply)
}
