package main

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"git.ve.home/nicolasc/linotte/libs/helpers"
	"git.ve.home/nicolasc/linotte/services/api/configuration"
	"git.ve.home/nicolasc/linotte/services/api/controllers"

	proto_job "git.ve.home/nicolasc/linotte/services/job/proto"
	proto_search "git.ve.home/nicolasc/linotte/services/search/proto"
	proto_taxref "git.ve.home/nicolasc/linotte/services/taxref/proto"
	proto_user "git.ve.home/nicolasc/linotte/services/user/proto"
)

type Api struct {
	Endpoint                 string
	authenticationController *controllers.AuthenticationController
	taxrefController         *controllers.TaxrefController
	jobController            *controllers.JobController
	jobReportController      *controllers.JobReportController
	jobResultController      *controllers.JobResultController
	user                     proto_user.UserClient
}

var instance *Api

func Initialize() (*Api, error) {
	var err error
	settings := config.Get()

	taxrefConnection, _ := grpc.Dial(settings.TaxrefServiceEndpoint, grpc.WithInsecure())
	searchConnection, _ := grpc.Dial(settings.SearchServiceEndpoint, grpc.WithInsecure())
	userConnection, _ := grpc.Dial(settings.UserServiceEndpoint, grpc.WithInsecure())
	jobConnection, _ := grpc.Dial(settings.JobServiceEndpoint, grpc.WithInsecure())

	job := proto_job.NewJobClient(jobConnection)
	user := proto_user.NewUserClient(userConnection)
	taxref := proto_taxref.NewTaxrefClient(taxrefConnection)
	search := proto_search.NewSearchClient(searchConnection)

	if instance == nil {
		instance = &Api{
			authenticationController: controllers.NewAuthenticationController(user),
			taxrefController: controllers.NewTaxrefController(
				taxref,
				search,
			),
			jobController:       controllers.NewJobController(job),
			jobReportController: controllers.NewJobReportController(job),
			jobResultController: controllers.NewJobResultController(
				job,
				taxref,
			),
			user: user,
		}
	}

	return instance, err
}

// Start starts the API
func (api *Api) Start(endpoint string) error {
	router := mux.NewRouter()

	router.HandleFunc("/login", api.authenticationController.Login).Methods(http.MethodPost)

	secure := router.PathPrefix("/secure").Subrouter()

	secure.HandleFunc("/job", api.jobController.All).Methods(http.MethodGet)
	secure.HandleFunc("/job/{id}", api.jobController.ByID).Methods(http.MethodGet)
	secure.HandleFunc("/job", api.jobController.UpdateJob).Methods(http.MethodPatch)
	secure.HandleFunc("/job-reports/{id}", api.jobReportController.ByJobID).Methods(http.MethodGet)
	secure.HandleFunc("/job-report/{id}", api.jobReportController.ByID).Methods(http.MethodGet)
	secure.HandleFunc("/job-report/statistics/{id}", api.jobReportController.StatisticsForID).Methods(http.MethodGet)

	secure.HandleFunc("/job-results/{report-id}/{state}", api.jobResultController.ResultsByState).Methods(http.MethodGet)
	secure.HandleFunc("/job-results", api.jobResultController.ResolveConflict).Methods(http.MethodPatch)

	secure.HandleFunc("/taxref/rank", api.taxrefController.AllRanks).Methods(http.MethodGet)
	secure.HandleFunc("/taxref/taxon/{id}", api.taxrefController.ByID).Methods(http.MethodGet)
	secure.HandleFunc("/taxref/taxon/all/{id}", api.taxrefController.ReferenceAndSynonymsForID).Methods(http.MethodGet)
	secure.HandleFunc("/taxref/taxon/classification/{id}/{depth}", api.taxrefController.ClassificationForID).Methods(http.MethodGet)
	secure.HandleFunc("/taxref/taxon/search", api.taxrefController.PaginatedSearch).Methods(http.MethodPost)

	global := http.NewServeMux()
	global.Handle("/", router)
	global.Handle("/secure/", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			reply, err := api.user.CheckAuthentication(
				context.Background(),
				&proto_user.CheckRequest{
					Token: r.Header.Get("Authorization"),
				},
			)

			if err == nil && reply.Success {
				next(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				helpers.WriteError(w, "Invalid token", err)
			}
		}),
		negroni.Wrap(router),
	))

	api.Endpoint = endpoint
	log.Printf("API start (%s)", endpoint)

	n := negroni.Classic()
	n.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		Debug:          false,
	}))
	n.UseHandler(global)

	return http.ListenAndServe(endpoint, n)
}

func main() {
	api, err := Initialize()
	helpers.HandleError(err)
	api.Start(config.Get().ApiServiceEndpoint)
}
