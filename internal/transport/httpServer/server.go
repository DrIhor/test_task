package httpServer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/DrIhor/test_task/internal/models/items"
	configs "github.com/DrIhor/test_task/internal/models/server"
	"github.com/DrIhor/test_task/internal/storage/elk"
	"github.com/DrIhor/test_task/internal/storage/memory"
	"github.com/DrIhor/test_task/internal/storage/postgres"
	"github.com/DrIhor/test_task/internal/storage/redis"
	"github.com/DrIhor/test_task/internal/transport/graphQL/graph"
	"github.com/DrIhor/test_task/internal/transport/graphQL/graph/generated"
	"github.com/DrIhor/test_task/internal/transport/httpServer/middleware"
	"github.com/DrIhor/test_task/internal/transport/httpServer/routes"

	"github.com/gorilla/mux"

	er "github.com/DrIhor/test_task/internal/errors"

	mg "github.com/DrIhor/test_task/internal/storage/mongo"
)

type Server struct {
	ctx     context.Context
	config  *configs.Config
	router  *mux.Router
	storage items.ItemStorageServices
}

func New(ctx context.Context) *Server {
	return &Server{
		router: mux.NewRouter(),
		ctx:    ctx,
	}
}

func (s *Server) ServerAddrConfig() error {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		return er.WrongPort
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

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:         s.getHttpAddress(),
		Handler:      middleware.JwtSessionCheck(middleware.JsonRespHeaders(s.router)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	cancelTimeout, errTimeout := strconv.Atoi(os.Getenv("Server_Cancel_Timeout"))
	if errTimeout != nil {
		return errTimeout
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	fmt.Println("Server is running on: " + s.getHttpAddress())

	<-ctx.Done()
	log.Println("Server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Duration(cancelTimeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutDown); err != nil {
		return err
	}

	log.Printf("Server exited properly")

	return nil
}
