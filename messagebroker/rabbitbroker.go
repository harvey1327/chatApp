package messagebroker

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker interface {
	CloseConnection()
	DeclareQueue(queueName string)
	getChannel() *amqp.Channel
}

type rabbitMessageBroker struct {
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

	return &rabbitMessageBroker{
		connection: connection,
		channel:    channel,
	}
}

func (rmq *rabbitMessageBroker) DeclareQueue(queueName string) {
	_, err := rmq.channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (rmq *rabbitMessageBroker) CloseConnection() {
	rmq.channel.Close()
	rmq.connection.Close()
}

func (rmq *rabbitMessageBroker) getChannel() *amqp.Channel {
	return rmq.channel
}
