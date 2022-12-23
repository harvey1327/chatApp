package client

import "time"

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
