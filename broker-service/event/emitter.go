package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Emiiter struct {
	conn *amqp.Connection
}

func (e *Emiiter) setup() error {
	channel, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (e *Emiiter) Push(event string, severity string) error {
	channel, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEmitter(conn *amqp.Connection) (Emiiter, error) {
	emitter := Emiiter{
		conn: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emiiter{}, err
	}

	return emitter, nil
}
