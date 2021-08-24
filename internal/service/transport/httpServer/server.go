package httpServer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/DrIhor/test_task/internal/models/items"
	configs "github.com/DrIhor/test_task/internal/models/server"
	"github.com/DrIhor/test_task/internal/service/transport/graphQL/graph"
	"github.com/DrIhor/test_task/internal/service/transport/graphQL/graph/generated"
	routes "github.com/DrIhor/test_task/internal/service/transport/httpServer/routes"

	elk "github.com/DrIhor/test_task/internal/storage/elk"

	"github.com/DrIhor/test_task/internal/storage/memory"
	mg "github.com/DrIhor/test_task/internal/storage/mongo"
	"github.com/DrIhor/test_task/internal/storage/postgres"
	"github.com/DrIhor/test_task/internal/storage/redis"

	"github.com/gorilla/mux"
)

type Server struct {
	ctx     context.Context
	config  *configs.Config
	router  *mux.Router
	storage items.ItemStorageServices
}

func New() *Server {
	return &Server{
		router: mux.NewRouter(),
		ctx:    context.Background(),
	}
}

func (s *Server) ServerAddrConfig() error {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return errors.New("Wrong port")
	}

	s.config = &configs.Config{
		Host: os.Getenv("SERVER_HOST"),
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

	case "redis":
		stor := redis.New()
		s.storage = stor
		fmt.Println("Start Redis")
		return nil

	case "elk":
		stor := elk.New()
		s.storage = stor
		fmt.Println("Start ELK")
		return nil
	}
	return errors.New("No such storage")
}

/**
 * all routes of http server
 * also contain connection to graphql server
 */

func (s *Server) GetRouters() {
	itemsHandler := routes.New(s.ctx, s.router, s.storage)
	itemsHandler.HandlerItems()
}

// work with graphql queries
func (s *Server) AddGraohQLRoutes() {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(s.storage)}))
	s.router.Handle("/query/", srv)
}

func (s *Server) getHttpAddress() string {
	return fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
}

func (s *Server) Start() error {
	fmt.Println("Server is running on " + s.config.Host)
	err := http.ListenAndServe(s.getHttpAddress(), routes.Middleware(s.router))
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:         s.getHttpAddress(),
		Handler:      routes.Middleware(s.router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return server.ListenAndServe()
}
