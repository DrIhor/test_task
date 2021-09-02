package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/DrIhor/test_task/internal/transport/httpServer"
)

func init() {
	os.Setenv("Server_Cancel_Timeout", "5")
	os.Setenv("STORAGE_TYPE", "")

	os.Setenv("CHECK_TOKEN", "false")

	// grpc
	os.Setenv("GRCP_PORT", "8080")
	os.Setenv("GRCP_HOST", "")
	os.Setenv("GRCP_ADDR", ":8081")

	// server
	os.Setenv("STORAGE", "postgres")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_HOST", "")

	// postgres
	os.Setenv("POSTGRE_HOST", "localhost")
	os.Setenv("POSTGRE_PORT", "5432")
	os.Setenv("POSTGRE_USER", "postgres")
	os.Setenv("POSTGRE_PASS", "postgres")
	os.Setenv("POSTGRE_DB", "postgres")

	// mongo
	os.Setenv("MONGO_ADDR", "localhost:27017/?readPreference=primary&ssl=false")

	//redis
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASS", "")
	os.Setenv("REDIS_DB", "0")

	// elk
	os.Setenv("ELASTIC_ADDR", "http://localhost:9200")
	os.Setenv("ELASTIC_USER", "")
	os.Setenv("ELASTIC_PASSWORD", "")

}

func main() {

	// init data
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx := context.Background()

	// read config
	server := httpServer.New(ctx)
	if err := server.ServerAddrConfig(); err != nil {
		log.Fatal("Can`t get config of server: ", err)
	}
	if err := server.ConfigStorage(); err != nil {
		log.Fatal("Can`t config storage: ", err)
	}

	server.GetRouters()
	server.AddGraohQLRoutes()

	// create shutdown
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := server.Start(ctx); err != nil {
		log.Fatal("Problems with server run: ", err)
	}
}
