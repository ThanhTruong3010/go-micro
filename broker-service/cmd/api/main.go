package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

type Config struct{}

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	app := Config{}

	port := GetEnv("BROKER_PORT", "8080")

	log.Printf("Starting broker service on port %s\n", port)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
