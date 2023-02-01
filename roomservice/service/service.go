package service

import (
	"github.com/harvey1327/chatapplib/database"
	"github.com/harvey1327/chatapplib/models/createroom"
	"github.com/harvey1327/chatapplib/proto/generated/roompb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceImpl struct {
	roompb.UnimplementedServiceServer
	eventCol database.EventCommands[createroom.Model]
}

func NewService(eventCol database.EventCommands[createroom.Model]) *ServiceImpl {
	return &ServiceImpl{
		eventCol: eventCol,
	}
}

func (s *ServiceImpl) GetByEventID(request *roompb.GetByEventIDRequest, stream roompb.Service_GetByEventIDServer) error {
	for event := range s.eventCol.ListenByEventID(request.GetEventID()) {
		err := stream.Send(&roompb.EventMessage{
			EventID: event.Data.EventID,
			Status:  roompb.Status(roompb.Status_value[string(event.Data.Status)]),
			Error:   event.Data.Error,
			Time:    timestamppb.New(event.Data.TimeStamp),
		})
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
	return nil
}
