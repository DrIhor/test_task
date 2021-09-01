package gRPC

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	configs "github.com/DrIhor/test_task/internal/models/server"
	"google.golang.org/grpc"

	er "github.com/DrIhor/test_task/internal/errors"

	"github.com/DrIhor/test_task/internal/storage/elk"
	mg "github.com/DrIhor/test_task/internal/storage/mongo"
	"github.com/DrIhor/test_task/internal/storage/redis"
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
		return er.WrongPort
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
		pg, err := postgres.New()
		if err != nil {
			return err
		}

		s.storage = pg
		fmt.Println("Start Postgres")
		return nil

	case "mongo":
		mong, err := mg.New()
		if err != nil {
			return err
		}

		s.storage = mong
		fmt.Println("Start Mongo")
		return nil

	case "redis":
		rd, err := redis.New()
		if err != nil {
			return err
		}

		stor := rd
		s.storage = stor
		fmt.Println("Start Redis")
		return nil

	case "elk":
		el, err := elk.New()
		if err != nil {
			return err
		}

		s.storage = el
		fmt.Println("Start ELK")
		return nil
	}

	return er.NoSuchStorage
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", os.Getenv("GRCP_ADDR"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	sr := grpc.NewServer(opts...)
	pb.RegisterItemStorageServer(sr, s)

	// shutdown
	go func() {
		if err := sr.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	fmt.Println("GRCP is running")

	<-ctx.Done()
	log.Println("Server stopped")

	sr.GracefulStop()

	log.Printf("Server exited properly")

	return nil
}
