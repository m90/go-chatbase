package chatbase

import (
	"context"
	"errors"
)

var errBadData = errors.New("unexpected return value from submit call")

func resultWithContext(ctx context.Context, thunk func() (interface{}, error)) (interface{}, error) {
	out := make(chan interface{})
	errors := make(chan error)
	go func() {
		data, err := thunk()
		if err != nil {
			errors <- err
		}
		out <- data
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errors:
		return <-out, err
	case res := <-out:
		return res, nil
	}
}

func withContext(ctx context.Context, thunk func() error) error {
	errors := make(chan error)
	go func() {
		errors <- thunk()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errors:
		return err
	}
}
