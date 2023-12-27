/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package status

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
)

func Test_Freeform_CheckStatuses(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	service0 := uuid.V4().String()
	service1 := uuid.V4().String()

	var calledService0, calledService1 bool
	ff := NewFreeform()
	ff.ServiceChecks[service0] = CheckerFunc(func(_ context.Context) error {
		calledService0 = true
		return nil
	})
	ff.ServiceChecks[service1] = CheckerFunc(func(_ context.Context) error {
		calledService1 = true
		return fmt.Errorf("oh noes")
	})

	statuses, err := ff.CheckStatuses(context.Background())
	its.Nil(err)
	its.Len(statuses, 2)
	its.True(calledService0)
	its.True(calledService1)
	its.True(statuses[service0])
	its.False(statuses[service1])
}

func Test_Freeform_CheckStatuses_servicesToCheck(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	service0 := uuid.V4().String()
	service1 := uuid.V4().String()
	service2 := uuid.V4().String()

	var calledService0, calledService1, calledService2 bool
	ff := NewFreeform()
	ff.ServiceChecks[service0] = CheckerFunc(func(_ context.Context) error {
		calledService0 = true
		return nil
	})
	ff.ServiceChecks[service1] = CheckerFunc(func(_ context.Context) error {
		calledService1 = true
		return fmt.Errorf("oh noes")
	})
	ff.ServiceChecks[service2] = CheckerFunc(func(_ context.Context) error {
		calledService2 = true
		return fmt.Errorf("oh noes")
	})

	statuses, err := ff.CheckStatuses(context.Background(), service0, service2)
	its.Nil(err)
	its.Len(statuses, 2)
	its.True(calledService0)
	its.False(calledService1)
	its.True(calledService2)
	its.True(statuses[service0])
	its.False(statuses[service2])
}

func Test_Freeform_CheckStatuses_missing(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	service0 := uuid.V4().String()
	service1 := uuid.V4().String()
	service2 := uuid.V4().String()

	var calledService0, calledService1 bool
	ff := NewFreeform()
	ff.ServiceChecks[service0] = CheckerFunc(func(_ context.Context) error {
		calledService0 = true
		return nil
	})
	ff.ServiceChecks[service1] = CheckerFunc(func(_ context.Context) error {
		calledService1 = true
		return fmt.Errorf("oh noes")
	})

	statuses, err := ff.CheckStatuses(context.Background(), service1, service2)
	its.NotNil(err)
	its.Empty(statuses)

	its.False(calledService0)
	its.False(calledService1)
}

func Test_Freeform_getCheckStatus(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	service0 := uuid.V4().String()

	var calledService0 bool

	ff := NewFreeform()
	res := ff.getCheckStatus(context.Background(), service0, CheckerFunc(func(_ context.Context) error {
		calledService0 = true
		return nil
	}))
	its.True(calledService0)
	its.True(res.Ok)
	its.Equal(service0, res.ServiceName)
}

func Test_Freeform_getCheckStatus_logsError(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	logBuffer := new(bytes.Buffer)
	log := logger.Memory(logBuffer,
		logger.OptText(
			logger.OptTextHideTimestamp(),
			logger.OptTextNoColor(),
		),
	)
	defer log.Close()

	service0 := uuid.V4().String()
	var calledService0 bool
	ff := NewFreeform(
		OptFreeformLog(log),
	)
	res := ff.getCheckStatus(context.Background(), service0, CheckerFunc(func(_ context.Context) error {
		calledService0 = true
		return fmt.Errorf("this is just a test")
	}))
	its.True(calledService0)
	its.False(res.Ok)
	its.Equal(service0, res.ServiceName)

	its.Equal("[error] this is just a test\n", logBuffer.String())
}

func Test_Freeform_getCheckStatus_timeout(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var calledService bool
	serviceID := uuid.V4().String()

	called := make(chan struct{})
	returned := make(chan struct{})

	ff := NewFreeform(
		OptFreeformTimeout(time.Millisecond),
	)

	var res freeformCheckResult
	go func() {
		defer close(returned)
		res = ff.getCheckStatus(context.Background(), serviceID, CheckerFunc(func(ctx context.Context) error {
			calledService = true
			close(called)
			<-ctx.Done()
			return context.Canceled
		}))
	}()
	<-called
	<-returned
	its.True(calledService)
	its.False(res.Ok)
	its.Equal(serviceID, res.ServiceName)
	its.Equal(context.Canceled, res.Err)
}

func Test_Freeform_getCheckStatus_panic(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var calledService bool
	serviceID := uuid.V4().String()

	ff := NewFreeform(
		OptFreeformTimeout(time.Millisecond),
	)

	res := ff.getCheckStatus(context.Background(), serviceID, CheckerFunc(func(ctx context.Context) error {
		calledService = true
		panic("just a test panic!")
	}))
	its.True(calledService)
	its.False(res.Ok)
	its.Equal(serviceID, res.ServiceName)
	its.NotNil(res.Err)
}

func Test_Freeform_timeoutOrDefault(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	its.Equal(DefaultFreeformTimeout, NewFreeform().timeoutOrDefault())
	its.Equal(DefaultFreeformTimeout<<1, NewFreeform(OptFreeformTimeout(DefaultFreeformTimeout<<1)).timeoutOrDefault())
}

func Test_Freeform_serviceChecksOrDefault(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	set0 := []string{"alpha", "bravo", "charlie"}
	set1 := []string{"alpha", "charlie"}
	set2 := []string{"alpha", uuid.V4().String()}

	ff := NewFreeform()
	for _, key := range set0 {
		ff.ServiceChecks[key] = CheckerFunc(func(_ context.Context) error { return nil })
	}

	serviceChecks, err := ff.serviceChecksOrDefault()
	its.Nil(err)
	its.Len(serviceChecks, 3)
	its.NotNil(serviceChecks["alpha"])
	its.NotNil(serviceChecks["bravo"])
	its.NotNil(serviceChecks["charlie"])

	serviceChecks, err = ff.serviceChecksOrDefault(set1...)
	its.Nil(err)
	its.Len(serviceChecks, 2)
	its.NotNil(serviceChecks["alpha"])
	its.Nil(serviceChecks["bravo"])
	its.NotNil(serviceChecks["charlie"])

	serviceChecks, err = ff.serviceChecksOrDefault(set2...)
	its.NotNil(err)
	its.Empty(serviceChecks)
}
