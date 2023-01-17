package client

import (
	"context"
	"fmt"
	"log"

	"github.com/harvey1327/chatapplib/proto/generated/roompb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RoomClient interface {
	GetByEventID(eventID string) (*roompb.EventMessage, error)
	Close() error
}

type roomClient struct {
	client roompb.ServiceClient
	conn   *grpc.ClientConn
}

func NewRoomClient(host string, port int) RoomClient {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), options...)
	if err != nil {
		log.Fatal(err)
	}
	return &roomClient{
		client: roompb.NewServiceClient(conn),
		conn:   conn,
	}
}

func (c *roomClient) GetByEventID(eventID string) (*roompb.EventMessage, error) {
	return retryNonPendingRoom(5, func(arguments ...interface{}) (*roompb.EventMessage, error) {
		return c.client.GetByEventID(context.TODO(), &roompb.GetByEventIDRequest{EventID: eventID})
	})
}

func (c *roomClient) Close() error {
	return c.conn.Close()
}
