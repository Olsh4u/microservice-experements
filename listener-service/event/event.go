package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",      //type
		true,         //durability
		false,        //auto deleted
		false,        // internal
		false,        // no wait
		nil,          // specific args
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, //durability
		false, //auto delete
		true,  //exclusive
		false, //no wait
		nil,   // specific args
	)
}
