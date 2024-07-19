/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package breaker

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
)

func TestNew(t *testing.T) {
	its := assert.New(t)

	b := New()
	its.Equal(DefaultHalfOpenMaxActions, b.HalfOpenMaxActions)
	its.Equal(DefaultOpenExpiryInterval, b.OpenExpiryInterval)
	its.Equal(DefaultClosedExpiryInterval, b.ClosedExpiryInterval)
}

func TestNewOptions(t *testing.T) {
	its := assert.New(t)

	b := New(
		OptHalfOpenMaxActions(5),
		OptOpenExpiryInterval(10*time.Second),
		OptClosedExpiryInterval(20*time.Second),
	)
	its.Equal(5, b.HalfOpenMaxActions)
	its.Equal(10*time.Second, b.OpenExpiryInterval)
	its.Equal(20*time.Second, b.ClosedExpiryInterval)
}

func createTestBreaker() *Breaker {
	return New(OptClosedExpiryInterval(0))
}

func succeed(b *Breaker) error {
	_, err := b.Intercept(ActionerFunc(func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, nil
	})).Action(context.Background(), nil)
	return err
}

func fail(b *Breaker) error {
	msg := "fail"
	_, err := b.Intercept(ActionerFunc(func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, fmt.Errorf(msg)
	})).Action(context.Background(), nil)
	if err != nil && err.Error() == msg {
		return nil
	}
	return err
}

func pseudoSleep(b *Breaker, period time.Duration) {
	if !b.stateExpiresAt.IsZero() {
		b.stateExpiresAt = b.stateExpiresAt.Add(-period)
	}
}

func Test_Breaker(t *testing.T) {
	its := assert.New(t)
	ctx := context.Background()

	b := createTestBreaker()

	for i := 0; i < 5; i++ {
		its.Nil(fail(b))
	}
	its.Equal(StateClosed, b.EvaluateState(ctx))
	its.Equal(Counts{5, 0, 5, 0, 5}, b.Counts)

	its.Nil(succeed(b))
	its.Equal(StateClosed, b.EvaluateState(ctx))
	its.Equal(Counts{6, 1, 5, 1, 0}, b.Counts)

	its.Nil(fail(b))
	its.Equal(StateClosed, b.EvaluateState(ctx))
	its.Equal(Counts{7, 1, 6, 0, 1}, b.Counts)

	// StateClosed to StateOpen
	for i := 0; i < 5; i++ {
		its.Nil(fail(b)) // 5 more consecutive failures
	}
	its.Equal(StateOpen, b.EvaluateState(ctx))
	its.Equal(Counts{0, 0, 0, 0, 0}, b.Counts)
	its.False(b.stateExpiresAt.IsZero())

	err := succeed(b)
	its.True(ErrIsOpen(err)) // this shouldn't have called the action, should yield closed
	err = fail(b)
	its.True(ErrIsOpen(err)) // this shouldn't have called the action either
	its.Equal(Counts{0, 0, 0, 0, 0}, b.Counts)

	pseudoSleep(b, 59*time.Second) // push forward time by 59s
	its.Equal(StateOpen, b.EvaluateState(ctx))

	// StateOpen to StateHalfOpen
	pseudoSleep(b, 2*time.Second) // over Timeout
	its.Equal(StateHalfOpen, b.EvaluateState(ctx))
	its.True(b.stateExpiresAt.IsZero())

	// StateHalfOpen to StateOpen
	// there are like, (3) calls queued here
	its.Nil(fail(b))
	its.Equal(StateOpen, b.EvaluateState(ctx))
	its.Equal(Counts{0, 0, 0, 0, 0}, b.Counts)
	its.False(b.stateExpiresAt.IsZero())

	// StateOpen to StateHalfOpen
	pseudoSleep(b, time.Duration(60)*time.Second)
	its.Equal(StateHalfOpen, b.EvaluateState(ctx))
	its.True(b.stateExpiresAt.IsZero())

	// StateHalfOpen to StateClosed
	its.Nil(succeed(b))
	its.Equal(StateClosed, b.EvaluateState(ctx))
	its.Equal(Counts{0, 0, 0, 0, 0}, b.Counts)
	its.True(b.stateExpiresAt.IsZero())
}

func Test_Breaker_ErrStateOpen(t *testing.T) {
	its := assert.New(t)

	var didCall bool
	b := New()
	b.state = StateOpen
	b.stateExpiresAt = time.Now().Add(time.Hour)
	_, err := b.Intercept(ActionerFunc(func(_ context.Context, _ interface{}) (interface{}, error) {
		didCall = true
		return nil, nil
	})).Action(context.Background(), nil)
	its.True(ex.Is(err, ErrOpenState), fmt.Sprintf("%v", err))
	its.False(didCall)
}

func Test_Breaker_ErrTooManyRequests(t *testing.T) {
	its := assert.New(t)

	var didCall bool
	b := New()

	b.state = StateHalfOpen
	b.Counts.Requests = 10
	b.HalfOpenMaxActions = 5

	_, err := b.Intercept(ActionerFunc(func(_ context.Context, _ interface{}) (interface{}, error) {
		didCall = true
		return nil, nil
	})).Action(context.Background(), nil)
	its.True(ex.Is(err, ErrTooManyRequests))
	its.False(didCall)
}

func Test_Breaker_callsOnOpenHandler(t *testing.T) {
	its := assert.New(t)

	var didCall, didCallOpen bool
	b := New(
		OptOpenAction(ActionerFunc(func(_ context.Context, _ interface{}) (interface{}, error) {
			didCallOpen = true
			return "on open", nil
		})),
	)

	b.state = StateOpen
	b.stateExpiresAt = time.Now().Add(time.Hour)

	res, err := b.Intercept(ActionerFunc(func(_ context.Context, _ interface{}) (interface{}, error) {
		didCall = true
		return nil, nil
	})).Action(context.Background(), nil)
	its.Nil(err)
	its.False(didCall)
	its.True(didCallOpen)
	its.Equal("on open", res)
}
