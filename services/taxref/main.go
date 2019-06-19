package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"git.ve.home/nicolasc/linotte/services/taxref/configuration"
	"git.ve.home/nicolasc/linotte/services/taxref/models"
	proto_taxref "git.ve.home/nicolasc/linotte/services/taxref/proto"
	"git.ve.home/nicolasc/linotte/services/taxref/services"
)

var database *models.Connection

type server struct {
	taxonService *services.TaxonService
	rankService  *services.RankService
}

func (s *server) AllRanks(ctx context.Context, in *proto_taxref.RanksRequest) (*proto_taxref.RanksReply, error) {
	log.Print("AllRanks")
	return s.rankService.AllRanks()
}

func (s *server) TaxonByID(ctx context.Context, in *proto_taxref.TaxonRequest) (*proto_taxref.TaxonReply, error) {
	log.Print("TaxonByID")
	return s.taxonService.TaxonByID(in)
}

func (s *server) ReferenceByID(ctx context.Context, in *proto_taxref.TaxonRequest) (*proto_taxref.TaxonReply, error) {
	log.Print("ReferenceByID")
	return s.taxonService.ReferenceByID(in)
}

func (s *server) TaxonsByIDs(ctx context.Context, in *proto_taxref.TaxonsRequest) (*proto_taxref.TaxonsReply, error) {
	log.Print("TaxonsByIDs")
	return s.taxonService.TaxonsByIDs(in)
}

// ReferenceAndSynonymsForTaxonID returns the reference and every synonyms of the given taxon
func (s *server) ReferenceAndSynonymsForTaxonID(ctx context.Context, in *proto_taxref.TaxonRequest) (*proto_taxref.TaxonsReply, error) {
	log.Print("ReferenceAndSynonymsForTaxonID")
	return s.taxonService.ReferenceAndSynonymsForTaxonID(in)
}

// TaxonClassificationForID returns the classification of the given taxon
func (s *server) TaxonClassificationForID(ctx context.Context, in *proto_taxref.TaxonClassificationRequest) (*proto_taxref.TaxonReply, error) {
	log.Print("TaxonClassificationForID")
	return s.taxonService.TaxonClassificationForID(in)
}

// ReferenceByVerb returns the reference taxon corresponding to the given verb
func (s *server) ReferenceByVerb(ctx context.Context, in *proto_taxref.TaxonRequest) (*proto_taxref.TaxonReply, error) {
	log.Print("ReferenceByVerb")
	return s.taxonService.ReferenceByVerb(in)
}

func initialize() (*server, error) {
	var (
		instance *server
		err      error
	)

	if database, err = models.Connect(config.Get()); err == nil {
		if err = models.MigrateAll(); err == nil {
			instance = &server{
				rankService:  services.NewRankService(models.NewRankAccessor(database.DB)),
				taxonService: services.NewTaxonService(models.NewTaxonAccessor(database.DB)),
			}
		}
	}

	return instance, err
}

func main() {
	instance, err := initialize()
	if err != nil {
		log.Fatalf("failed to initialize: %v", err)
	}

	defer database.Close()

	lis, err := net.Listen("tcp", config.Get().ServiceEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto_taxref.RegisterTaxrefServer(s, instance)
	reflection.Register(s)

	log.Println("Starting Taxref service")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
