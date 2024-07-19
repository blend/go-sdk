/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/ex"
)

const (
	updateRows = "UPDATE might_deadlock SET counter = counter + 1 WHERE key = $1;"
	separator  = "=================================================="
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

// contendReads introduces two reads (in a transaction) with a sleep in
// between.
// H/T to https://www.citusdata.com/blog/2018/02/22/seven-tips-for-dealing-with-postgres-locks/
// for the idea on how to "easily" introduce a deadlock.
func contendReads(ctx context.Context, wg *sync.WaitGroup, tx *sql.Tx, key1, key2 string, cfg *config) error {
	defer wg.Done()

	_, err := tx.ExecContext(ctx, updateRows, key1)
	if err != nil {
		return err
	}

	time.Sleep(cfg.TxSleep)
	_, err = tx.ExecContext(ctx, updateRows, key2)
	if err == context.DeadlineExceeded {
		return nest(err, ex.New("Context cancel in between queries"))
	}
	return err
}

func intentionalContention(ctx context.Context, pool *db.Connection, cfg *config) (err error) {
	var tx1, tx2 *sql.Tx
	defer func() {
		err = txFinalize(tx1, err)
		err = txFinalize(tx2, err)
	}()

	log.Println("Starting transactions")
	tx1, err = pool.BeginContext(ctx)
	if err != nil {
		return
	}
	tx2, err = pool.BeginContext(ctx)
	if err != nil {
		return
	}
	log.Println("Transactions opened")

	// Kick off two goroutines that contend with each other.
	wg := sync.WaitGroup{}
	errLock := sync.Mutex{}
	wg.Add(2)
	go func() {
		contendErr := contendReads(ctx, &wg, tx1, "hello", "world", cfg)
		errLock.Lock()
		defer errLock.Unlock()
		err = nest(err, contendErr)
	}()
	go func() {
		contendErr := contendReads(ctx, &wg, tx2, "world", "hello", cfg)
		errLock.Lock()
		defer errLock.Unlock()
		err = nest(err, contendErr)
	}()
	wg.Wait()

	if err != nil {
		return
	}

	// Make sure to commit both transactions before moving on.
	err = nest(err, tx1.Commit())
	err = nest(err, tx2.Commit())
	return
}

func main() {
	log.SetFlags(0)
	log.SetOutput(newLogWriter())
	cfg := getConfig()

	// 1. Set the `DB_LOCK_TIMEOUT` environment variable.
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
	// 3. Make sure `lock_timeout` is in the connection string (it gets printed)
	log.Println(separator)
	pool, err := createConn(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanUp(pool)

	// 4. Demonstrate that the observed lock timeout on an open connection is
	//    `LockTimeout`.
	log.Println(separator)
	timeout, err := ensureLockTimeout(ctx, pool, cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("lock_timeout=%s\n", timeout)

	// 5. Create a table schema and insert data to seed the database.
	err = seedDatabase(ctx, pool)
	if err != nil {
		log.Fatal(err)
	}

	// 6. Create two goroutines that intentionally contend with transactions.
	log.Println(separator)
	err = intentionalContention(ctx, pool, cfg)
	if err == nil {
		log.Fatal(ex.New("Expected lock contention to occur"))
	}

	// 7. Display the error / errors in as verbose a way as possible.
	log.Println("***")
	err = displayError(err)
	if err != nil {
		log.Fatal(err)
	}
}
