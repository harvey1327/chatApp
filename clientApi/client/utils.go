package client

import (
	"time"

	"github.com/harvey1327/chatapplib/proto/generated/roompb"
	"github.com/harvey1327/chatapplib/proto/generated/userpb"
)

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

func retryNonPendingRoom(attempts int, f func(arguments ...interface{}) (*roompb.EventMessage, error)) (res *roompb.EventMessage, err error) {
	for i := 0; i < attempts; i++ {
		res, err = f()
		if err == nil && res.Status != roompb.Status_PENDING {
			return res, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil, err
}

func retryNonPendingUser(attempts int, f func(arguments ...interface{}) (*userpb.EventMessage, error)) (res *userpb.EventMessage, err error) {
	for i := 0; i < attempts; i++ {
		res, err = f()
		if err == nil && res.Status != userpb.Status_PENDING {
			return res, nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil, err
}
