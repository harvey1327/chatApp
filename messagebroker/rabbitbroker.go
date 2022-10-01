package messagebroker

import (
	"log"

	"github.com/chatapp/messagebroker/events/createuser"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBroker interface {
	CloseConnection()
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

	broker := &rabbitMessageBroker{
		connection: connection,
		channel:    channel,
	}

	broker.declareQueue(createuser.QUEUE_NAME)

	return broker
}

func (rmq *rabbitMessageBroker) declareQueue(queueName string) {
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
