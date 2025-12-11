package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"listener-service/utils"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

var logServiceURL string = utils.GetServiceURL("logger")

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		ch.QueueBind(
			q.Name,
			topic,
			"logs_topic",
			false,
			nil,
		)

		// if err != nil {
		// 	return err
		// }
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for messages [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	case "auth":
		// authenticate
	}

	// we can have as many cases as we want, as long as we write the logic
}

func logEvent(entry Payload) error {
	// create some json we'll send to the logger microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/log", logServiceURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// make sure we get back the correct status code
	if res.StatusCode == http.StatusUnauthorized {
		return err
	}

	return nil
}
