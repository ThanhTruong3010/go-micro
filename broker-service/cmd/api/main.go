package main

import (
	"broker-service/utils"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	port := utils.GetEnv("BROKER_PORT", "8080")

	log.Printf("Starting broker service on port %s\n", port)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		var url string = utils.GetRabbitMQURL()
		c, err := amqp.Dial(url)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Printf("RabbitMQ not ready after 5 attempts: %s", err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second
		log.Printf("backing off for %s", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
