package main

import (
	"log"
	"os"

	"github.com/DrIhor/test_task/internal/service/transport/gRPC"
)

func init() {
	os.Setenv("STORAGE_TYPE", "grpc")
	os.Setenv("GRCP_PORT", "8080")
	os.Setenv("GRCP_HOST", "")

	os.Setenv("GRCP_ADDR", ":8081")

	os.Setenv("STORAGE", "mongo")
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
}

func main() {

	// read config
	server := gRPC.New()
	if err := server.ServerAddrConfig(); err != nil {
		log.Fatal("Can`t get config of server: ", err)
	}
	if err := server.ConfigStorage(); err != nil {
		log.Fatal("Can`t config storage: ", err)
	}

	if err := server.Start(); err != nil {
		log.Fatal("Problems with server run: ", err)
	}
}