package main

import (
	"authentication/data"
	"authentication/database"
	"authentication/utils"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := utils.GetEnv("APP_PORT", "8081")

	log.Printf("Starting authentication service on port %s\n", port)

	// Connect database
	conn := database.New()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// Setup config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// Define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// Start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
