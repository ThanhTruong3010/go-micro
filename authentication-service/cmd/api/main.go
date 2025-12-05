package main

import (
	"authentication/data"
	"authentication/database"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	// "github.com/joho/godotenv"
)

const WEB_PORT = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	// Load .env file if it exists
	// if err := godotenv.Load(); err != nil {
	// 	log.Println("No .env file found, using environment variables")
	// }

	log.Println("Starting authention service")

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
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	// Start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
