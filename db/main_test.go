/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/blend/go-sdk/uuid"
)

//------------------------------------------------------------------------------------------------
// Testing Entrypoint
//------------------------------------------------------------------------------------------------

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	conn, err := OpenTestConnection()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}

	setDefaultDB(conn)
	os.Exit(m.Run())
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

// OpenTestConnection opens a test connection from the environment, disabling ssl.
//
// You should not use this function in production like settings, this is why it is kept in the _test.go file.
func OpenTestConnection(opts ...Option) (*Connection, error) {
	defaultOptions := []Option{OptConfigFromEnv(), OptSSLMode(SSLModeDisable)}
	conn, err := Open(New(append(defaultOptions, opts...)...))
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//------------------------------------------------------------------------------------------------
// Util Types
//------------------------------------------------------------------------------------------------

type upsertObj struct {
	UUID		uuid.UUID	`db:"uuid,pk,auto"`
	Timestamp	time.Time	`db:"timestamp_utc"`
	Category	string		`db:"category"`
}

func (uo upsertObj) TableName() string {
	return "upsert_object"
}

func createUpsertObjectTable(tx *sql.Tx) error {
	createSQL := `CREATE TABLE IF NOT EXISTS upsert_object (uuid uuid primary key default gen_random_uuid(), timestamp_utc timestamp, category varchar(255));`
	return IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec(createSQL))
}

type upsertNoAutosObj struct {
	UUID		uuid.UUID	`db:"uuid,pk"`
	Timestamp	time.Time	`db:"timestamp_utc"`
	Category	string		`db:"category"`
}

func (uo upsertNoAutosObj) TableName() string {
	return "upsert_no_autos_object"
}

func createUpsertNoAutosObjectTable(tx *sql.Tx) error {
	createSQL := `CREATE TABLE IF NOT EXISTS upsert_no_autos_object (uuid varchar(255) primary key, timestamp_utc timestamp, category varchar(255));`
	return IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec(createSQL))
}

//------------------------------------------------------------------------------------------------
// Benchmarking
//------------------------------------------------------------------------------------------------

type benchObj struct {
	ID		int		`db:"id,pk,auto"`
	UUID		string		`db:"uuid"`
	Name		string		`db:"name,uk"`
	Timestamp	time.Time	`db:"timestamp_utc"`
	Amount		float32		`db:"amount"`
	Pending		bool		`db:"pending"`
	Category	string		`db:"category"`
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

func createIndex(tx *sql.Tx) error {
	createSQL := `CREATE UNIQUE INDEX ON bench_object (name)`
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
		Name:		fmt.Sprintf("test_object_%d", index),
		UUID:		uuid.V4().String(),
		Timestamp:	time.Now().UTC(),
		Amount:		1000.0 + (5.0 * float32(index)),
		Pending:	index%2 == 0,
		Category:	fmt.Sprintf("category_%d", index),
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
	defer func() { _ = rows.Close() }()
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

type mockTracer struct {
	PrepareHandler		func(context.Context, Config, string)
	QueryHandler		func(context.Context, Config, string, string) TraceFinisher
	FinishPrepareHandler	func(context.Context, error)
	FinishQueryHandler	func(context.Context, sql.Result, error)
}

func (mt mockTracer) Prepare(ctx context.Context, cfg Config, statement string) TraceFinisher {
	if mt.PrepareHandler != nil {
		mt.PrepareHandler(ctx, cfg, statement)
	}
	return mockTraceFinisher{
		FinishPrepareHandler:	mt.FinishPrepareHandler,
		FinishQueryHandler:	mt.FinishQueryHandler,
	}
}

func (mt mockTracer) Query(ctx context.Context, cfg Config, label, statement string) TraceFinisher {
	if mt.PrepareHandler != nil {
		mt.PrepareHandler(ctx, cfg, statement)
	}
	return mockTraceFinisher{
		FinishPrepareHandler:	mt.FinishPrepareHandler,
		FinishQueryHandler:	mt.FinishQueryHandler,
	}
}

type mockTraceFinisher struct {
	FinishPrepareHandler	func(context.Context, error)
	FinishQueryHandler	func(context.Context, sql.Result, error)
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

var (
	_ Tracer = (*captureStatementTracer)(nil)
)

type captureStatementTracer struct {
	Tracer

	Label		string
	Statement	string
	Err		error
}

func (cst *captureStatementTracer) Query(_ context.Context, cfg Config, label string, statement string) TraceFinisher {
	cst.Label = label
	cst.Statement = statement
	return &captureStatementTracerFinisher{cst}
}

type captureStatementTracerFinisher struct {
	*captureStatementTracer
}

func (cstf *captureStatementTracerFinisher) FinishPrepare(context.Context, error)	{}
func (cstf *captureStatementTracerFinisher) FinishQuery(_ context.Context, _ sql.Result, err error) {
	cstf.captureStatementTracer.Err = err
}

var failInterceptorError = "this is just an interceptor error"

func failInterceptor(_ context.Context, _, statement string) (string, error) {
	return "", fmt.Errorf(failInterceptorError)
}

type uniqueObj struct {
	ID	int	`db:"id,pk"`
	Name	string	`db:"name"`
}

// TableName returns the mapped table name.
func (uo uniqueObj) TableName() string {
	return "unique_obj"
}

type uuidTest struct {
	ID	uuid.UUID	`db:"id"`
	Name	string		`db:"name"`
}

func (ut uuidTest) TableName() string {
	return "uuid_test"
}

type EmbeddedTestMeta struct {
	ID		uuid.UUID	`db:"id,pk"`
	TimestampUTC	time.Time	`db:"timestamp_utc"`
}

type embeddedTest struct {
	EmbeddedTestMeta	`db:",inline"`
	Name			string	`db:"name"`
}

func (et embeddedTest) TableName() string {
	return "embedded_test"
}

type jsonTestChild struct {
	Label string `json:"label"`
}

type jsonTest struct {
	ID	int	`db:"id,pk,auto"`
	Name	string	`db:"name"`

	NotNull		jsonTestChild	`db:"not_null,json"`
	Nullable	[]string	`db:"nullable,json"`
}

func (jt jsonTest) TableName() string {
	return "json_test"
}

func secondArgErr(_ interface{}, err error) error {
	return err
}

func createJSONTestTable(tx *sql.Tx) error {
	return IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("create table json_test (id serial primary key, name varchar(255), not_null json, nullable json)"))
}

func dropJSONTextTable(tx *sql.Tx) error {
	return IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("drop table if exists json_test"))
}

func createUpsertAutosRegressionTable(tx *sql.Tx) error {
	schemaDefinition := `CREATE TABLE upsert_auto_regression (
		id uuid not null,
		status smallint not null,
		required boolean not null default false,
		created_at timestamp default current_timestamp,
		updated_at timestamp,
		migrated_at timestamp
	);`
	schemaPrimaryKey := "ALTER TABLE upsert_auto_regression ADD CONSTRAINT pk_upsert_auto_regression_id PRIMARY KEY (id);"
	if _, err := defaultDB().Invoke(OptTx(tx)).Exec(schemaDefinition); err != nil {
		return err
	}
	if _, err := defaultDB().Invoke(OptTx(tx)).Exec(schemaPrimaryKey); err != nil {
		return err
	}
	return nil
}

func dropUpsertRegressionTable(tx *sql.Tx) error {
	_, err := defaultDB().Invoke(OptTx(tx)).Exec("DROP TABLE upsert_auto_regression")
	return err
}

func createUpsertSerialPKTable(tx *sql.Tx) error {
	schemaDefinition := `CREATE TABLE upsert_serial_pk (
		id serial not null primary key,
		status smallint not null,
		required boolean not null default false,
		created_at timestamp default current_timestamp,
		updated_at timestamp,
		migrated_at timestamp
	);`
	if _, err := defaultDB().Invoke(OptTx(tx)).Exec(schemaDefinition); err != nil {
		return err
	}
	return nil
}

func dropUpsertSerialPKTable(tx *sql.Tx) error {
	_, err := defaultDB().Invoke(OptTx(tx)).Exec("DROP TABLE upsert_serial_pk")
	return err
}

// upsertAutoRegression contains all data associated with an envelope of documents.
type upsertAutoRegression struct {
	ID		uuid.UUID	`db:"id,pk"`
	Status		uint8		`db:"status"`
	Required	bool		`db:"required"`
	CreatedAt	*time.Time	`db:"created_at,auto"`
	UpdatedAt	*time.Time	`db:"updated_at,auto"`
	MigratedAt	*time.Time	`db:"migrated_at"`
	ReadOnly	string		`db:"read_only,readonly"`
}

// TableName returns the table name.
func (uar upsertAutoRegression) TableName() string {
	return "upsert_auto_regression"
}

type upsertSerialPK struct {
	ID		int		`db:"id,pk,serial"`
	Status		uint8		`db:"status"`
	Required	bool		`db:"required"`
	CreatedAt	*time.Time	`db:"created_at,auto"`
	UpdatedAt	*time.Time	`db:"updated_at,auto"`
	MigratedAt	*time.Time	`db:"migrated_at"`
	ReadOnly	string		`db:"read_only,readonly"`
}

// TableName returns the table name.
func (uar upsertSerialPK) TableName() string {
	return "upsert_serial_pk"
}
