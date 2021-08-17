package graphQL

import (
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
	"github.com/DrIhor/test_task/internal/storage/memory"
	"github.com/DrIhor/test_task/internal/storage/postgres"
	"github.com/gorilla/mux"
)

type Server struct {
	config  *configs.Config
	router  *mux.Router
	storage items.ItemStorageServices
}

func New() *Server {
	return &Server{
		router: mux.NewRouter(),
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

	}

	return errors.New("No such storage")
}

func (s *Server) getHttpAddress() string {
	return fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
}

func (s *Server) Start() error {
	fmt.Println("Server is running on " + s.config.Host)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(s.storage)}))
	http.Handle("/", srv)

	server := &http.Server{
		Addr:         s.getHttpAddress(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return server.ListenAndServe()
}
