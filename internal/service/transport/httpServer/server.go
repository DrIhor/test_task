package httpServer

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/DrIhor/test_task/internal/models/iteams"
	configs "github.com/DrIhor/test_task/internal/models/server"
	routes "github.com/DrIhor/test_task/internal/routes"

	iteamsStor "github.com/DrIhor/test_task/internal/storage/memory"

	"github.com/gorilla/mux"
)

type Server struct {
	config  *configs.Config
	router  *mux.Router
	storage iteams.IteamStorageServices
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
		stor := iteamsStor.New()
		s.storage = stor
		return nil
	}
	return errors.New("No such storage")
}

func (s *Server) GetRouters() {
	routes.HandlerItems(s.router, s.storage)
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

	return nil
}
