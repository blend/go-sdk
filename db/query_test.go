package db

import (
	"sync"
	"testing"
	"time"

	"github.com/blend/go-sdk/uuid"

	"github.com/blend/go-sdk/assert"
)

func TestQueryExecute(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer tx.Rollback()

	seedErr := seedObjects(10, tx)
	a.Nil(seedErr)

	rows, err := defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Do()
	a.Nil(err)
	defer rows.Close()
	a.True(rows.Next())
	a.Nil(rows.Err())
}

func TestQueryEach(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer tx.Rollback()

	seedErr := seedObjects(10, tx)
	a.Nil(seedErr)

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
	a.Nil(err)
	a.NotEmpty(all)
}

func TestQueryAny(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer tx.Rollback()

	err = seedObjects(10, tx)
	a.Nil(err)

	var all []benchObj
	allErr := defaultDB().Invoke(OptTx(tx)).All(&all)
	a.Nil(allErr)
	a.NotEmpty(all)

	obj := all[0]

	exists, err := defaultDB().Invoke(OptTx(tx)).Query("select 1 from bench_object where id = $1", obj.ID).Any()
	a.Nil(err)
	a.True(exists)

	notExists, err := defaultDB().Invoke(OptTx(tx)).Query("select 1 from bench_object where id = $1", -1).Any()
	a.Nil(err)
	a.False(notExists)
}

func TestQueryNone(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer tx.Rollback()

	seedErr := seedObjects(10, tx)
	a.Nil(seedErr)

	var all []benchObj
	allErr := defaultDB().Invoke(OptTx(tx)).All(&all)
	a.Nil(allErr)
	a.NotEmpty(all)

	obj := all[0]

	exists, existsErr := defaultDB().Invoke(OptTx(tx)).Query("select 1 from bench_object where id = $1", obj.ID).None()
	a.Nil(existsErr)
	a.False(exists)

	notExists, notExistsErr := defaultDB().Invoke(OptTx(tx)).Query("select 1 from bench_object where id = $1", -1).None()
	a.Nil(notExistsErr)
	a.True(notExists)
}

func TestQueryPanicHandling(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer tx.Rollback()

	err = seedObjects(10, tx)
	a.Nil(err)

	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Each(func(r Rows) error {
		panic("THIS IS A TEST PANIC")
	})
	a.NotNil(err) // this should have the result of the panic

	// we now test to see if the connection is still in a good state, i.e. that we recovered from the panic
	// and closed the connection / rows / statement
	hasRows, err := defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Any()
	a.Nil(err)
	a.True(hasRows)
}

func TestMultipleQueriesPerTransaction(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer tx.Rollback()

	wg := sync.WaitGroup{}
	wg.Add(3)

	a.NotNil(defaultDB().Connection)

	err = seedObjects(10, nil)
	a.Nil(err)

	go func() {
		defer wg.Done()
		hasRows, err := defaultDB().Query("select * from bench_object").Any()
		a.Nil(err)
		a.True(hasRows)
	}()

	go func() {
		defer wg.Done()
		hasRows, err := defaultDB().Query("select * from bench_object").Any()
		a.Nil(err)
		a.True(hasRows)
	}()

	go func() {
		defer wg.Done()
		hasRows, err := defaultDB().Query("select * from bench_object").Any()
		a.Nil(err)
		a.True(hasRows)
	}()

	wg.Wait()

	hasRows, err := defaultDB().Query("select * from bench_object").Any()
	a.Nil(err)
	a.True(hasRows)
}

// Note: this test assumes that `bench_object` DOES NOT EXIST.
// It also is generally skipped as it barfs a bunch of errors into the
// postgres log.
func TestMultipleQueriesPerTransactionWithFailure(t *testing.T) {
	t.Skip()

	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer tx.Rollback()

	wg := sync.WaitGroup{}
	wg.Add(3)

	a.NotNil(defaultDB().Connection)

	go func() {
		defer wg.Done()
		hasRows, err := defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Any()
		a.NotNil(err)
		a.False(hasRows)
	}()

	go func() {
		defer wg.Done()
		hasRows, err := defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Any()
		a.NotNil(err)
		a.False(hasRows)
	}()

	go func() {
		defer wg.Done()
		hasRows, err := defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Any()
		a.NotNil(err)
		a.False(hasRows)
	}()

	wg.Wait()
	hasRows, err := defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").Any()

	a.NotNil(err)
	a.False(hasRows)
}

func TestQueryFirst(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer tx.Rollback()

	seedErr := seedObjects(10, tx)
	a.Nil(seedErr)

	var first benchObj
	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").First(func(r Rows) error {
		return first.Populate(r)
	})
	a.Nil(err)
	a.Equal(1, first.ID)
}

func TestQueryScan(t *testing.T) {
	a := assert.New(t)
	tx, err := defaultDB().Begin()
	a.Nil(err)
	defer tx.Rollback()

	seedErr := seedObjects(10, tx)
	a.Nil(seedErr)

	var first benchObj
	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").First(func(r Rows) error {
		return first.Populate(r)
	})
	a.Nil(err)
	a.Equal(1, first.ID)
}

func TestQueryExists(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	var first benchObj
	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").First(func(r Rows) error {
		return first.Populate(r)
	})
	assert.Nil(err)
	assert.Equal(1, first.ID)

	exists, err := defaultDB().Invoke(OptTx(tx)).Exists(&first)
	assert.Nil(err)
	assert.True(exists)

	var invalid benchObj
	exists, err = defaultDB().Invoke(OptTx(tx)).Exists(&invalid)
	assert.Nil(err)
	assert.False(exists)
}

func TestQueryQueryPopulateByname(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	var first benchObj
	cols := Columns(first)
	err = defaultDB().Invoke(OptTx(tx)).Query("select * from bench_object").First(func(r Rows) error {
		return PopulateByName(&first, r, cols)
	})
	assert.Nil(err)
	assert.Equal(1, first.ID)
}

type benchWithPointer struct {
	ID        int        `db:"id,pk,auto"`
	UUID      string     `db:"uuid,nullable,uk"`
	Name      string     `db:"name"`
	Timestamp *time.Time `db:"timestamp_utc"`
	Amount    float32    `db:"amount"`
	Pending   bool       `db:"pending"`
	Category  string     `db:"category"`
}

func (t benchWithPointer) TableName() string {
	return "bench_object"
}

func TestOutWithDirtyStructs(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

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
		ID:        192,
		UUID:      uuid.V4().ToFullString(),
		Name:      "Widget",
		Timestamp: &timeObj,
		Amount:    4.99,
		Category:  "Baz",
	}

	b, err := defaultDB().Invoke(OptTx(tx)).Query("SELECT * FROM bench_object WHERE uuid = $1", uniq).Out(&dirty)
	assert.Nil(err)
	assert.True(b)
	assert.Nil(dirty.Timestamp)
	assert.True(dirty.Amount == 0)
}

func TestIntoWithDirtyStructs(t *testing.T) {
	assert := assert.New(t)
	tx, err := defaultDB().Begin()
	assert.Nil(err)
	defer tx.Rollback()

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
		ID:        192,
		UUID:      uuid.V4().ToFullString(),
		Name:      "Widget",
		Timestamp: &timeObj,
		Amount:    4.99,
		Category:  "Baz",
	}

	b, err := defaultDB().Invoke(OptTx(tx)).Query("SELECT * FROM bench_object WHERE uuid = $1", uniq).Out(&dirty)
	assert.Nil(err)
	assert.True(b)
	assert.Nil(dirty.Timestamp)
	assert.Zero(dirty.Amount)
}
