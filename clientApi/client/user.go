package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/harvey1327/chatapplib/proto/generated/userpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient interface {
	GetByEventID(eventID string) (*userpb.EventMessage, error)
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

func (c *userClient) GetByEventID(eventID string) (*userpb.EventMessage, error) {
	return retry(5, func(arguments ...interface{}) (*userpb.EventMessage, error) {
		return c.client.GetByEventID(context.TODO(), &userpb.GetByEventIDRequest{EventID: eventID})
	})
}

func (c *userClient) Close() error {
	return c.conn.Close()
}

func retry[T any](attempts int, f func(arguments ...interface{}) (T, error)) (res T, err error) {
	for i := 0; i < attempts; i++ {
		res, err = f()
		if err == nil {
			return res, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return *new(T), err
}
