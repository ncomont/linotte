package main

import (
	"fmt"
	"log"
	"net"

	elastic "gopkg.in/olivere/elastic.v5"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"git.ve.home/nicolasc/linotte/services/search/configuration"
	proto_search "git.ve.home/nicolasc/linotte/services/search/proto"
	"git.ve.home/nicolasc/linotte/services/search/services"
)

type server struct {
	searchService *services.SearchService
}

// Search returns every taxons ids corresponding to the search request
func (s *server) Search(ctx context.Context, in *proto_search.SearchRequest) (*proto_search.SearchReply, error) {
	log.Print("Search")
	return s.searchService.Search(in)
}

func initialize() (*server, error) {
	var (
		instance *server
		err      error
		client   *elastic.Client
	)

	settings := config.Get()

	if instance == nil {
		url := fmt.Sprintf("http://%s:%d", settings.ElasticHost, settings.ElasticPort)
		if client, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(true)); err == nil {
			instance = &server{
				searchService: services.NewSearchService(client, settings.ElasticIndex),
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

	lis, err := net.Listen("tcp", config.Get().ServiceEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto_search.RegisterSearchServer(s, instance)
	reflection.Register(s)

	log.Println("Starting Search service")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
