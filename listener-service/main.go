package main

import (
	"fmt"
	"listener-service/event"
	"listener-service/utils"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	if err := consumer.Listen([]string{"auth.INFO", "log.INFO", "log.WARNING", "log.ERROR"}); err != nil {
		panic(err)
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
