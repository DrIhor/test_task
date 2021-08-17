package gRPC

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	configs "github.com/DrIhor/test_task/internal/models/server"
	"google.golang.org/grpc"

	mg "github.com/DrIhor/test_task/internal/storage/mongo"
	pb "github.com/DrIhor/test_task/pkg/grpc"

	"github.com/DrIhor/test_task/internal/models/items"
	"github.com/DrIhor/test_task/internal/storage/memory"
	"github.com/DrIhor/test_task/internal/storage/postgres"
)

type Server struct {
	config *configs.Config
	pb.UnimplementedItemStorageServer
	storage items.ItemStorageServices
}

func New() *Server {
	return &Server{}
}

func (s *Server) ServerAddrConfig() error {
	port := os.Getenv("GRCP_PORT")
	if port == "" {
		return errors.New("Wrong port")
	}

	s.config = &configs.Config{
		Host: os.Getenv("GRCP_HOST"),
		Port: port,
	}
	return nil
}

func (s *Server) ConfigStorage() error {
	storageType := os.Getenv("STORAGE")
	switch storageType {
	case "in-memory":
		stor := memory.New()
		s.storage = stor
		fmt.Println("Start in-memory DB")
		return nil

	case "postgres":
		stor := postgres.New()
		s.storage = stor
		fmt.Println("Start Postgres")
		return nil

	case "mongo":
		stor := mg.New()
		s.storage = stor
		fmt.Println("Start Mongo")
		return nil
	}
	return errors.New("No such storage")
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", os.Getenv("GRCP_ADDR"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	sr := grpc.NewServer(opts...)
	pb.RegisterItemStorageServer(sr, s)
	if err := sr.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}
