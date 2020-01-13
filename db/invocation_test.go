package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
	"time"
	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
)

type jsonTestChild struct {
	Label string `json:"label"`
}

type jsonTest struct {
	ID   int    `db:"id,pk,auto"`
	Name string `db:"name"`

	NotNull  jsonTestChild `db:"not_null,json"`
	Nullable []string      `db:"nullable,json"`
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

func TestInvocationJSONNulls(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()
	defer dropJSONTextTable(tx)

	assert.Nil(createJSONTestTable(tx))

	// try creating fully set object and reading it out
	obj0 := jsonTest{Name: uuid.V4().String(), NotNull: jsonTestChild{Label: uuid.V4().String()}, Nullable: []string{uuid.V4().String()}}
	assert.Nil(defaultDB().Invoke(OptTx(tx)).Create(&obj0))

	var verify0 jsonTest
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify0, obj0.ID)
	assert.Nil(err)

	assert.Equal(obj0.ID, verify0.ID)
	assert.Equal(obj0.Name, verify0.Name)
	assert.Equal(obj0.Nullable, verify0.Nullable)
	assert.Equal(obj0.NotNull.Label, verify0.NotNull.Label)

	// try creating partially set object and reading it out
	obj1 := jsonTest{Name: uuid.V4().String(), NotNull: jsonTestChild{Label: uuid.V4().String()}} //note `Nullable` isn't set

	columns := Columns(obj1)
	values := columns.ColumnValues(obj1)
	assert.Len(values, 4)
	assert.Nil(values[3], "we shouldn't emit a literal 'null' here")
	assert.NotEqual("null", values[3], "we shouldn't emit a literal 'null' here")

	assert.Nil(defaultDB().Invoke(OptTx(tx)).Create(&obj1))

	var verify1 jsonTest
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify1, obj1.ID)
	assert.Nil(err)

	assert.Equal(obj1.ID, verify1.ID)
	assert.Equal(obj1.Name, verify1.Name)
	assert.Nil(verify1.Nullable)
	assert.Equal(obj1.NotNull.Label, verify1.NotNull.Label)

	any, err := defaultDB().Invoke(OptTx(tx)).Query("select 1 from json_test where id = $1 and nullable is null", obj1.ID).Any()
	assert.Nil(err)
	assert.True(any, "we should have written a sql null, not a literal string 'null'")

	// set it to literal 'null' to test this is backward compatible
	err = IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("update json_test set nullable = 'null' where id = $1", obj1.ID))
	assert.Nil(err)

	var verify2 jsonTest
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify2, obj1.ID)
	assert.Nil(err)
	assert.Equal(obj1.ID, verify2.ID)
	assert.Equal(obj1.Name, verify2.Name)
	assert.Nil(verify2.Nullable, "even if we set it to literal 'null' it should come out golang nil")
	assert.Equal(obj1.NotNull.Label, verify2.NotNull.Label)
}

type uniqueObj struct {
	ID   int    `db:"id,pk"`
	Name string `db:"name"`
}

// TableName returns the mapped table name.
func (uo uniqueObj) TableName() string {
	return "unique_obj"
}

func TestInvocationCreateRepeatInTx(t *testing.T) {
	assert := assert.New(t)

	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	assert.Nil(IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("CREATE TABLE IF NOT EXISTS unique_obj (id int not null primary key, name varchar)")))
	assert.Nil(defaultDB().Invoke(OptTx(tx)).Create(&uniqueObj{ID: 1, Name: "one"}))
	var verify uniqueObj
	assert.Nil(secondArgErr(defaultDB().Invoke(OptTx(tx)).Get(&verify, 1)))
	// Make sure it fails if we collide on keys
	assert.NotNil(defaultDB().Invoke(OptTx(tx)).Create(&uniqueObj{ID: 1, Name: "one"}))
	assert.Equal("one", verify.Name)
	assert.NotNil(defaultDB().Invoke(OptTx(tx)).Create(&uniqueObj{ID: 1, Name: "two"}))
}

type uuidTest struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
}

func (ut uuidTest) TableName() string {
	return "uuid_test"
}

func TestInvocationUUIDs(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	assert.Nil(IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("CREATE TABLE IF NOT EXISTS uuid_test (id uuid not null, name varchar(255) not null)")))

	assert.Nil(defaultDB().Invoke(OptTx(tx)).Create(&uuidTest{ID: uuid.V4(), Name: "foo"}))
	assert.Nil(defaultDB().Invoke(OptTx(tx)).Create(&uuidTest{ID: uuid.V4(), Name: "foo2"}))

	var objs []uuidTest
	assert.Nil(defaultDB().Invoke(OptTx(tx)).All(&objs))

	assert.Len(objs, 2)
}

type EmbeddedTestMeta struct {
	ID           uuid.UUID `db:"id,pk"`
	TimestampUTC time.Time `db:"timestamp_utc"`
}

type embeddedTest struct {
	EmbeddedTestMeta `db:",inline"`
	Name             string `db:"name"`
}

func (et embeddedTest) TableName() string {
	return "embedded_test"
}

func TestInlineMeta(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	test := &embeddedTest{EmbeddedTestMeta: EmbeddedTestMeta{ID: uuid.V4(), TimestampUTC: time.Now().UTC()}, Name: "foo"}
	cols := Columns(test)
	assert.NotEmpty(cols.PrimaryKeys().Columns())
	assert.Equal("id", cols.Columns()[0].ColumnName)
	assert.Equal("timestamp_utc", cols.Columns()[1].ColumnName)
	assert.Equal("name", cols.Columns()[2].ColumnName)

	values := cols.NotReadOnly().NotAutos().ColumnValues(test)
	assert.Len(values, 3)
	assert.Equal(test.ID, values[0])
	assert.False(values[1].(time.Time).IsZero())
	assert.Equal("foo", values[2])

	id0 := uuid.V4()
	id1 := uuid.V4()
	assert.Nil(IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("CREATE TABLE IF NOT EXISTS embedded_test (id uuid not null primary key, timestamp_utc timestamp not null, name varchar(255) not null)")))
	assert.Nil(defaultDB().Invoke(OptTx(tx)).Create(&embeddedTest{EmbeddedTestMeta: EmbeddedTestMeta{ID: id0, TimestampUTC: time.Now().UTC()}, Name: "foo"}))
	assert.Nil(defaultDB().Invoke(OptTx(tx)).Create(&embeddedTest{EmbeddedTestMeta: EmbeddedTestMeta{ID: id1, TimestampUTC: time.Now().UTC()}, Name: "foo2"}))

	var objs []embeddedTest
	assert.Nil(defaultDB().Invoke(OptTx(tx)).All(&objs))

	assert.Len(objs, 2)
	assert.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).ID.Equal(id0)
	})
	assert.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).ID.Equal(id1)
	})
	assert.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).Name == "foo"
	})
	assert.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).Name == "foo2"
	})
	assert.All(objs, func(v interface{}) bool {
		return !v.(embeddedTest).TimestampUTC.IsZero()
	})
}

func TestInvocationStatementInterceptor(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	invocation := defaultDB().Invoke(OptInvocationStatementInterceptor(func(statementID, statement string) string {
		return statement + "; -- foo"
	}))
	assert.NotNil(invocation.StatementInterceptor)

	err = IgnoreExecResult(invocation.Exec("select 'ok!'"))
	assert.Nil(err)
}

func TestConnectionCreate(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createTable(tx)
	assert.Nil(err)

	obj := &benchObj{
		Name:      fmt.Sprintf("test_object_0"),
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Amount:    1000.0 + (5.0 * float32(0)),
		Pending:   true,
		Category:  fmt.Sprintf("category_%d", 0),
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(obj)
	assert.Nil(err)
}

func TestConnectionCreateParallel(t *testing.T) {
	assert := assert.New(t)

	err := createTable(nil)
	assert.Nil(err)
	defer dropTableIfExists(nil)

	wg := sync.WaitGroup{}
	wg.Add(5)
	for x := 0; x < 5; x++ {
		go func() {
			defer wg.Done()
			obj := &benchObj{
				Name:      fmt.Sprintf("test_object_0"),
				UUID:      uuid.V4().String(),
				Timestamp: time.Now().UTC(),
				Amount:    1000.0 + (5.0 * float32(0)),
				Pending:   true,
				Category:  fmt.Sprintf("category_%d", 0),
			}
			innerErr := defaultDB().Invoke().Create(obj)
			assert.Nil(innerErr)
		}()
	}
	wg.Wait()
}

func TestConnectionGetMiss(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createUpserObjectTable(tx)
	assert.Nil(err)

	obj := &upsertObj{
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Category:  uuid.V4().String(),
	}
	found, err := defaultDB().Invoke(OptTx(tx)).Get(obj, uuid.V4().String())
	assert.Nil(err)
	assert.False(found)
	assert.Equal("", obj.UUID)
	assert.True(obj.Timestamp.IsZero())
	assert.Equal("", obj.Category)
}

func TestConnectionDelete(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createUpserObjectTable(tx)
	assert.Nil(err)

	obj := &upsertObj{
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Category:  uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(obj)
	assert.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)

	deleted, err := defaultDB().Invoke(OptTx(tx)).Delete(obj)
	assert.Nil(err)
	assert.True(deleted)
}

func TestConnectionDeleteMiss(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createUpserObjectTable(tx)
	assert.Nil(err)

	obj := &upsertObj{
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Category:  uuid.V4().String(),
	}
	deleted, err := defaultDB().Invoke(OptTx(tx)).Delete(obj)
	assert.Nil(err)
	assert.False(deleted)
}

func TestConnectionUpdate(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createUpserObjectTable(tx)
	assert.Nil(err)

	obj := &upsertObj{
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Category:  uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(obj)
	assert.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	updated, err := defaultDB().Invoke(OptTx(tx)).Update(obj)
	assert.Nil(err)
	assert.True(updated)

	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)
}

func TestConnectionUpdateMiss(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createUpserObjectTable(tx)
	assert.Nil(err)

	obj := &upsertObj{
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Category:  uuid.V4().String(),
	}
	updated, err := defaultDB().Invoke(OptTx(tx)).Update(obj)
	assert.Nil(err)
	assert.False(updated)
}

func TestConnectionUpsert(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createUpserObjectTable(tx)
	assert.Nil(err)

	obj := &upsertObj{
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Category:  uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Upsert(obj)
	assert.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	err = defaultDB().Invoke(OptTx(tx)).Upsert(obj)
	assert.Nil(err)

	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)
}

func TestConnectionUpsertWithSerial(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createTable(tx)
	assert.Nil(err)

	obj := &benchObj{
		Name:      "test_object_0",
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Amount:    1005.0,
		Pending:   true,
		Category:  "category_0",
	}
	err = defaultDB().Invoke(OptTx(tx)).Upsert(obj)
	assert.Nil(err, fmt.Sprintf("%+v", err))
	assert.NotZero(obj.ID)

	var verify benchObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.ID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	err = defaultDB().Invoke(OptTx(tx)).Upsert(obj)
	assert.Nil(err)
	assert.NotZero(obj.ID)

	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.ID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)
}

func TestConnectionCreateMany(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createTable(tx)
	assert.Nil(err)

	var objects []DatabaseMapped
	for x := 0; x < 10; x++ {
		objects = append(objects, benchObj{
			Name:      fmt.Sprintf("test_object_%d", x),
			UUID:      uuid.V4().String(),
			Timestamp: time.Now().UTC(),
			Amount:    1005.0,
			Pending:   true,
			Category:  fmt.Sprintf("category_%d", x),
		})
	}

	err = defaultDB().Invoke(OptTx(tx)).CreateMany(objects)
	assert.Nil(err)

	var verify []benchObj
	err = defaultDB().Invoke(OptTx(tx)).Query(`select * from bench_object`).OutMany(&verify)
	assert.Nil(err)
	assert.NotEmpty(verify)
}

func TestConnectionCreateIfNotExists(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	err = createUpserObjectTable(tx)
	assert.Nil(err)

	obj := &upsertObj{
		UUID:      uuid.V4().String(),
		Timestamp: time.Now().UTC(),
		Category:  uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).CreateIfNotExists(obj)
	assert.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	assert.Nil(err)
	assert.Equal(obj.Category, verify.Category)

	oldCategory := obj.Category
	obj.Category = "test"

	err = defaultDB().Invoke(OptTx(tx)).CreateIfNotExists(obj)
	assert.Nil(err)

	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	assert.Nil(err)
	assert.Equal(oldCategory, verify.Category)
}

func TestInvocationMetrics(t *testing.T) {
	assert := assert.New(t)

	log := logger.All(logger.OptOutput(ioutil.Discard))
	defer log.Close()

	done := make(chan struct{})
	var elapsed time.Duration
	log.Listen(QueryFlag, "test", NewQueryEventListener(func(ctx context.Context, qe QueryEvent) {
		elapsed = qe.Elapsed
		close(done)
	}))

	_, err := defaultDB().Invoke(OptInvocationLog(log)).Query("select 'ok!'").Any()
	assert.Nil(err)
	<-done
	assert.NotZero(elapsed)
}
