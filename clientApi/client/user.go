package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/harvey1327/chatapplib/proto/generated/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient interface {
	GetByEventID(eventID string) (<-chan *userpb.EventMessage, error)
	Close() error
}

type userClient struct {
	client userpb.ServiceClient
	conn   *grpc.ClientConn
}

func NewUserClient(host string, port int) UserClient {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), options...)
	if err != nil {
		log.Fatal(err)
	}
	return &userClient{
		client: userpb.NewServiceClient(conn),
		conn:   conn,
	}
}

func (c *userClient) GetByEventID(eventID string) (<-chan *userpb.EventMessage, error) {
	stream, err := c.client.GetByEventID(context.Background(), &userpb.GetByEventIDRequest{EventID: eventID})
	if err != nil {
		return nil, err
	}
	res := make(chan *userpb.EventMessage)

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

func (c *userClient) Close() error {
	return c.conn.Close()
}
