package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"git.ve.home/nicolasc/linotte/libs/helpers"
	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
	proto_taxref "git.ve.home/nicolasc/linotte/services/taxref/proto"
)

type JobResultController struct {
	job    proto_job.JobClient
	taxref proto_taxref.TaxrefClient
}

type reply struct {
	Id         uint32                   `json:"id"`
	SearchData string                   `json:"search_data"`
	Value      string                   `json:"value"`
	ReportId   uint32                   `json:"report_id"`
	TaxonId    uint32                   `json:"taxon_id"`
	State      string                   `json:"state"`
	Taxon      *proto_taxref.TaxonReply `json:"taxon"`
}

func NewJobResultController(job proto_job.JobClient, taxref proto_taxref.TaxrefClient) *JobResultController {
	return &JobResultController{job, taxref}
}

func (controller *JobResultController) ResultsByState(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		id       int
		results  *proto_job.ResultsReply
		taxons   *proto_taxref.TaxonsReply
		ids      []uint32
		computed []*reply
	)

	if id, err = strconv.Atoi(mux.Vars(r)["report-id"]); err != nil {
		helpers.WriteError(w, "Illegal report id", err)
		return
	}

	results, err = controller.job.ResultsByStateAndReportID(
		context.Background(),
		&proto_job.ResultRequest{
			ReportId: uint32(id),
			State:    mux.Vars(r)["state"],
		},
	)

	if err != nil {
		helpers.WriteError(w, "Error retrieving results by state", err)
		return
	}

	for _, result := range results.Results {
		if result.TaxonId > 0 {
			ids = append(ids, result.TaxonId)
		}
	}

	taxons, err = controller.taxref.TaxonsByIDs(
		context.Background(),
		&proto_taxref.TaxonsRequest{Ids: ids},
	)

	for _, result := range results.Results {
		var t *proto_taxref.TaxonReply
		for _, taxon := range taxons.Taxons {
			if result.TaxonId == taxon.Id {
				t = taxon
			}
		}

		computed = append(computed, &reply{
			Id:         result.Id,
			SearchData: result.SearchData,
			Value:      result.Value,
			ReportId:   result.ReportId,
			TaxonId:    result.TaxonId,
			State:      result.State,
			Taxon:      t,
		})
	}

	if err != nil {
		helpers.WriteError(w, "Error retrieving results by state", err)
	} else {
		json.NewEncoder(w).Encode(computed)
	}
}

func (controller *JobResultController) ResolveConflict(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		reply *proto_job.ResultReply
	)

	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		helpers.WriteError(w, "Invalid parameters", err)
		return
	}

	parameters := &struct {
		ConflictID uint `json:"conflictId"`
		TaxonID    uint `json:"taxonId"`
	}{}

	if err = json.Unmarshal(p, &parameters); err != nil {
		helpers.WriteError(w, "Error parsing parameters in ResolveConflict", err)
		return
	}

	taxon, err := controller.taxref.ReferenceByID(
		context.Background(),
		&proto_taxref.TaxonRequest{
			Id: uint32(parameters.TaxonID),
		},
	)

	if err != nil {
		helpers.WriteError(w, "Error getting reference taxon", err)
		return
	}

	if taxon.Id == 0 {
		helpers.WriteError(w, "Taxon not found in ResolveConflict", err)
		return
	}

	reply, err = controller.job.ResolveConflict(
		context.Background(),
		&proto_job.ResultRequest{
			Id:      uint32(parameters.ConflictID),
			TaxonId: uint32(taxon.Id),
		},
	)

	if err != nil {
		helpers.WriteError(w, "Error during conflict resolution", err)
	} else {
		json.NewEncoder(w).Encode(reply)
	}
}
