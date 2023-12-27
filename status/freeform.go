/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package status

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/web"
)

// NewFreeform returns a new freeform check aggregator.
func NewFreeform(opts ...FreeformOption) *Freeform {
	ff := &Freeform{
		ServiceChecks: make(map[string]Checker),
	}
	for _, opt := range opts {
		opt(ff)
	}
	return ff
}

// FreeformOption mutates a freeform check.
type FreeformOption func(*Freeform)

// OptFreeformTimeout sets the timeout.
func OptFreeformTimeout(d time.Duration) FreeformOption {
	return func(ff *Freeform) {
		ff.Timeout = d
	}
}

// OptFreeformLog sets the logger.
func OptFreeformLog(log logger.Log) FreeformOption {
	return func(ff *Freeform) {
		ff.Log = log
	}
}

// Freeform holds sla check actions.
type Freeform struct {
	// Timeout serves as an overall timeout, but is
	// enforced per action.
	Timeout	time.Duration
	// ServiceChecks are the individual checks
	// we'll try as part of the status.
	ServiceChecks	map[string]Checker
	// Log is a reference to a logger for
	// situations where there are errors.
	Log	logger.Log
}

// Endpoint implements the handler for a given list of services.
func (f Freeform) Endpoint(servicesToCheck ...string) web.Action {
	return func(r *web.Ctx) web.Result {
		results, err := f.CheckStatuses(r.Context(), servicesToCheck...)
		if err != nil {
			return web.JSON.InternalError(err)
		}

		statusCode := http.StatusOK
		for _, status := range results {
			if !status {
				statusCode = http.StatusServiceUnavailable
				break
			}
		}
		return web.JSON.Status(statusCode, results)
	}
}

// CheckStatuses runs the check statuses for a given list of service names.
func (f Freeform) CheckStatuses(ctx context.Context, servicesToCheck ...string) (FreeformResult, error) {
	servicesOrDefault, err := f.serviceChecksOrDefault(servicesToCheck...)
	if err != nil {
		return nil, err
	}

	results := make(chan freeformCheckResult, len(servicesOrDefault))
	wg := sync.WaitGroup{}
	wg.Add(len(servicesOrDefault))
	for serviceName, check := range servicesOrDefault {
		go func(ictx context.Context, sn string, c Checker) {
			defer wg.Done()
			results <- f.getCheckStatus(ictx, sn, c)
		}(ctx, serviceName, check)
	}
	wg.Wait()

	// handle results
	output := make(FreeformResult)
	resultCount := len(results)
	for x := 0; x < resultCount; x++ {
		res := <-results
		output[res.ServiceName] = res.Ok
	}
	return output, nil
}

// getCheckStatus runs a check for a given serviceName.
func (f Freeform) getCheckStatus(ctx context.Context, serviceName string, checkAction Checker) (res freeformCheckResult) {
	res.ServiceName = serviceName
	defer func() {
		if r := recover(); r != nil {
			res.Err = ex.Append(res.Err, ex.New(r))
		}
		if res.Err != nil {
			logger.MaybeErrorContext(ctx, f.Log, res.Err)
		} else {
			res.Ok = true
		}
		return
	}()
	timeoutCtx, cancel := context.WithTimeout(ctx, f.timeoutOrDefault())
	defer cancel()
	res.Err = checkAction.Check(timeoutCtx)
	return
}

//
// Private / Internal
//

// timeoutOrDefault returns the timeout or a default.
func (f Freeform) timeoutOrDefault() time.Duration {
	if f.Timeout > 0 {
		return f.Timeout
	}
	return DefaultFreeformTimeout
}

// serviceChecksOrDefault returns either the full ServiceCheck's list
// or a subset of those checks. If any element of the subset is not found,
// an error will be returned.
func (f Freeform) serviceChecksOrDefault(servicesToCheck ...string) (map[string]Checker, error) {
	if len(servicesToCheck) == 0 {
		return f.ServiceChecks, nil
	}

	servicesChecks := make(map[string]Checker)
	for _, serviceName := range servicesToCheck {
		check, ok := f.ServiceChecks[serviceName]
		if !ok {
			return nil, ex.New(ErrServiceCheckNotDefined, ex.OptMessagef("service: %s", serviceName))
		}
		servicesChecks[serviceName] = check
	}
	return servicesChecks, nil
}
