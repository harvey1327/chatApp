package messagebroker

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publish interface {
	Publish(message publishMessage) error
}

type rabbitPublish struct {
	channel *amqp.Channel
}

func NewRabbitPublish(broker MessageBroker) Publish {
	return &rabbitPublish{
		channel: broker.getChannel(),
	}
}

func (rbtp *rabbitPublish) Publish(message publishMessage) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return rbtp.channel.Publish("", message.queueName, false, false, amqp.Publishing{ContentType: message.contentType, Body: bytes})
}
