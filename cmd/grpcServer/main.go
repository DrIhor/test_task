package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/DrIhor/test_task/internal/transport/gRPC"
)

func init() {
	os.Setenv("STORAGE_TYPE", "grpc")

	// grpc
	os.Setenv("GRCP_PORT", "8080")
	os.Setenv("GRCP_HOST", "")
	os.Setenv("GRCP_ADDR", ":8081")

	// grpc
	os.Setenv("STORAGE", "elk")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_HOST", "")

	// postgres
	os.Setenv("POSTGRE_HOST", "localhost")
	os.Setenv("POSTGRE_PORT", "5432")
	os.Setenv("POSTGRE_USER", "postgres")
	os.Setenv("POSTGRE_PASS", "postgres")
	os.Setenv("POSTGRE_DB", "postgres")

	//redis
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASS", "")
	os.Setenv("REDIS_DB", "0")

	// mongo
	os.Setenv("MONGO_ADDR", "localhost:27017/?readPreference=primary&ssl=false")

	// elk
	os.Setenv("ELASTIC_ADDR", "http://localhost:9200")
	os.Setenv("ELASTIC_USER", "")
	os.Setenv("ELASTIC_PASSWORD", "")
}

func main() {
	// init data
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// read config
	server := gRPC.New()
	if err := server.ServerAddrConfig(); err != nil {
		log.Fatal("Can`t get config of server: ", err)
	}
	if err := server.ConfigStorage(); err != nil {
		log.Fatal("Can`t config storage: ", err)
	}

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
