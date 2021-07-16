package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/DrIhor/test_task/internal/routes"
)

func main() {
	// read config
	errLoad := godotenv.Load("../../config/.env")
	if errLoad != nil {
		log.Fatal("Error loading config")
	}

	port := os.Getenv("SERVER_PORT")

	router := routes.Handler()
	fmt.Println("Start HTTP server")
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), routes.Middleware(router)); err != nil {
		log.Fatal(err)
	}
}
