package service

import (
	"context"

	"github.com/harvey1327/chatapplib/database"
	"github.com/harvey1327/chatapplib/messagebroker"
	"github.com/harvey1327/chatapplib/models/createroom"
	"github.com/harvey1327/chatapplib/proto/generated/roompb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceImpl struct {
	roompb.UnimplementedServiceServer
	commands database.CollectionCommands[messagebroker.EventMessage[createroom.Model]]
}

func NewService(commands database.CollectionCommands[messagebroker.EventMessage[createroom.Model]]) *ServiceImpl {
	return &ServiceImpl{
		commands: commands,
	}
}

func (s *ServiceImpl) GetByEventID(ctx context.Context, request *roompb.GetByEventIDRequest) (*roompb.EventMessage, error) {
	eventID := request.GetEventID()
	res, err := s.commands.FindSingleByQuery(database.Query("eventID", eventID))
	if err != nil {
		if err == database.EMPTY {
			return nil, status.Error(codes.NotFound, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &roompb.EventMessage{
		EventID: res.Data.EventID,
		Status:  roompb.Status(roompb.Status_value[string(res.Data.Status)]),
		Body:    &roompb.Model{DisplayName: res.Data.Body.DisplayName},
		Error:   res.Data.Error,
		Time:    timestamppb.New(res.Data.TimeStamp),
	}, nil
}
