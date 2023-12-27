/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package db

import (
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/uuid"
)

func Test_Query_OutMany(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var all []benchObj
	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").OutMany(&all)
	its.Nil(err)
	its.NotEmpty(all)
}

func Test_Query_OutMany_statementInterceptorFailure(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var all []benchObj
	err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Query("select * from bench_object").OutMany(&all)
	its.Equal(failInterceptorError, err.Error())
	its.Empty(all)
}

func Test_Query_Out(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var out benchObj
	_, err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object limit 1").Out(&out)
	its.Nil(err)
	its.NotZero(out.ID)
}

func Test_Query_Out_statementInterceptorFailure(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var out benchObj
	_, err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Query("select * from bench_object limit 1").Out(&out)
	its.Equal(failInterceptorError, err.Error())
	its.Zero(out.ID)
}

func Test_Query_Do(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	rows, err := defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Do()
	its.Nil(err)
	defer rows.Close()
	its.True(rows.Next())
	its.Nil(rows.Err())
}

func Test_Query_Do_StatementInterceptorFailure(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	_, err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Query("select * from bench_object").Do()
	its.Equal(failInterceptorError, err.Error())
}

func Test_Query_Each(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var all []benchObj
	var popErr error
	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Each(func(r Rows) error {
		bo := benchObj{}
		popErr = bo.Populate(r)
		if popErr != nil {
			return popErr
		}
		all = append(all, bo)
		return nil
	})
	its.Nil(err)
	its.NotEmpty(all)
}

func Test_Query_Each_statementInterceptorFailure(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var all []benchObj
	var popErr error
	err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Query("select * from bench_object").Each(func(r Rows) error {
		bo := benchObj{}
		popErr = bo.Populate(r)
		if popErr != nil {
			return popErr
		}
		all = append(all, bo)
		return nil
	})
	its.Equal(failInterceptorError, err.Error())
	its.Empty(all)
}

func Test_Query_Any(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = seedObjects(10, tx)
	its.Nil(err)

	var all []benchObj
	allErr := defaultDB().Invoke(OptTx(tx)).All(&all)
	its.Nil(allErr)
	its.NotEmpty(all)

	obj := all[0]

	exists, err := defaultDB().Invoke(OptTx(tx)).Query("select 1 from bench_object where id = $1", obj.ID).Any()
	its.Nil(err)
	its.True(exists)

	notExists, err := defaultDB().Invoke(OptTx(tx)).Query("select 1 from bench_object where id = $1", -1).Any()
	its.Nil(err)
	its.False(notExists)
}

func Test_Query_Any_statementInterceptorFailure(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = seedObjects(10, tx)
	its.Nil(err)

	var all []benchObj
	allErr := defaultDB().Invoke(OptTx(tx)).All(&all)
	its.Nil(allErr)
	its.NotEmpty(all)

	obj := all[0]

	exists, err := defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Query("select 1 from bench_object where id = $1", obj.ID).Any()
	its.Equal(failInterceptorError, err.Error())
	its.False(exists)
}

func Test_Query_None(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var all []benchObj
	allErr := defaultDB().Invoke(OptTx(tx)).All(&all)
	its.Nil(allErr)
	its.NotEmpty(all)

	obj := all[0]

	exists, existsErr := defaultDB().Invoke(OptTx(tx)).Query("select 1 from bench_object where id = $1", obj.ID).None()
	its.Nil(existsErr)
	its.False(exists)

	notExists, notExistsErr := defaultDB().Invoke(OptTx(tx)).Query("select 1 from bench_object where id = $1", -1).None()
	its.Nil(notExistsErr)
	its.True(notExists)
}

func Test_Query_None_statementInterceptorFailure(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var all []benchObj
	allErr := defaultDB().Invoke(OptTx(tx)).All(&all)
	its.Nil(allErr)
	its.NotEmpty(all)

	obj := all[0]

	exists, err := defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Query("select 1 from bench_object where id = $1", obj.ID).None()
	its.Equal(failInterceptorError, err.Error())
	its.False(exists)
}

func Test_Query_PanicHandling(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = seedObjects(10, tx)
	its.Nil(err)

	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Each(func(r Rows) error {
		panic("THIS IS A TEST PANIC")
	})
	its.NotNil(err)	// this should have the result of the panic

	// we now test to see if the connection is still in a good state, i.e. that we recovered from the panic
	// and closed the connection / rows / statement
	hasRows, err := defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Any()
	its.Nil(err)
	its.True(hasRows)
}

func Test_Query_Any_MultipleQueriesPerTransaction(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	wg := sync.WaitGroup{}
	wg.Add(3)

	its.NotNil(defaultDB().Connection)

	err = seedObjects(10, nil)
	its.Nil(err)

	go func() {
		defer wg.Done()
		hasRows, err := defaultDB().Query("select * from bench_object").Any()
		its.Nil(err)
		its.True(hasRows)
	}()

	go func() {
		defer wg.Done()
		hasRows, err := defaultDB().Query("select * from bench_object").Any()
		its.Nil(err)
		its.True(hasRows)
	}()

	go func() {
		defer wg.Done()
		hasRows, err := defaultDB().Query("select * from bench_object").Any()
		its.Nil(err)
		its.True(hasRows)
	}()

	wg.Wait()

	hasRows, err := defaultDB().Query("select * from bench_object").Any()
	its.Nil(err)
	its.True(hasRows)
}

func Test_Query_First(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var first benchObj
	var found bool
	found, err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").First(func(r Rows) error {
		return first.Populate(r)
	})
	its.Nil(err)
	its.True(found)
	its.Equal(1, first.ID)
}

func Test_Query_First_statementInterceptorFailure(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var first benchObj
	var found bool
	found, err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Query("select * from bench_object").First(func(r Rows) error {
		return first.Populate(r)
	})
	its.Equal(failInterceptorError, err.Error())
	its.False(found)
	its.Zero(first.ID)
}

func Test_Query_Scan(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var id int
	_, err = defaultDB().Invoke(OptTx(tx)).Query("select id from bench_object limit 1").Scan(&id)
	its.Nil(err)
	its.Equal(1, id)
}

func Test_Query_Scan_statementInterceptorFailure(t *testing.T) {
	its := assert.New(t)
	tx, err := defaultDB().Begin()
	its.Nil(err)
	defer func() { _ = tx.Rollback() }()

	seedErr := seedObjects(10, tx)
	its.Nil(seedErr)

	var id int
	_, err = defaultDB().Invoke(
		OptTx(tx),
		OptInvocationStatementInterceptor(failInterceptor),
	).Query("select id from bench_object limit 1").Scan(&id)
	its.Equal(failInterceptorError, err.Error())
	its.Zero(id)
}

func Test_Query_PopulateByname(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer func() { _ = tx.Rollback() }()

	var first benchObj
	cols := Columns(first)
	_, err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").First(func(r Rows) error {
		return PopulateByName(&first, r, cols)
	})
	assert.Nil(err)
	assert.Equal(1, first.ID)
}

type benchWithPointer struct {
	ID		int		`db:"id,pk,auto"`
	UUID		string		`db:"uuid,nullable,uk"`
	Name		string		`db:"name"`
	Timestamp	*time.Time	`db:"timestamp_utc"`
	Amount		float32		`db:"amount"`
	Pending		bool		`db:"pending"`
	Category	string		`db:"category"`
}

func (t benchWithPointer) TableName() string {
	return "bench_object"
}

func Test_Query_Out_WithDirtyStructs(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer func() { _ = tx.Rollback() }()

	err = createTable(tx)
	assert.Nil(err)

	uniq := uuid.V4().ToFullString()

	i, err := defaultDB().Invoke(OptTx(tx)).Exec("INSERT INTO bench_object (uuid, name, category) VALUES ($1, $2, $3)",
		uniq, "Foo", "Bar")
	assert.Nil(err)
	ra, _ := i.RowsAffected()
	assert.Equal(1, ra)

	timeObj := time.Now()

	dirty := benchWithPointer{
		ID:		192,
		UUID:		uuid.V4().ToFullString(),
		Name:		"Widget",
		Timestamp:	&timeObj,
		Amount:		4.99,
		Category:	"Baz",
	}

	b, err := defaultDB().Invoke(OptTx(tx)).Query("SELECT * FROM bench_object WHERE uuid = $1", uniq).Out(&dirty)
	assert.Nil(err)
	assert.True(b)
	assert.Nil(dirty.Timestamp)
	assert.True(dirty.Amount == 0)
}
