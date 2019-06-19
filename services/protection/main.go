package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"git.ve.home/nicolasc/linotte/services/protection/configuration"
	"git.ve.home/nicolasc/linotte/services/protection/models"
	proto_protection "git.ve.home/nicolasc/linotte/services/protection/proto"
	"git.ve.home/nicolasc/linotte/services/protection/services"
)

var database *models.Connection

type server struct {
	protectionService *services.ProtectionService
}

// ProtectionsByTaxonIDs
func (s *server) ProtectionsByTaxonIDs(ctx context.Context, in *proto_protection.ProtectionRequest) (*proto_protection.ProtectionsReply, error) {
	log.Print("ProtectionsByTaxonIDs")
	return nil, nil
}

func initialize() (*server, error) {
	var (
		instance *server
		err      error
	)

	if database, err = models.Connect(config.Get()); err == nil {
		if err = models.MigrateAll(); err == nil {
			instance = &server{
				protectionService: services.NewProtectionService(models.NewProtectionAccessor(database.DB)),
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
	proto_protection.RegisterProtectionServer(s, instance)
	reflection.Register(s)

	log.Println("Starting Protection service")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
