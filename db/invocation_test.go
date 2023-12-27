/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package db

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/logger"
	"github.com/blend/go-sdk/uuid"
)

func Test_Invocation_StatementInterceptor(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	statementTracer := new(captureStatementTracer)
	invocation := defaultDB().Invoke(OptInvocationStatementInterceptor(func(_ context.Context, label, statement string) (string, error) {
		return statement + "; -- foo", nil
	}), OptInvocationTracer(statementTracer))
	its.NotNil(invocation.StatementInterceptor)

	err = IgnoreExecResult(invocation.Exec("select 'ok!'"))
	its.Nil(err)

	its.Equal("select 'ok!'; -- foo", statementTracer.Statement)
}

func Test_Invocation_Query(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	its.Equal(DefaultSchema, defaultDB().Config.SchemaOrDefault())
	err = seedObjects(100, tx)
	its.Nil(err)

	objs := []benchObj{}
	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").OutMany(&objs)
	its.Nil(err)
	its.NotEmpty(objs)

	err = defaultDB().Invoke(OptTx(tx), OptInvocationStatementInterceptor(failInterceptor)).Query("select * from bench_object").OutMany(&objs)
	its.Equal("this is just an interceptor error", err.Error())
}

func Test_Invocation_Create(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createTable(tx)
	its.Nil(err)

	obj := &benchObj{
		Name:		"test_object_0",
		UUID:		uuid.V4().String(),
		Timestamp:	time.Now().UTC(),
		Amount:		1000.0 + (5.0 * float32(0)),
		Pending:	true,
		Category:	fmt.Sprintf("category_%d", 0),
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(obj)
	its.Nil(err)
}

func Test_Invocation_Create_statementInterceptorFailure(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createTable(tx)
	its.Nil(err)

	obj := &benchObj{
		Name:		"test_object_0",
		UUID:		uuid.V4().String(),
		Timestamp:	time.Now().UTC(),
		Amount:		1000.0 + (5.0 * float32(0)),
		Pending:	true,
		Category:	fmt.Sprintf("category_%d", 0),
	}
	err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Create(obj)
	its.Equal(failInterceptorError, err.Error())
}

func Test_Invocation_Create_jsonNulls(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()
	defer func() { _ = dropJSONTextTable(tx) }()

	its.Nil(createJSONTestTable(tx))

	// try creating fully set object and reading it out
	obj0 := jsonTest{Name: uuid.V4().String(), NotNull: jsonTestChild{Label: uuid.V4().String()}, Nullable: []string{uuid.V4().String()}}
	its.Nil(defaultDB().Invoke(OptTx(tx)).Create(&obj0))

	var verify0 jsonTest
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify0, obj0.ID)
	its.Nil(err)

	its.Equal(obj0.ID, verify0.ID)
	its.Equal(obj0.Name, verify0.Name)
	its.Equal(obj0.Nullable, verify0.Nullable)
	its.Equal(obj0.NotNull.Label, verify0.NotNull.Label)

	// try creating partially set object and reading it out
	obj1 := jsonTest{Name: uuid.V4().String(), NotNull: jsonTestChild{Label: uuid.V4().String()}}	//note `Nullable` isn't set

	columns := Columns(obj1)
	values := columns.ColumnValues(obj1)
	its.Len(values, 4)
	its.Nil(values[3], "we shouldn't emit a literal 'null' here")
	its.NotEqual("null", values[3], "we shouldn't emit a literal 'null' here")

	its.Nil(defaultDB().Invoke(OptTx(tx)).Create(&obj1))

	var verify1 jsonTest
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify1, obj1.ID)
	its.Nil(err)

	its.Equal(obj1.ID, verify1.ID)
	its.Equal(obj1.Name, verify1.Name)
	its.Nil(verify1.Nullable)
	its.Equal(obj1.NotNull.Label, verify1.NotNull.Label)

	any, err := defaultDB().Invoke(OptTx(tx)).Query("select 1 from json_test where id = $1 and nullable is null", obj1.ID).Any()
	its.Nil(err)
	its.True(any, "we should have written a sql null, not a literal string 'null'")

	// set it to literal 'null' to test this is backward compatible
	err = IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("update json_test set nullable = 'null' where id = $1", obj1.ID))
	its.Nil(err)

	var verify2 jsonTest
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify2, obj1.ID)
	its.Nil(err)
	its.Equal(obj1.ID, verify2.ID)
	its.Equal(obj1.Name, verify2.Name)
	its.Nil(verify2.Nullable, "even if we set it to literal 'null' it should come out golang nil")
	its.Equal(obj1.NotNull.Label, verify2.NotNull.Label)
}

func Test_Invocation_Create_repeatInTx(t *testing.T) {
	its := assert.New(t)

	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	its.Nil(IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("CREATE TABLE IF NOT EXISTS unique_obj (id int not null primary key, name varchar)")))
	its.Nil(defaultDB().Invoke(OptTx(tx)).Create(&uniqueObj{ID: 1, Name: "one"}))
	var verify uniqueObj
	its.Nil(secondArgErr(defaultDB().Invoke(OptTx(tx)).Get(&verify, 1)))
	// Make sure it fails if we collide on keys
	its.NotNil(defaultDB().Invoke(OptTx(tx)).Create(&uniqueObj{ID: 1, Name: "one"}))
	its.Equal("one", verify.Name)
	its.NotNil(defaultDB().Invoke(OptTx(tx)).Create(&uniqueObj{ID: 1, Name: "two"}))
}

func Test_Invocation_Create_uuids(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	its.Nil(IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("CREATE TABLE IF NOT EXISTS uuid_test (id uuid not null, name varchar(255) not null)")))

	its.Nil(defaultDB().Invoke(OptTx(tx)).Create(&uuidTest{ID: uuid.V4(), Name: "foo"}))
	its.Nil(defaultDB().Invoke(OptTx(tx)).Create(&uuidTest{ID: uuid.V4(), Name: "foo2"}))

	var objs []uuidTest
	its.Nil(defaultDB().Invoke(OptTx(tx)).All(&objs))

	its.Len(objs, 2)
}

func Test_Invocation_Create_inlineMeta(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	test := &embeddedTest{EmbeddedTestMeta: EmbeddedTestMeta{ID: uuid.V4(), TimestampUTC: time.Now().UTC()}, Name: "foo"}
	cols := Columns(test)
	its.NotEmpty(cols.PrimaryKeys().Columns())
	its.Equal("id", cols.Columns()[0].ColumnName)
	its.Equal("timestamp_utc", cols.Columns()[1].ColumnName)
	its.Equal("name", cols.Columns()[2].ColumnName)

	values := cols.NotReadOnly().NotAutos().ColumnValues(test)
	its.Len(values, 3)
	its.Equal(test.ID, values[0])
	its.False(values[1].(time.Time).IsZero())
	its.Equal("foo", values[2])

	id0 := uuid.V4()
	id1 := uuid.V4()
	its.Nil(IgnoreExecResult(defaultDB().Invoke(OptTx(tx)).Exec("CREATE TABLE IF NOT EXISTS embedded_test (id uuid not null primary key, timestamp_utc timestamp not null, name varchar(255) not null)")))
	its.Nil(defaultDB().Invoke(OptTx(tx)).Create(&embeddedTest{EmbeddedTestMeta: EmbeddedTestMeta{ID: id0, TimestampUTC: time.Now().UTC()}, Name: "foo"}))
	its.Nil(defaultDB().Invoke(OptTx(tx)).Create(&embeddedTest{EmbeddedTestMeta: EmbeddedTestMeta{ID: id1, TimestampUTC: time.Now().UTC()}, Name: "foo2"}))

	var objs []embeddedTest
	its.Nil(defaultDB().Invoke(OptTx(tx)).All(&objs))

	its.Len(objs, 2)
	its.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).ID.Equal(id0)
	})
	its.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).ID.Equal(id1)
	})
	its.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).Name == "foo"
	})
	its.Any(objs, func(v interface{}) bool {
		return v.(embeddedTest).Name == "foo2"
	})
	its.All(objs, func(v interface{}) bool {
		return !v.(embeddedTest).TimestampUTC.IsZero()
	})
}

func Test_Invocation_Create_parallel(t *testing.T) {
	its := assert.New(t)

	err := createTable(nil)
	its.Nil(err)
	defer func() { _ = dropTableIfExists(nil) }()

	wg := sync.WaitGroup{}
	wg.Add(5)
	for x := 0; x < 5; x++ {
		go func() {
			defer wg.Done()
			obj := &benchObj{
				Name:		"test_object_0",
				UUID:		uuid.V4().String(),
				Timestamp:	time.Now().UTC(),
				Amount:		1000.0 + (5.0 * float32(0)),
				Pending:	true,
				Category:	fmt.Sprintf("category_%d", 0),
			}
			innerErr := defaultDB().Invoke().Create(obj)
			its.Nil(innerErr)
		}()
	}
	wg.Wait()
}

func Test_Invocation_Create_withAutos(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertAutosRegressionTable(tx)
	its.Nil(err)
	defer func() { _ = dropUpsertRegressionTable(tx) }()

	// NOTE; postgres truncates nanos
	ts1 := time.Date(2020, 12, 23, 12, 11, 10, 0, time.UTC)
	ts2 := time.Date(2020, 12, 23, 13, 12, 11, 0, time.UTC)

	// create initial value
	value := upsertAutoRegression{
		ID:		uuid.V4(),
		Status:		1,
		Required:	true,
		// CreatedAt:  &ts0,
		UpdatedAt:	&ts1,
		MigratedAt:	&ts2,
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(&value)
	its.Nil(err)

	var verify upsertAutoRegression
	var found bool
	found, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, value.ID)
	its.Nil(err)
	its.True(found)

	its.Equal(value.Status, verify.Status)
	its.Equal(value.Required, verify.Required)

	its.NotNil(value.CreatedAt)
	its.False(value.CreatedAt.IsZero())
	its.Equal((*value.UpdatedAt).UTC(), (*verify.UpdatedAt).UTC())
	its.Equal((*value.MigratedAt).UTC(), (*verify.MigratedAt).UTC())
}

func Test_Invocation_Get(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	i := defaultDB().Invoke(OptTx(tx))
	err = i.Create(obj)
	its.Nil(err)
	its.Equal("upsert_object_create", i.Label)

	var verify upsertObj
	i = defaultDB().Invoke(OptTx(tx))
	found, err := i.Get(&verify, obj.UUID)
	its.Nil(err)
	its.True(found)
	its.Equal(verify.UUID, obj.UUID)
	its.Equal("upsert_object_get", i.Label)

	// Perform same get, but set a label on the invocation
	verify = upsertObj{}
	i = defaultDB().Invoke(OptTx(tx), OptLabel("bespoke_upsert"))
	found, err = i.Get(&verify, obj.UUID)
	its.Nil(err)
	its.True(found)
	its.Equal(verify.UUID, obj.UUID)
	its.Equal("bespoke_upsert", i.Label)
}

func Test_Invocation_Get_statementInterceptor(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(obj)
	its.Nil(err)

	var verify upsertObj

	found, err := defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Get(&verify, obj.UUID)
	its.Equal(failInterceptorError, err.Error())
	its.False(found)
	its.Empty(verify.UUID)
}

func Test_Invocation_Get_notFound(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	found, err := defaultDB().Invoke(OptTx(tx)).Get(obj, uuid.V4().String())
	its.Nil(err)
	its.False(found)
	its.Empty(obj.UUID)
	its.True(obj.Timestamp.IsZero())
	its.Equal("", obj.Category)
}

func Test_Invocation_Delete(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(obj)
	its.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)

	deleted, err := defaultDB().Invoke(OptTx(tx)).Delete(obj)
	its.Nil(err)
	its.True(deleted)
}

func Test_Invocation_Delete_statementInterceptor(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(obj)
	its.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)

	deleted, err := defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Delete(obj)
	its.Equal(failInterceptorError, err.Error())
	its.False(deleted)
}

func Test_Invocation_Delete_notFound(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	deleted, err := defaultDB().Invoke(OptTx(tx)).Delete(obj)
	its.Nil(err)
	its.False(deleted)
}

func Test_Invocation_Update(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(obj)
	its.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	updated, err := defaultDB().Invoke(OptTx(tx)).Update(obj)
	its.Nil(err)
	its.True(updated)

	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)
}

func Test_Invocation_Update_statementInterceptor(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Create(obj)
	its.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	updated, err := defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Update(obj)
	its.Equal(failInterceptorError, err.Error())
	its.False(updated)
}

func Test_Invocation_Update_notFound(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),	// this will be mostly impossible to exist
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	updated, err := defaultDB().Invoke(OptTx(tx)).Update(obj)
	its.Nil(err)
	its.False(updated)
}

func Test_Invocation_Upsert(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Upsert(obj)
	its.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	err = defaultDB().Invoke(OptTx(tx)).Upsert(obj)
	its.Nil(err)

	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)
}

func Test_Invocation_Upsert_noAutos_logging(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertNoAutosObjectTable(tx)
	its.Nil(err)

	buf := new(bytes.Buffer)
	log := logger.Memory(buf)

	obj := &upsertNoAutosObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationLog(log),
	).Upsert(obj)
	its.Nil(err)

	instances := strings.Count(buf.String(), "upsert_no_autos_object_upsert")
	its.Equal(2, instances, "should have a [db.query.start] and a [db.query]")

	var verify upsertNoAutosObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationLog(log),
	).Upsert(obj)
	its.Nil(err)

	instances = strings.Count(buf.String(), "upsert_no_autos_object_upsert")
	its.Equal(4, instances, "should have a [db.query.start] and a [db.query]")

	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)
}

func Test_Invocation_Upsert_statementInterceptor(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).Upsert(obj)
	its.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)

	obj.Category = "test"

	err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Upsert(obj)
	its.Equal(failInterceptorError, err.Error())
}

func Test_Invocation_Upsert_withAutos(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertAutosRegressionTable(tx)
	its.Nil(err)
	defer func() { _ = dropUpsertRegressionTable(tx) }()

	// NOTE; postgres truncates nanos
	ts0 := time.Date(2020, 12, 23, 11, 10, 9, 0, time.UTC)
	ts1 := time.Date(2020, 12, 23, 12, 11, 10, 0, time.UTC)
	ts2 := time.Date(2020, 12, 23, 13, 12, 11, 0, time.UTC)

	// create initial value
	value := upsertAutoRegression{
		ID:		uuid.V4(),
		Status:		1,
		Required:	true,
		CreatedAt:	&ts0,
		UpdatedAt:	&ts1,
		MigratedAt:	&ts2,
	}
	err = defaultDB().Invoke(OptTx(tx)).Upsert(&value)
	its.Nil(err)

	var verify upsertAutoRegression
	var found bool
	found, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, value.ID)
	its.Nil(err)
	its.True(found)

	its.Equal(value.Status, verify.Status)
	its.Equal(value.Required, verify.Required)

	its.Equal((*value.CreatedAt).UTC(), (*verify.CreatedAt).UTC())
	its.Equal((*value.UpdatedAt).UTC(), (*verify.UpdatedAt).UTC())
	its.Equal((*value.MigratedAt).UTC(), (*verify.MigratedAt).UTC())

	value.CreatedAt = &ts1
	value.UpdatedAt = &ts2
	value.MigratedAt = &ts0

	err = defaultDB().Invoke(OptTx(tx)).Upsert(&value)
	its.Nil(err)

	found, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, value.ID)
	its.Nil(err)
	its.True(found)

	its.Equal(value.Status, verify.Status)
	its.Equal(value.Required, verify.Required)

	its.Equal((*value.CreatedAt).UTC(), (*verify.CreatedAt).UTC())
	its.Equal((*value.UpdatedAt).UTC(), (*verify.UpdatedAt).UTC())
	its.Equal((*value.MigratedAt).UTC(), (*verify.MigratedAt).UTC())
}

func Test_Invocation_Upsert_withAutos_Unset(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertAutosRegressionTable(tx)
	its.Nil(err)
	defer func() { _ = dropUpsertRegressionTable(tx) }()

	tsMig := time.Date(2020, 12, 23, 11, 10, 9, 0, time.UTC)

	// create initial value but let created_at be set by the default
	value := upsertAutoRegression{
		ID:		uuid.V4(),
		Status:		1,
		Required:	true,
		MigratedAt:	&tsMig,
	}

	err = defaultDB().Invoke(OptTx(tx)).Upsert(&value)
	its.Nil(err)

	var verify upsertAutoRegression
	var found bool
	found, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, value.ID)
	its.Nil(err)
	its.True(found)

	its.Equal(value.Status, verify.Status)
	its.Equal(value.Required, verify.Required)
	its.Equal((*value.MigratedAt).UTC(), (*verify.MigratedAt).UTC())
	its.NotNil(verify.CreatedAt)
	recorded := *verify.CreatedAt
	its.True(recorded.Unix() > time.Date(2021, 1, 30, 0, 0, 0, 0, time.UTC).Unix())
}

func Test_Invocation_Upsert_withSerialPKAndOtherAutos(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertSerialPKTable(tx)
	its.Nil(err)
	defer func() { _ = dropUpsertSerialPKTable(tx) }()

	tsMig := time.Date(2020, 12, 23, 11, 10, 9, 0, time.UTC)

	// create initial value but let created_at be set by the default
	value := upsertSerialPK{
		Status:		1,
		Required:	true,
		MigratedAt:	&tsMig,
	}

	err = defaultDB().Invoke(OptTx(tx)).Upsert(&value)
	its.Nil(err)
	its.NotZero(value.ID)
	insertID := value.ID

	var verify upsertSerialPK
	var found bool
	found, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, insertID)
	its.Nil(err)
	its.True(found)

	its.Equal(value.Status, verify.Status)
	its.Equal(value.Required, verify.Required)
	its.Equal((*value.MigratedAt).UTC(), (*verify.MigratedAt).UTC())
	its.NotNil(verify.CreatedAt)
	recorded := *verify.CreatedAt
	its.True(recorded.Unix() > time.Date(2021, 1, 30, 0, 0, 0, 0, time.UTC).Unix())

	newTime := time.Now().UTC()
	value.UpdatedAt = &newTime

	err = defaultDB().Invoke(OptTx(tx)).Upsert(&value)
	its.Nil(err)
	its.NotZero(value.ID)
	its.Equal(insertID, value.ID)

	found, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, insertID)
	its.Nil(err)
	its.True(found)

	// Using Unix so that this doesn't fail due to loss of nanosecond precision in postgres
	its.Equal(newTime.Unix(), verify.UpdatedAt.Unix())
}

func Test_Invocation_CreateMany(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createTable(tx)
	its.Nil(err)

	var objects []DatabaseMapped
	for x := 0; x < 10; x++ {
		objects = append(objects, benchObj{
			Name:		fmt.Sprintf("test_object_%d", x),
			UUID:		uuid.V4().String(),
			Timestamp:	time.Now().UTC(),
			Amount:		1005.0,
			Pending:	true,
			Category:	fmt.Sprintf("category_%d", x),
		})
	}

	err = defaultDB().Invoke(OptTx(tx)).CreateMany(objects)
	its.Nil(err)

	var verify []benchObj
	err = defaultDB().Invoke(OptTx(tx)).Query(`select * from bench_object`).OutMany(&verify)
	its.Nil(err)
	its.NotEmpty(verify)
}

func Test_Invocation_CreateMany_statementInterceptor(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createTable(tx)
	its.Nil(err)

	var objects []DatabaseMapped
	for x := 0; x < 10; x++ {
		objects = append(objects, benchObj{
			Name:		fmt.Sprintf("test_object_%d", x),
			UUID:		uuid.V4().String(),
			Timestamp:	time.Now().UTC(),
			Amount:		1005.0,
			Pending:	true,
			Category:	fmt.Sprintf("category_%d", x),
		})
	}

	err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).CreateMany(objects)
	its.Equal(failInterceptorError, err.Error())
}

func Test_Invocation_UpsertMany(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	its.Nil(createTable(tx))
	its.Nil(createIndex(tx))
	currentTime := time.Now().UTC().Truncate(time.Second)

	// Test using upsertMany for insertion.
	var objects []benchObj
	objects = append(objects, benchObj{
		ID:		1,
		Name:		"test_object",
		UUID:		uuid.V4().ToFullString(),
		Timestamp:	currentTime,
		Amount:		1005.0,
		Pending:	true,
		Category:	"category",
	})
	its.Nil(defaultDB().Invoke(OptTx(tx)).UpsertMany(objects))

	var verify []benchObj
	its.Nil(defaultDB().Invoke(OptTx(tx)).Query(`select * from bench_object`).OutMany(&verify))

	// TODO: Convert the type of Timestamp attribute on benchObj to String, so that the
	//       comparison happens via the value, not the address of the pointer.
	//       Currently asserting the equality of the object always fails since the
	//       comparison of Timestamp field is pointer address comparison.
	its.Equal(len(objects), len(verify))
	its.Equal(objects[0].ID, verify[0].ID)
	its.Equal(objects[0].Name, verify[0].Name)
	its.Equal(objects[0].UUID, verify[0].UUID)
	its.True(objects[0].Timestamp.Equal(verify[0].Timestamp))
	its.Equal(objects[0].Amount, verify[0].Amount)
	its.Equal(objects[0].Pending, verify[0].Pending)
	its.Equal(objects[0].Category, verify[0].Category)

	// Confirm that conflict on uk column, name, results in an update.
	var updatedObjects []benchObj
	updatedObjects = append(updatedObjects, benchObj{
		ID:		1,
		Name:		"test_object",
		UUID:		uuid.V4().ToFullString(),
		Timestamp:	currentTime,
		Amount:		2000,
		Pending:	true,
		Category:	"category",
	})
	its.Nil(defaultDB().Invoke(OptTx(tx)).UpsertMany(updatedObjects))

	var updateVerify []benchObj
	its.Nil(defaultDB().Invoke(OptTx(tx)).Query(`select * from bench_object`).OutMany(&updateVerify))

	its.Equal(len(updatedObjects), len(updateVerify))
	its.Equal(updatedObjects[0].ID, updateVerify[0].ID)
	its.Equal(updatedObjects[0].Name, updateVerify[0].Name)
	its.Equal(updatedObjects[0].UUID, updateVerify[0].UUID)
	its.True(updatedObjects[0].Timestamp.Equal(updateVerify[0].Timestamp))
	its.Equal(updatedObjects[0].Amount, updateVerify[0].Amount)
	its.Equal(updatedObjects[0].Pending, updateVerify[0].Pending)
	its.Equal(updatedObjects[0].Category, updateVerify[0].Category)
}

func Test_Invocation_UpsertMany_statementInterceptor(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	its.Nil(createTable(tx))
	its.Nil(createIndex(tx))
	currentTime := time.Now().UTC().Truncate(time.Second)

	// Test using upsertMany for insertion.
	var objects []benchObj
	objects = append(objects, benchObj{
		ID:		1,
		Name:		"test_object",
		UUID:		uuid.V4().ToFullString(),
		Timestamp:	currentTime,
		Amount:		1005.0,
		Pending:	true,
		Category:	"category",
	})
	err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).UpsertMany(objects)
	its.Equal(failInterceptorError, err.Error())
}

func Test_Invocation_CreateIfNotExists(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(OptTx(tx)).CreateIfNotExists(obj)
	its.Nil(err)

	var verify upsertObj
	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(obj.Category, verify.Category)

	oldCategory := obj.Category
	obj.Category = "test"

	err = defaultDB().Invoke(OptTx(tx)).CreateIfNotExists(obj)
	its.Nil(err)

	_, err = defaultDB().Invoke(OptTx(tx)).Get(&verify, obj.UUID)
	its.Nil(err)
	its.Equal(oldCategory, verify.Category)
}

func Test_Invocation_CreateIfNotExists_statementInterceptor(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createUpsertObjectTable(tx)
	its.Nil(err)

	obj := &upsertObj{
		UUID:		uuid.V4(),
		Timestamp:	time.Now().UTC(),
		Category:	uuid.V4().String(),
	}
	err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).CreateIfNotExists(obj)
	its.Equal(failInterceptorError, err.Error())
}

func Test_Invocation_Exists(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var first benchObj
	_, err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").First(func(r Rows) error {
		return first.Populate(r)
	})
	its.Nil(err)
	its.Equal(1, first.ID)

	exists, err := defaultDB().Invoke(OptTx(tx)).Exists(&first)
	its.Nil(err)
	its.True(exists)

	var invalid benchObj
	exists, err = defaultDB().Invoke(OptTx(tx)).Exists(&invalid)
	its.Nil(err)
	its.False(exists)
}

func Test_Invocation_Exists_statementInterceptor(t *testing.T) {
	t.Parallel()
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var first benchObj
	_, err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Query("select * from bench_object").First(func(r Rows) error {
		return first.Populate(r)
	})
	its.Equal(failInterceptorError, err.Error())
}

func Test_Invocation_metrics(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	log := logger.All(logger.OptOutput(io.Discard))
	defer log.Close()

	done := make(chan struct{})
	var elapsed time.Duration
	log.Listen(QueryFlag, "test", NewQueryEventListener(func(ctx context.Context, qe QueryEvent) {
		elapsed = qe.Elapsed
		close(done)
	}))

	_, err := defaultDB().Invoke(OptInvocationLog(log)).Query("select 'ok!'").Any()
	its.Nil(err)
	<-done
	its.NotZero(elapsed)
}

func Test_Invocation_generateCreateMany(t *testing.T) {
	t.Parallel()
	its := assert.New(t)
	var objects []DatabaseMapped
	for x := 0; x < 10; x++ {
		objects = append(objects, benchObj{
			Name:		fmt.Sprintf("test_object_%d", x),
			UUID:		uuid.V4().String(),
			Timestamp:	time.Now().UTC(),
			Amount:		1005.0,
			Pending:	true,
			Category:	fmt.Sprintf("category_%d", x),
		})
	}
	invocation := defaultDB().Invoke()

	// TODO: Add assertions for the two other values returned.
	queryBody, _, _ := invocation.generateCreateMany(objects)
	its.Equal(
		`INSERT INTO bench_object (uuid,name,timestamp_utc,amount,pending,category) VALUES ($1,$2,$3,$4,$5,$6),($7,$8,$9,$10,$11,$12),($13,$14,$15,$16,$17,$18),($19,$20,$21,$22,$23,$24),($25,$26,$27,$28,$29,$30),($31,$32,$33,$34,$35,$36),($37,$38,$39,$40,$41,$42),($43,$44,$45,$46,$47,$48),($49,$50,$51,$52,$53,$54),($55,$56,$57,$58,$59,$60)`,
		queryBody,
	)
}

func Test_Invocation_generateUpsertMany(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	var objects []DatabaseMapped
	for x := 0; x < 10; x++ {
		objects = append(objects, benchObj{
			Name:		fmt.Sprintf("test_object_%d", x),
			UUID:		uuid.V4().String(),
			Timestamp:	time.Now().UTC(),
			Amount:		1005.0,
			Pending:	true,
			Category:	fmt.Sprintf("category_%d", x),
		})
	}
	invocation := defaultDB().Invoke()

	// TODO: Add assertions for the two other values returned.
	queryBody, _, _ := invocation.generateUpsertMany(objects)
	its.Equal(
		`INSERT INTO bench_object (uuid,name,timestamp_utc,amount,pending,category) VALUES ($1,$2,$3,$4,$5,$6),($7,$8,$9,$10,$11,$12),($13,$14,$15,$16,$17,$18),($19,$20,$21,$22,$23,$24),($25,$26,$27,$28,$29,$30),($31,$32,$33,$34,$35,$36),($37,$38,$39,$40,$41,$42),($43,$44,$45,$46,$47,$48),($49,$50,$51,$52,$53,$54),($55,$56,$57,$58,$59,$60) ON CONFLICT (name) DO UPDATE SET uuid=Excluded.uuid,name=Excluded.name,timestamp_utc=Excluded.timestamp_utc,amount=Excluded.amount,pending=Excluded.pending,category=Excluded.category`,
		queryBody,
	)
}

func Test_Invocation_start(t *testing.T) {
	t.Parallel()
	its := assert.New(t)

	statement, err := new(Invocation).start("test-statement")
	its.Equal(ErrConnectionClosed.Error(), ex.ErrClass(err).Error())
	its.Empty(statement)

	buf := new(bytes.Buffer)
	log := logger.Memory(buf)
	statement, err = (&Invocation{
		DB:		defaultDB().Connection,
		Context:	context.Background(),
		Log:		log,
	}).start("select 1")
	its.Nil(err)
	its.Equal("select 1", statement)
	its.NotEmpty(buf.String())
}
