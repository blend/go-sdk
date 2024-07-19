/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package testutil_test

import (
	"context"
	"testing"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/testutil"
)

func TestMain(m *testing.M) {
	testutil.MarkUpdateGoldenFlag(testutil.OptUpdateGoldenFlag(testUpdateGoldenFlag))

	testutil.New(
		m,
		testutil.OptLog(logger.All()),
		ensureConnectionOption(),
	).Run()
}

func ensureConnectionOption() testutil.Option {
	return func(s *testutil.Suite) {
		s.Before = append(s.Before, ensureConnection)
	}
}

// ensureConnection makes sure that a valid database connection exists
// before running tests. It uses the helpers **from this package** to
// configure and validation the connection, then closes it since further
// usage will not be needed.
func ensureConnection(ctx context.Context) error {
	c := db.Config{}
	err := testutil.ResolveDBConfig(ctx, &c)
	if err != nil {
		return err
	}

	pool, err := db.New(db.OptConfig(c))
	if err != nil {
		return err
	}

	err = testutil.ValidatePool(ctx, pool, "")
	if err != nil {
		return err
	}

	return pool.Close()
}
