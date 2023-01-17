package service

import (
	"context"

	"github.com/harvey1327/chatapplib/database"
	"github.com/harvey1327/chatapplib/models/createroom"
	"github.com/harvey1327/chatapplib/models/message"
	"github.com/harvey1327/chatapplib/proto/generated/roompb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceImpl struct {
	roompb.UnimplementedServiceServer
	eventCol database.CollectionCommands[message.EventMessage[createroom.Model]]
}

func NewService(eventCol database.CollectionCommands[message.EventMessage[createroom.Model]]) *ServiceImpl {
	return &ServiceImpl{
		eventCol: eventCol,
	}
}

func (s *ServiceImpl) GetByEventID(ctx context.Context, request *roompb.GetByEventIDRequest) (*roompb.EventMessage, error) {
	eventID := request.GetEventID()
	res, err := s.eventCol.FindSingleByQuery(database.Query("eventID", eventID))
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
		Error:   res.Data.Error,
		Time:    timestamppb.New(res.Data.TimeStamp),
	}, nil
}
