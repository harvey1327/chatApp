package messagebroker

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageBrokerCommands[T any] interface {
	Publish(message message[T]) error
	Subscribe(queueName string) <-chan eventMessage[T]
}

type rabbitMessageBrokerCommands[T any] struct {
	channel *amqp.Channel
}

func NewRabbitMQCommands[T any](broker MessageBroker) MessageBrokerCommands[T] {
	return &rabbitMessageBrokerCommands[T]{
		channel: broker.getChannel(),
	}
}

func (rmq *rabbitMessageBrokerCommands[T]) Publish(message message[T]) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return rmq.channel.Publish("", message.queueName, false, false, amqp.Publishing{ContentType: message.contentType, Body: bytes})
}

func (rmq *rabbitMessageBrokerCommands[T]) Subscribe(queueName string) <-chan eventMessage[T] {
	results := make(chan eventMessage[T])
	msgs, err := rmq.channel.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			received, ok := <-msgs
			if !ok {
				break
			}
			event := message[T]{}
			err := json.Unmarshal(received.Body, &event)
			if err != nil {
				log.Fatal(err)
			}
			results <- event.EventMessage
		}
		close(results)
	}()
	return results
}
