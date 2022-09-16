package messagebroker

import "github.com/google/uuid"

type message struct {
	queueName    string
	contentType  string
	EventMessage eventMessage `json:"eventMessage"`
}

type eventMessage struct {
	ID     string      `json:"id"`
	Status status      `json:"status"`
	Body   interface{} `json:"body"`
}

type status string

const (
	PENDING  status = "PENDING"
	COMPLETE status = "COMPLETE"
	FAILED   status = "FAILED"
)

func Message(body interface{}, queueName string) message {
	return message{
		contentType:  "application/json",
		EventMessage: eventMessage{Status: PENDING, Body: body, ID: uuid.New().String()},
		queueName:    queueName,
	}
}
