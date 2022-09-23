package messagebroker

import "github.com/google/uuid"

type subscribeMessage[T any] struct {
	queueName    string
	contentType  string
	EventMessage eventMessage[T] `json:"eventMessage"`
}

type publishMessage struct {
	queueName    string
	contentType  string
	EventMessage eventMessage[interface{}] `json:"eventMessage"`
}

type eventMessage[T any] struct {
	EventID string `json:"eventID"`
	Status  status `json:"status"`
	Body    T      `json:"body"`
}

type status string

const (
	PENDING  status = "PENDING"
	COMPLETE status = "COMPLETE"
	FAILED   status = "FAILED"
)

func PublishMessage(body interface{}, queueName string) publishMessage {
	return publishMessage{
		contentType:  "application/json",
		EventMessage: eventMessage[interface{}]{Status: PENDING, Body: body, EventID: uuid.New().String()},
		queueName:    queueName,
	}
}
