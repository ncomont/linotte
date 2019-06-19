package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"git.ve.home/nicolasc/linotte/libs/helpers"
	"github.com/gorilla/mux"

	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
)

type JobController struct {
	job proto_job.JobClient
}

func NewJobController(job proto_job.JobClient) *JobController {
	return &JobController{job}
}

func (controller *JobController) ByID(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		id    int
		reply *proto_job.JobReply
	)

	if id, err = strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		log.Fatal("Illegal id for ReportsByJobID")
	}

	if reply, err = controller.job.JobByID(
		context.Background(),
		&proto_job.JobRequest{Id: uint32(id)},
	); err != nil {
		helpers.WriteError(w, "Error retrieving job by id", err)
	} else {
		json.NewEncoder(w).Encode(reply)
	}
}

func (controller *JobController) All(w http.ResponseWriter, r *http.Request) {
	if reply, err := controller.job.AllJobs(
		context.Background(),
		&proto_job.JobsRequest{},
	); err != nil {
		helpers.WriteError(w, "Error retrieving jobs", err)
	} else {
		json.NewEncoder(w).Encode(reply.Jobs)
	}
}

// UpdateJob updates the given job's status and put it in the stack if necessary
func (controller *JobController) UpdateJob(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		reply  *proto_job.JobReply
		params proto_job.JobRequest
		buffer []byte
	)

	if buffer, err = ioutil.ReadAll(r.Body); err != nil {
		helpers.WriteError(w, "Error parsing job request", err)
		return
	}
	json.Unmarshal(buffer, &params)

	if reply, err = controller.job.UpdateJob(context.Background(), &params); err != nil {
		helpers.WriteError(w, "Error updating job", err)
	} else {
		json.NewEncoder(w).Encode(reply)
	}
}
