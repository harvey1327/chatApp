package service

import (
	"context"

	"github.com/harvey1327/chatapplib/database"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/models/createuser"
	"github.com/harvey1327/chatapplib/proto/generated/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceImpl struct {
	userpb.UnimplementedServiceServer
	commands database.CollectionCommands[messagebroker.EventMessage[createuser.Model]]
}

func NewService(commands database.CollectionCommands[messagebroker.EventMessage[createuser.Model]]) *ServiceImpl {
	return &ServiceImpl{
		commands: commands,
	}
}

func (s *ServiceImpl) GetByEventID(ctx context.Context, request *userpb.GetByEventIDRequest) (*userpb.EventMessage, error) {
	eventID := request.GetEventID()
	res, err := s.commands.FindSingleByQuery(database.Query("eventID", eventID))
	if err != nil {
		if err == database.EMPTY {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &userpb.EventMessage{
		EventID: res.Data.EventID,
		Status:  userpb.Status(userpb.Status_value[string(res.Data.Status)]),
		Body:    &userpb.Model{DisplayName: res.Data.Body.DisplayName},
		Error:   res.Data.Error,
		Time:    timestamppb.New(res.Data.TimeStamp),
	}, nil
}
