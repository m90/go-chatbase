package chatbase

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestResultWithContext(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		result, err := resultWithContext(ctx, func() (interface{}, error) {
			time.Sleep(time.Second)
			return 12345, nil
		})
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		if result != 12345 {
			t.Errorf("Unexpected result %v", result)
		}
	})
	t.Run("error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		result, err := resultWithContext(ctx, func() (interface{}, error) {
			time.Sleep(time.Second)
			return nil, errors.New("broke")
		})
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if result != nil {
			t.Errorf("Unexpected result %v", result)
		}
	})
	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		result, err := resultWithContext(ctx, func() (interface{}, error) {
			time.Sleep(time.Minute)
			return 12345, nil
		})
		if err == nil {
			t.Error("Expected error, got nil")
		}
		if result != nil {
			t.Errorf("Unexpected result %v", result)
		}
	})
}

func TestWithContext(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := withContext(ctx, func() error {
			time.Sleep(time.Second)
			return nil
		})
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
	})
	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := withContext(ctx, func() error {
			time.Sleep(time.Minute)
			return nil
		})
		if err == nil {
			t.Error("Expected error, got nil")
		}
	})
}
