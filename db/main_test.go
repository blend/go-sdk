package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
	// tests uses postgres
	_ "github.com/lib/pq"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

//------------------------------------------------------------------------------------------------
// Testing Entrypoint
//------------------------------------------------------------------------------------------------

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	conn, err := New(OptConfigFromEnv())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
	err = openDefaultDB(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
	assert.Main(m)
}

// BenchmarkMain is the benchmarking entrypoint.
func BenchmarkMain(b *testing.B) {
	tx, err := defaultDB().Begin()
	if err != nil {
		b.Error("Unable to create transaction")
		b.FailNow()
	}
	if tx == nil {
		b.Error("`tx` is nil")
		b.FailNow()
	}

	defer func() {
		if tx != nil {
			if err := tx.Rollback(); err != nil {
				b.Errorf("Error rolling back transaction: %v", err)
				b.FailNow()
			}
		}
	}()

	err = seedObjects(10000, tx)
	if err != nil {
		b.Errorf("Error seeding objects: %v", err)
		b.FailNow()
	}

	var manual time.Duration
	for x := 0; x < b.N*10; x++ {
		manualStart := time.Now()
		_, err = readManual(tx)
		if err != nil {
			b.Errorf("Error using manual query: %v", err)
			b.FailNow()
		}
		manual += time.Since(manualStart)
	}

	var orm time.Duration
	for x := 0; x < b.N*10; x++ {
		ormStart := time.Now()
		_, err = readOrm(tx)
		if err != nil {
			b.Errorf("Error using orm: %v", err)
			b.FailNow()
		}
		orm += time.Since(ormStart)
	}

	var ormCached time.Duration
	for x := 0; x < b.N*10; x++ {
		ormCachedStart := time.Now()
		_, err = readCachedOrm(tx)
		if err != nil {
			b.Errorf("Error using orm: %v", err)
			b.FailNow()
		}
		ormCached += time.Since(ormCachedStart)
	}

	b.Logf("Benchmark Test Results:\nManual: %v \nOrm: %v\nOrm (Cached Plan): %v\n", manual, orm, ormCached)
}

//------------------------------------------------------------------------------------------------
// Util Types
//------------------------------------------------------------------------------------------------

type upsertObj struct {
	UUID      string    `db:"uuid,pk"`
	Timestamp time.Time `db:"timestamp_utc"`
	Category  string    `db:"category"`
}

func (uo upsertObj) TableName() string {
	return "upsert_object"
}

func createUpserObjectTable(tx *sql.Tx) error {
	createSQL := `CREATE TABLE IF NOT EXISTS upsert_object (uuid varchar(255) primary key, timestamp_utc timestamp, category varchar(255));`
	return IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec(createSQL))
}

//------------------------------------------------------------------------------------------------
// Benchmarking
//------------------------------------------------------------------------------------------------

type benchObj struct {
	ID        int       `db:"id,pk,auto"`
	UUID      string    `db:"uuid"`
	Name      string    `db:"name"`
	Timestamp time.Time `db:"timestamp_utc"`
	Amount    float32   `db:"amount"`
	Pending   bool      `db:"pending"`
	Category  string    `db:"category"`
}

func (b *benchObj) Populate(rows Scanner) error {
	return rows.Scan(&b.ID, &b.UUID, &b.Name, &b.Timestamp, &b.Amount, &b.Pending, &b.Category)
}

func (b benchObj) TableName() string {
	return "bench_object"
}

func createTable(tx *sql.Tx) error {
	createSQL := `CREATE TABLE IF NOT EXISTS bench_object (
		id serial not null primary key
		, uuid uuid not null
		, name varchar(255)
		, timestamp_utc timestamp
		, amount real
		, pending boolean
		, category varchar(255)
	);`
	return IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec(createSQL))
}

func dropTableIfExists(tx *sql.Tx) error {
	dropSQL := `DROP TABLE IF EXISTS bench_object;`
	return IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec(dropSQL))
}

func ensureUUIDExtension() error {
	uuidCreate := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	return IgnoreExecResult(defaultDB().Exec(uuidCreate))
}

func createObject(index int, tx *sql.Tx) error {
	obj := benchObj{
		Name:      fmt.Sprintf("test_object_%d", index),
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Amount:    1000.0 + (5.0 * float32(index)),
		Pending:   index%2 == 0,
		Category:  fmt.Sprintf("category_%d", index),
	}
	return defaultDB().Invoke(OptTx(tx)).Create(&obj)
}

func seedObjects(count int, tx *sql.Tx) error {
	if err := ensureUUIDExtension(); err != nil {
		return err
	}
	if err := dropTableIfExists(tx); err != nil {
		return err
	}

	if err := createTable(tx); err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		if err := createObject(i, tx); err != nil {
			return err
		}
	}
	return nil
}

func readManual(tx *sql.Tx) ([]benchObj, error) {
	var objs []benchObj
	readSQL := `select id,uuid,name,timestamp_utc,amount,pending,category from bench_object`
	readStmt, err := defaultDB().PrepareContext(context.Background(), readSQL, tx)
	if err != nil {
		return nil, err
	}
	defer readStmt.Close()

	rows, err := readStmt.Query()
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		obj := &benchObj{}
		err = obj.Populate(rows)
		if err != nil {
			return nil, err
		}
		objs = append(objs, *obj)
	}

	return objs, nil
}

func readOrm(tx *sql.Tx) ([]benchObj, error) {
	var objs []benchObj
	allErr := defaultDB().Invoke(OptTx(tx)).Query(fmt.Sprintf("select %s from bench_object", ColumnNamesCSV(benchObj{}))).OutMany(&objs)
	return objs, allErr
}

func readCachedOrm(tx *sql.Tx) ([]benchObj, error) {
	var objs []benchObj
	allErr := defaultDB().Invoke(OptTx(tx), OptLabel("get_all_bench_object")).Query(fmt.Sprintf("select %s from bench_object", ColumnNamesCSV(benchObj{}))).OutMany(&objs)
	return objs, allErr
}

var (
	defaultConnection *Connection
)

func setDefaultDB(conn *Connection) {
	defaultConnection = conn
}

func defaultDB() *Connection {
	return defaultConnection
}

func openDefaultDB(conn *Connection) error {
	err := conn.Open()
	if err != nil {
		return err
	}
	setDefaultDB(conn)
	return nil
}

type mockTracer struct {
	PrepareHandler       func(context.Context, Config, string)
	QueryHandler         func(context.Context, Config, string, string) TraceFinisher
	FinishPrepareHandler func(context.Context, error)
	FinishQueryHandler   func(context.Context, sql.Result, error)
}

func (mt mockTracer) Prepare(ctx context.Context, cfg Config, statement string) TraceFinisher {
	if mt.PrepareHandler != nil {
		mt.PrepareHandler(ctx, cfg, statement)
	}
	return mockTraceFinisher{
		FinishPrepareHandler: mt.FinishPrepareHandler,
		FinishQueryHandler:   mt.FinishQueryHandler,
	}
}

func (mt mockTracer) Query(ctx context.Context, cfg Config, label, statement string) TraceFinisher {
	if mt.PrepareHandler != nil {
		mt.PrepareHandler(ctx, cfg, statement)
	}
	return mockTraceFinisher{
		FinishPrepareHandler: mt.FinishPrepareHandler,
		FinishQueryHandler:   mt.FinishQueryHandler,
	}
}

type mockTraceFinisher struct {
	FinishPrepareHandler func(context.Context, error)
	FinishQueryHandler   func(context.Context, sql.Result, error)
}

func (mtf mockTraceFinisher) FinishPrepare(ctx context.Context, err error) {
	if mtf.FinishPrepareHandler != nil {
		mtf.FinishPrepareHandler(ctx, err)
	}
}

func (mtf mockTraceFinisher) FinishQuery(ctx context.Context, res sql.Result, err error) {
	if mtf.FinishQueryHandler != nil {
		mtf.FinishQueryHandler(ctx, res, err)
	}
}
