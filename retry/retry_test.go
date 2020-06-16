package retry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func TestRetry(t *testing.T) {
	assert := assert.New(t)

	results := make(chan interface{}, 5)

	var internalAttempt int
	action := Action(func(ctx context.Context) (interface{}, error) {
		defer func() { internalAttempt++ }()

		if internalAttempt < 4 {
			err := fmt.Errorf("attempt %d", internalAttempt)
			results <- err
			return nil, err
		}
		result := "OK!"
		results <- result
		return result, nil
	})

	result, err := Retry(context.Background(), action, OptConstantDelay(time.Millisecond))
	assert.Nil(err)
	assert.Equal("OK!", result)
	assert.Len(results, 5)
	assert.Equal(fmt.Errorf("attempt 0"), <-results)
	assert.Equal(fmt.Errorf("attempt 1"), <-results)
	assert.Equal(fmt.Errorf("attempt 2"), <-results)
	assert.Equal(fmt.Errorf("attempt 3"), <-results)
	assert.Equal("OK!", <-results)
}

func TestRetry_ShouldRetryProvider(t *testing.T) {
	assert := assert.New(t)

	results := make(chan interface{}, 5)

	var internalAttempt int
	action := Action(func(ctx context.Context) (interface{}, error) {
		defer func() { internalAttempt++ }()

		if internalAttempt < 4 {
			err := fmt.Errorf("attempt %d", internalAttempt)
			results <- err
			return nil, err
		}
		result := "OK!"
		results <- result
		return result, nil
	})

	result, err := Retry(
		context.Background(),
		action,
		OptConstantDelay(time.Millisecond),
		OptShouldRetryProvider(func(err error) bool {
			return err.Error() != "attempt 3"
		}),
	)
	assert.NotNil(err)
	assert.Empty(result)
	assert.Len(results, 4)
	assert.Equal(fmt.Errorf("attempt 0"), <-results)
	assert.Equal(fmt.Errorf("attempt 1"), <-results)
	assert.Equal(fmt.Errorf("attempt 2"), <-results)
	assert.Equal(fmt.Errorf("attempt 3"), <-results)
}
