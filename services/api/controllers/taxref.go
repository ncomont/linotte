package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"git.ve.home/nicolasc/linotte/libs/helpers"
	proto_search "git.ve.home/nicolasc/linotte/services/search/proto"
	proto_taxref "git.ve.home/nicolasc/linotte/services/taxref/proto"
)

type TaxrefController struct {
	taxref proto_taxref.TaxrefClient
	search proto_search.SearchClient
}

func NewTaxrefController(taxref proto_taxref.TaxrefClient, search proto_search.SearchClient) *TaxrefController {
	return &TaxrefController{
		taxref: taxref,
		search: search,
	}
}

func (controller *TaxrefController) PaginatedSearch(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		pageNumber int
		limit      int
		params     proto_search.SearchRequest
		reply      *proto_search.SearchReply
		results    *proto_taxref.TaxonsReply
	)

	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &params)

	if len(params.Ids) > 0 {
		results, err = controller.taxref.TaxonsByIDs(
			context.Background(),
			&proto_taxref.TaxonsRequest{Ids: params.Ids},
		)

		if err != nil {
			helpers.WriteError(w, "Error getting taxons by IDs", err)
		} else {
			json.NewEncoder(w).Encode(results.Taxons)
		}
		return
	}

	if reply, err = controller.search.Search(context.Background(), &params); err != nil {
		helpers.WriteError(w, "Error during search", err)
		return
	}

	results, err = controller.taxref.TaxonsByIDs(
		context.Background(),
		&proto_taxref.TaxonsRequest{
			Ids: reply.Results,
		},
	)

	if err != nil {
		helpers.WriteError(w, "Error gettings taxons related to search results", err)
	} else {
		json.NewEncoder(w).Encode(struct {
			PageNumber int                        `json:"pageNumber"`
			PageSize   int                        `json:"pageSize"`
			Total      int64                      `json:"total"`
			Results    []*proto_taxref.TaxonReply `json:"results"`
		}{
			PageNumber: pageNumber,
			PageSize:   limit,
			Total:      reply.Total,
			Results:    results.Taxons,
		})
	}
}

func (controller *TaxrefController) ByID(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		id    int
		reply *proto_taxref.TaxonReply
	)

	if id, err = strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		log.Fatal("Illegal id for ByID")
	}

	reply, _ = controller.taxref.TaxonByID(
		context.Background(),
		&proto_taxref.TaxonRequest{Id: uint32(id)},
	)
	json.NewEncoder(w).Encode(reply)
}

func (controller *TaxrefController) ReferenceAndSynonymsForID(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		id    int
		reply *proto_taxref.TaxonsReply
	)

	if id, err = strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		helpers.WriteError(w, "Illegal id for ReferenceAndSynonymsForID", err)
		return
	}

	if reply, err = controller.taxref.ReferenceAndSynonymsForTaxonID(
		context.Background(),
		&proto_taxref.TaxonRequest{Id: uint32(id)},
	); err != nil {
		helpers.WriteError(w, "Error retrieving ReferenceAndSynonymsForTaxonID", err)
	} else {
		json.NewEncoder(w).Encode(reply.Taxons)
	}
}

func (controller *TaxrefController) ClassificationForID(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		id    int
		depth int
		reply *proto_taxref.TaxonReply
	)

	if id, err = strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		helpers.WriteError(w, "Illegal id for ReferenceAndSynonymsForID", err)
		return
	}
	if depth, err = strconv.Atoi(mux.Vars(r)["depth"]); err != nil {
		helpers.WriteError(w, "Illegal depth for ReferenceAndSynonymsForID", err)
		return
	}

	if reply, err = controller.taxref.TaxonClassificationForID(
		context.Background(),
		&proto_taxref.TaxonClassificationRequest{
			Id:    uint32(id),
			Depth: uint32(depth),
		},
	); err != nil {
		helpers.WriteError(w, "Error retrieving TaxonClassificationForID", err)
	} else {
		json.NewEncoder(w).Encode(reply)
	}
}

func (controller *TaxrefController) AllRanks(w http.ResponseWriter, r *http.Request) {
	if reply, err := controller.taxref.AllRanks(context.Background(), &proto_taxref.RanksRequest{}); err != nil {
		helpers.WriteError(w, "Error retrieving ranks", err)
	} else {
		json.NewEncoder(w).Encode(reply.Ranks)
	}
}
