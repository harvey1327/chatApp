package main

import (
	"fmt"
	"log"
	"net"

	"github.com/harvey1327/chatapp/roomservice/config"
	"github.com/harvey1327/chatapp/roomservice/interceptor"
	"github.com/harvey1327/chatapp/roomservice/service"
	"github.com/harvey1327/chatapplib/database"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/messagebroker/events/createroom"
	"github.com/harvey1327/chatapplib/proto/generated/roompb"
	"google.golang.org/grpc"
)

func main() {
	conf := config.Load()
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", conf.HOST, conf.PORT))
	if err != nil {
		log.Fatalf("failed to listen: %s:%d", conf.HOST, conf.PORT)
	} else {
		log.Printf("listening on %s", lis.Addr().String())
	}

	db := database.NewDB(database.USER, database.DBConfig(conf.DB_HOST, conf.DB_PORT, conf.DB_USERNAME, conf.DB_PASSWORD))
	defer db.Close()
	commands := database.NewCollection[messagebroker.EventMessage[createroom.Model]](db, createroom.QUEUE_NAME)

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.UnaryLoggerInterceptor()))

	roompb.RegisterServiceServer(grpcServer, service.NewService(commands))

	if error := grpcServer.Serve(lis); error != nil {
		log.Fatalf("failed to serve: %s", error)
	}
}
