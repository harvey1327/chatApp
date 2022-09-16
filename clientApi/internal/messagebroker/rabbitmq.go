package messagebroker

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker interface {
	CloseConnection()
	DeclareQueue(queueName string) amqp.Queue
	Publish(message message) error
}

type rabbitImpl struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitMQ() MessageBroker {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &rabbitImpl{
		connection: connection,
		channel:    channel,
	}
}

func (rmq *rabbitImpl) Publish(message message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return rmq.channel.Publish("", message.queueName, false, false, amqp.Publishing{ContentType: message.contentType, Body: bytes})
}

func (rmq *rabbitImpl) DeclareQueue(queueName string) amqp.Queue {
	queue, err := rmq.channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	return queue
}

func (rmq *rabbitImpl) CloseConnection() {
	rmq.channel.Close()
	rmq.connection.Close()
}
