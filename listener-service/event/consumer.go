package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

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

func (c *Consumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
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

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for messages [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		//log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		//authenticate
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, err := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}
