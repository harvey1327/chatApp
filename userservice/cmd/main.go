package main

import (
	"fmt"
	"log"
	"net"

	"github.com/harvey1327/chatapp/userservice/interceptor"
	"github.com/harvey1327/chatapp/userservice/service"
	"github.com/harvey1327/chatapplib/database"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/messagebroker/events/createuser"
	"github.com/harvey1327/chatapplib/proto/generated/userpb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50052))
	if err != nil {
		log.Fatal("failed to listen: 50052")
	} else {
		log.Printf("listening on %s", lis.Addr().String())
	}

	db := database.NewDB(database.USER)
	defer db.Close()
	commands := database.NewCollection[messagebroker.EventMessage[createuser.Model]](db, createuser.QUEUE_NAME)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryLoggerInterceptor()))

	userpb.RegisterServiceServer(grpcServer, service.NewService(commands))

	if error := grpcServer.Serve(lis); error != nil {
		log.Fatalf("failed to serve: %s", error)
	}
}
