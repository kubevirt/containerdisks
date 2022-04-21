package tests

import (
	"context"
	"time"
)

const (
	maxRetries    = 10
	retryDuration = 10 * time.Second
)

func retryTest(ctx context.Context, testFn func() error) (err error) {
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(retryDuration):
			}
		}

		err = testFn()
		if err == nil {
			return nil
		}
	}

	return err
}
