package main

import (
	"log"
	"os"

	sr "github.com/DrIhor/test_task/internal/service/transport/httpServer"
)

func init() {
	os.Setenv("STORAGE_TYPE", "grpc")
	os.Setenv("GRCP_ADDR", "localhost:8081")

	os.Setenv("STORAGE", "postgres")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_HOST", "")

	// postgres
	os.Setenv("POSTGRE_HOST", "localhost")
	os.Setenv("POSTGRE_PORT", "5432")
	os.Setenv("POSTGRE_USER", "postgres")
	os.Setenv("POSTGRE_PASS", "postgres")
	os.Setenv("POSTGRE_DB", "postgres")

}

func main() {
	// read config
	server := sr.New()
	if err := server.ServerAddrConfig(); err != nil {
		log.Fatal("Can`t get config of server: ", err)
	}
	if err := server.ConfigStorage(); err != nil {
		log.Fatal("Can`t config storage: ", err)
	}

	server.AddGraohQLRoutes() // add only graphql routes

	if err := server.Start(); err != nil {
		log.Fatal("Problems with server run: ", err)
	}
}
