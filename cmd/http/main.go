package main

import (
	"log"
	"os"

	"github.com/DrIhor/test_task/internal/service/transport/httpServer"
)

func init() {
	os.Setenv("storage-type", "in-memory")
	os.Setenv("server-port", "8080")
	os.Setenv("storage-host", "")

}

func main() {
	// read config
	server := httpServer.New()
	if err := server.ServerAddrConfig(); err != nil {
		log.Fatal("Can`t get config of server: ", err)
	}
	if err := server.ConfigStorage(); err != nil {
		log.Fatal("Can`t config storage: ", err)
	}

	server.GetRouters()
	if err := server.Start(); err != nil {
		log.Fatal("Problems with server run: ", err)
	}
}
