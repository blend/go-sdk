/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package main

import (
	"context"
	"time"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
)

func verifyLockTimeout(ctx context.Context, pool *db.Connection) (time.Duration, error) {
	type lockTimeoutRow struct {
		LockTimeout string `db:"lock_timeout"`
	}

	q := pool.QueryContext(ctx, "SHOW lock_timeout;")
	r := lockTimeoutRow{}
	found, err := q.Out(&r)
	if err != nil {
		return 0, err
	}
	if !found {
		return 0, ex.New("`SHOW lock_timeout;` query returned no results")
	}

	d, err := time.ParseDuration(r.LockTimeout)
	if err != nil {
		return 0, err
	}

	return d, nil
}

func ensureLockTimeout(ctx context.Context, pool *db.Connection, cfg *config) (time.Duration, error) {
	timeout, err := verifyLockTimeout(ctx, pool)
	if err != nil {
		return 0, err
	}

	if timeout != cfg.LockTimeout {
		err = ex.New(
			"Expected the default lock timeout to be set",
			ex.OptMessagef("Timeout: %s, Expected: %s", timeout, cfg.LockTimeout),
		)
		return 0, err
	}

	return timeout, nil
}
