package main

import (
	"fmt"
	"log"
	"net"

	"github.com/harvey1327/chatapp/libdatabase"
	"github.com/harvey1327/chatapp/libmessagebroker"
	"github.com/harvey1327/chatapp/libmessagebroker/events/createuser"
	"github.com/harvey1327/chatapp/libproto/generated/userpb"
	"github.com/harvey1327/chatapp/userservice/interceptor"
	"github.com/harvey1327/chatapp/userservice/service"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50052))
	if err != nil {
		log.Fatal("failed to listen: 50052")
	} else {
		log.Printf("listening on %s", lis.Addr().String())
	}

	db := libdatabase.NewDB(libdatabase.USER)
	defer db.Close()
	commands := libdatabase.NewCollection[libmessagebroker.EventMessage[createuser.Model]](db, createuser.QUEUE_NAME)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryLoggerInterceptor()))

	userpb.RegisterServiceServer(grpcServer, service.NewService(commands))

	if error := grpcServer.Serve(lis); error != nil {
		log.Fatalf("failed to serve: %s", error)
	}
}
