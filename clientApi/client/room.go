package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/harvey1327/chatapplib/proto/generated/roompb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RoomClient interface {
	GetByEventID(eventID string) (<-chan *roompb.EventMessage, error)
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

func (c *roomClient) GetByEventID(eventID string) (<-chan *roompb.EventMessage, error) {
	stream, err := c.client.GetByEventID(context.Background(), &roompb.GetByEventIDRequest{EventID: eventID})
	if err != nil {
		return nil, err
	}
	res := make(chan *roompb.EventMessage)

	go func() {
		for {
			event, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			res <- event
		}
		close(res)
	}()

	return res, nil
}

func (c *roomClient) Close() error {
	return c.conn.Close()
}
