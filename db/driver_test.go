package db

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

// PGXTimestampTest is a test object.
type PGXTimestampTest struct {
	ID uuid.UUID `db:"id,pk"`

	IntValue    int  `db:"int_value"`
	IntPtrValue *int `db:"int_ptr_value"`

	JSONValue PGXjson `db:"json_value,json"`

	Local time.Time
	UTC   time.Time

	LocalWithTimezone    time.Time `db:"local_with_timezone"`
	LocalWithoutTimezone time.Time `db:"local_without_timezone"`

	UTCWithTimezone    time.Time `db:"utc_with_timezone"`
	UTCWithoutTimezone time.Time `db:"utc_without_timezone"`
}

// TableName returns the mapped table name.
func (pgt PGXTimestampTest) TableName() string { return "pgx_timestamp_test" }

// PGXjson is a json object test.
type PGXjson struct {
	Foo string
	Bar string
}

func createPGXTimestampTestTable(conn *Connection, tx *sql.Tx) error {
	_, err := conn.Invoke(OptTx(tx)).Exec(`
	CREATE TABLE pgx_timestamp_test (
		id uuid not null primary key,

		int_value int not null,
		int_ptr_value int,
		json_value jsonb,

		local timestamp not null,
		utc timestamp not null,
		local_with_timezone timestamp with time zone not null,
		local_without_timezone timestamp without time zone not null,
		utc_with_timezone timestamp with time zone not null,
		utc_without_timezone timestamp without time zone not null
	)
	`)
	return err
}

func dropPGXTimestampTestTable(conn *Connection, tx *sql.Tx) error {
	_, err := conn.Invoke(OptTx(tx)).Exec(`DROP TABLE pgx_timestamp_test`)
	return err
}

func Test_PGX_Timestamp(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)

	// create the test table
	assert.Nil(createPGXTimestampTestTable(defaultDB(), tx))
	defer func() { _ = dropPGXTimestampTestTable(defaultDB(), tx) }()

	now := time.Now()
	intValue := 1234

	testObj := PGXTimestampTest{
		ID:          uuid.V4(),
		IntValue:    intValue,
		IntPtrValue: &intValue,
		JSONValue: PGXjson{
			Foo: "foo",
			Bar: "bar",
		},
		Local:                now,
		UTC:                  now.UTC(),
		LocalWithTimezone:    now,
		LocalWithoutTimezone: now,
		UTCWithTimezone:      now.UTC(),
		UTCWithoutTimezone:   now.UTC(),
	}
	assert.Nil(defaultDB().Invoke(OptTx(tx)).Create(&testObj))

	var verify PGXTimestampTest
	found, err := defaultDB().Invoke(OptTx(tx)).Get(&verify, testObj.ID)
	assert.Nil(err)
	assert.True(found)
	assert.True(verify.ID.Equal(testObj.ID), "ID should be equal")

	assert.Equal(intValue, verify.IntValue)
	assert.NotNil(verify.IntPtrValue)
	assert.Equal(intValue, *verify.IntPtrValue)

	assert.Equal("foo", verify.JSONValue.Foo)
	assert.Equal("bar", verify.JSONValue.Bar)

	assertTimeEqual(assert, now, testObj.Local)
	assertTimeEqual(assert, now.UTC(), testObj.Local.UTC())

	assertTimeEqual(assert, now, verify.LocalWithTimezone)
	assertTimeEqual(assert, testObj.LocalWithTimezone, verify.LocalWithTimezone)
	// assertTimeNotEqual(assert, testObj.LocalWithoutTimezone, verify.LocalWithoutTimezone)
	assertTimeEqual(assert, testObj.UTCWithTimezone, verify.UTCWithTimezone)
	assertTimeEqual(assert, testObj.UTCWithoutTimezone, verify.UTCWithoutTimezone)
}

func assertTimeEqual(a *assert.Assertions, expected, actual time.Time) {
	a.InTimeDelta(expected, actual, time.Second, fmt.Sprintf("actual delta: %v", expected.Sub(actual)))
}
