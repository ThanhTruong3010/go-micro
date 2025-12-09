package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Mailer Mail
}

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	app := Config{
		Mailer: createMail(),
	}

	port := GetEnv("MAILER_PORT", "8083")

	log.Printf("Starting mailer service on port %s\n", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func createMail() Mail {
	port, _ := strconv.Atoi(GetEnv("MAIL_PORT", "1025"))
	m := Mail{
		Domain:      GetEnv("MAIL_DOMAIN", "localhost"),
		Host:        GetEnv("MAIL_HOST", "mailhog"),
		Port:        port,
		Username:    GetEnv("MAIL_USERNAME", ""),
		Password:    GetEnv("MAIL_PASSWORD", ""),
		Encryption:  GetEnv("MAIL_ENCRYPTION", "none"),
		FromName:    GetEnv("MAIL_FROM_NAME", "Thanh Truong"),
		FromAddress: GetEnv("MAIL_FROM_ADDRESS", "thanh.truong@example.com"),
	}

	return m
}
