package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"git.ve.home/nicolasc/linotte/services/user/configuration"
	"git.ve.home/nicolasc/linotte/services/user/models"
	proto_user "git.ve.home/nicolasc/linotte/services/user/proto"
	"git.ve.home/nicolasc/linotte/services/user/services"
)

var database *models.Connection

type server struct {
	userService *services.UserService
}

func (s *server) Login(ctx context.Context, in *proto_user.LoginRequest) (*proto_user.UserReply, error) {
	log.Print("Login")
	return s.userService.Login(in)
}

func (s *server) Logout(ctx context.Context, in *proto_user.LogoutRequest) (*proto_user.UserReply, error) {
	log.Print("Logout")
	return s.userService.Logout(in)
}

func (s *server) CheckAuthentication(ctx context.Context, in *proto_user.CheckRequest) (*proto_user.CheckReply, error) {
	log.Print("CheckAuthentication")
	return s.userService.CheckAuthentication(in)
}

func initialize() (*server, error) {
	var (
		instance *server
		err      error
	)

	if database, err = models.Connect(config.Get()); err == nil {
		if err = models.MigrateAll(); err == nil {
			instance = &server{
				userService: services.NewUserService(models.NewUserAccessor(database.DB)),
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
	proto_user.RegisterUserServer(s, instance)
	reflection.Register(s)

	log.Println("Starting User service")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
