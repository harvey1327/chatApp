package messagebroker

import "github.com/google/uuid"

type message[T any] struct {
	queueName    string
	contentType  string
	EventMessage eventMessage[T] `json:"eventMessage"`
}

type eventMessage[T any] struct {
	ID     string `json:"id"`
	Status status `json:"status"`
	Body   T      `json:"body"`
}

type status string

const (
	PENDING  status = "PENDING"
	COMPLETE status = "COMPLETE"
	FAILED   status = "FAILED"
)

func Message[T any](body T, queueName string) message[T] {
	return message[T]{
		contentType:  "application/json",
		EventMessage: eventMessage[T]{Status: PENDING, Body: body, ID: uuid.New().String()},
		queueName:    queueName,
	}
}
