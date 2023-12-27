/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package retry

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

func Test_Retry(t *testing.T) {
	its := assert.New(t)

	passedArgs := make(chan interface{}, 5)
	results := make(chan interface{}, 5)

	var internalAttempt int
	action := ActionerFunc(func(ctx context.Context, args interface{}) (interface{}, error) {
		defer func() { internalAttempt++ }()
		passedArgs <- args
		if internalAttempt < 4 {
			err := fmt.Errorf("attempt %d", internalAttempt)
			results <- err
			return nil, err
		}
		result := "OK!"
		results <- result
		return result, nil
	})

	result, err := New(OptConstantDelay(time.Millisecond)).Intercept(action).Action(its.Background(), "args")
	its.Nil(err)
	its.Equal("OK!", result)
	its.Len(passedArgs, 5)
	its.Equal("args", <-passedArgs)
	its.Equal("args", <-passedArgs)
	its.Equal("args", <-passedArgs)
	its.Equal("args", <-passedArgs)
	its.Equal("args", <-passedArgs)
	its.Len(results, 5)
	its.Equal(fmt.Errorf("attempt 0"), <-results)
	its.Equal(fmt.Errorf("attempt 1"), <-results)
	its.Equal(fmt.Errorf("attempt 2"), <-results)
	its.Equal(fmt.Errorf("attempt 3"), <-results)
	its.Equal("OK!", <-results)
}

func Test_Retry_ShouldRetryProvider(t *testing.T) {
	its := assert.New(t)

	passedArgs := make(chan interface{}, 5)
	results := make(chan interface{}, 5)

	var internalAttempt int
	action := ActionerFunc(func(ctx context.Context, args interface{}) (interface{}, error) {
		defer func() { internalAttempt++ }()
		passedArgs <- args

		if internalAttempt < 4 {
			err := fmt.Errorf("attempt %d", internalAttempt)
			results <- err
			return nil, err
		}
		result := "OK!"
		results <- result
		return result, nil
	})

	result, err := New(
		OptConstantDelay(time.Millisecond),
		OptShouldRetryProvider(func(err error) bool {
			return err.Error() != "attempt 3"
		}),
	).Intercept(action).Action(its.Background(), "args")
	its.NotNil(err)
	its.Empty(result)
	its.Len(passedArgs, 4)
	its.Equal("args", <-passedArgs)
	its.Equal("args", <-passedArgs)
	its.Equal("args", <-passedArgs)
	its.Equal("args", <-passedArgs)
	its.Len(results, 4)
	its.Equal(fmt.Errorf("attempt 0"), <-results)
	its.Equal(fmt.Errorf("attempt 1"), <-results)
	its.Equal(fmt.Errorf("attempt 2"), <-results)
	its.Equal(fmt.Errorf("attempt 3"), <-results)
}
