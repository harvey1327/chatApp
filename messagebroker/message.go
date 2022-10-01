package messagebroker

import (
	"time"

	"github.com/google/uuid"
)

type subscribeMessage[T any] struct {
	queueName    string
	contentType  string
	EventMessage EventMessage[T] `json:"eventMessage"`
}

type publishMessage struct {
	queueName    string
	contentType  string
	EventMessage EventMessage[interface{}] `json:"eventMessage"`
}

type EventMessage[T any] struct {
	EventID   string    `json:"eventID" bson:"eventID"`
	Status    status    `json:"status" bson:"status"`
	Body      T         `json:"body" bson:",inline"`
	Error     string    `json:"error,omitempty" bson:"error"`
	TimeStamp time.Time `json:"time" bson:"time"`
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
		EventMessage: EventMessage[interface{}]{Status: PENDING, Body: body, EventID: uuid.New().String(), TimeStamp: time.Now().UTC()},
		queueName:    queueName,
	}
}

func (event *EventMessage[T]) Complete() {
	event.Status = COMPLETE
}

func (event *EventMessage[T]) Failed(reason string) {
	event.Status = FAILED
	event.Error = reason
}
