package service

import (
	"github.com/harvey1327/chatapplib/database"
	"github.com/harvey1327/chatapplib/models/createuser"
	"github.com/harvey1327/chatapplib/proto/generated/userpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ServiceImpl struct {
	userpb.UnimplementedServiceServer
	eventCol database.EventCommands[createuser.Model]
}

func NewService(eventCol database.EventCommands[createuser.Model]) *ServiceImpl {
	return &ServiceImpl{
		eventCol: eventCol,
	}
}

func (s *ServiceImpl) GetByEventID(request *userpb.GetByEventIDRequest, stream userpb.Service_GetByEventIDServer) error {
	for event := range s.eventCol.ListenByEventID(request.GetEventID()) {
		err := stream.Send(&userpb.EventMessage{
			EventID: event.Data.EventID,
			Status:  userpb.Status(userpb.Status_value[string(event.Data.Status)]),
			Error:   event.Data.Error,
			Time:    timestamppb.New(event.Data.TimeStamp),
		})
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
	return nil
}
