/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
)

const (
	longQueryTemplate = "SELECT id, pg_sleep(%f) FROM might_sleep WHERE id = 1337;"
	separator         = "=================================================="
)

func createConn(ctx context.Context) (*db.Connection, error) {
	pool, err := db.New(db.OptConfigFromEnv())
	if err != nil {
		return nil, err
	}

	err = pool.Open()
	if err != nil {
		return nil, err
	}

	err = pool.Connection.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	log.Printf("DSN=%q\n", pool.Config.CreateDSN())
	return pool, nil
}

func intentionallyLongQuery(ctx context.Context, pool *db.Connection, cfg *config) error {
	type resultRow struct {
		ID      int    `db:"id"`
		PGSleep string `db:"pg_sleep"`
	}

	s := float64(cfg.PGSleep) / float64(time.Second)

	statement := fmt.Sprintf(longQueryTemplate, s)
	q := pool.QueryContext(ctx, statement)

	r := resultRow{}
	log.Println("Starting query")
	found, err := q.Out(&r)
	if err != nil {
		return err
	}
	if !found {
		return ex.New("`SELECT id, pg_sleep(%f) ...;` query returned no results")
	}

	return nil
}

func main() {
	log.SetFlags(0)
	log.SetOutput(newLogWriter())
	cfg := getConfig()

	// 1. Set the `DB_STATEMENT_TIMEOUT` environment variable.
	log.Println(separator)
	cfg.Print()
	err := cfg.SetEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	deadline := time.Now().Add(cfg.ContextTimeout)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// 2. Parse config / open / ping
	// 3. Make sure `statement_timeout` is in the connection string (it gets printed)
	log.Println(separator)
	pool, err := createConn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanUp(pool)

	// 4. Demonstrate that the observed statement timeout on an open connection is
	//    `StatementTimeout`.
	log.Println(separator)
	timeout, err := ensureStatementTimeout(ctx, pool, cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("statement_timeout=%s\n", timeout)

	// 5. Create a table schema and insert data to seed the database.
	err = seedDatabase(ctx, pool)
	if err != nil {
		log.Fatal(err)
	}

	// 6. Run query that intentionally runs for a long time.
	log.Println(separator)
	err = intentionallyLongQuery(ctx, pool, cfg)
	if err == nil {
		log.Fatal(ex.New("Expected statement contention to occur"))
	}

	// 7. Display the error / errors in as verbose a way as possible.
	log.Println("***")
	err = displayError(err)
	if err != nil {
		log.Fatal(err)
	}
}
