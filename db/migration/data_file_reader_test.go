package migration

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"testing"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/logger"

	"github.com/blend/go-sdk/assert"
)

func TestDataFileReaderReadLineSingleLine(t *testing.T) {
	assert := assert.New(t)

	fileBuffer := bytes.NewReader([]byte(`this is a test line`))
	dfr := &DataFileReader{}

	var cursor int64
	var readBuffer = make([]byte, 32)
	var readErr error
	var lineBuffer = bytes.NewBuffer([]byte{})

	cursor, readErr = dfr.readLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.NotNil(readErr)
	assert.Equal(io.EOF, readErr)
	assert.Equal(19, cursor)
	assert.Equal(19, lineBuffer.Len())
	assert.Equal("this is a test line", lineBuffer.String())
}

func TestDataFileReaderReadLine(t *testing.T) {
	assert := assert.New(t)

	fileBuffer := bytes.NewReader([]byte(`this is a test line
this is another test line
this is 3rd test line
this is a 4th test line

this is after a blank line
`))
	dfr := &DataFileReader{}

	var cursor int64
	var readBuffer = make([]byte, 32)
	var readErr error
	var lineBuffer = bytes.NewBuffer([]byte{})

	cursor, readErr = dfr.readLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.Nil(readErr)
	assert.Equal(20, cursor)
	assert.Equal(19, lineBuffer.Len())
	assert.Equal("this is a test line", lineBuffer.String())

	lineBuffer.Reset()
	cursor, readErr = dfr.readLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.Nil(readErr)
	assert.Equal(46, cursor)
	assert.Equal(25, lineBuffer.Len())
	assert.Equal("this is another test line", lineBuffer.String())

	lineBuffer.Reset()
	cursor, readErr = dfr.readLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.Nil(readErr)
	assert.Equal(68, cursor)
	assert.Equal(21, lineBuffer.Len())
	assert.Equal("this is 3rd test line", lineBuffer.String())

	lineBuffer.Reset()
	cursor, readErr = dfr.readLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.Nil(readErr)
	assert.Equal(92, cursor)
	assert.Equal(23, lineBuffer.Len())
	assert.Equal("this is a 4th test line", lineBuffer.String())

	lineBuffer.Reset()
	cursor, readErr = dfr.readLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.Nil(readErr, fmt.Sprintf("Total Buffer Len: %d, cusor: %d", fileBuffer.Len(), cursor))
	assert.Equal(93, cursor)
	assert.Equal(0, lineBuffer.Len())
	assert.Empty(lineBuffer.String())

	lineBuffer.Reset()
	cursor, readErr = dfr.readLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.NotNil(readErr)
	assert.Equal(io.EOF, readErr)
	assert.Equal(120, cursor)
	assert.Equal(26, lineBuffer.Len())
	assert.Equal("this is after a blank line", lineBuffer.String())
}

func TestDataFileReaderReadTabLine(t *testing.T) {
	assert := assert.New(t)

	fileBuffer := bytes.NewReader([]byte(`hello world	123	\N	testing
1	2	3	4	5	6	7	8	9	10
this is a line that ends in a tab	`))

	dfr := &DataFileReader{}

	var cursor int64
	var readBuffer = make([]byte, 32)
	var readErr error
	var lineBuffer = bytes.NewBuffer([]byte{})
	var pieces []interface{}

	pieces, cursor, readErr = dfr.readTabLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.Nil(readErr)
	assert.Len(pieces, 4, fmt.Sprintf("%#v", pieces))
	assert.Equal("hello world", pieces[0])
	assert.Equal("123", pieces[1])
	assert.Nil(pieces[2])

	lineBuffer.Reset()
	pieces, cursor, readErr = dfr.readTabLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.Nil(readErr)
	assert.Len(pieces, 10)
	assert.Equal("1", pieces[0])
	assert.Equal("10", pieces[9])

	lineBuffer.Reset()
	pieces, cursor, readErr = dfr.readTabLine(fileBuffer, cursor, readBuffer, lineBuffer)
	assert.NotNil(readErr)
	assert.Equal(io.EOF, readErr)
	assert.Len(pieces, 1)
	assert.Equal("this is a line that ends in a tab", pieces[0])
}

func TestCopyIn(t *testing.T) {
	a := assert.New(t)
	s := CopyIn("test_table", "col_1", "col_2", "col_3", "col_4")
	a.Equal(`COPY "test_table" ("col_1", "col_2", "col_3", "col_4") FROM STDIN`, s)
}

type Data struct {
	Col1 string `db:"col_1"`
	Col2 string `db:"col_2"`
	Col3 string `db:"col_3"`
	Col4 string `db:"col_4"`
}

func TestDataFileReaderAction(t *testing.T) {
	a := assert.New(t)
	conn := getSchemaConnection(db.DefaultSchema, a)
	testSchemaName := buildTestSchemaName()
	err := conn.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", testSchemaName))
	a.Nil(err)
	dfr := ReadDataFile("./testdata/data_file_reader_test.sql")
	s := New(OptLog(logger.None()), OptGroups(createDataFileMigrations(testSchemaName)...))
	defer func() {
		c := conn
		if c == nil {
			c = defaultDB()
		}
		// pq can't parameterize Drop
		err := c.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE;", testSchemaName))
		a.Nil(err)
	}()
	err = s.Apply(context.Background(), conn)
	a.Nil(err)
	app, _, _, _ := s.Results()
	a.Equal(2, app)
	conn.Close()
	conn = getSchemaConnection(testSchemaName, a)

	s = New(OptLog(logger.None()), OptGroups(NewGroupWithActions(NewStep(Always(), dfr.Action))))
	err = s.Apply(context.Background(), conn)
	a.Nil(err)

	var count int
	var found bool
	found, err = conn.Query("Select count(1) from test_table_one").Scan(&count)
	a.Nil(err)
	a.True(found)
	a.Equal(3, count)

	var data []Data
	err = conn.Query("Select * from test_table_one ORDER BY col_1 ASC").OutMany(&data)
	for i, d := range data {
		a.Equal(fmt.Sprintf("data_%d_1", i+1), d.Col1)
		a.Equal(fmt.Sprintf("data_%d_2", i+1), d.Col2)
		a.Equal(fmt.Sprintf("data_%d_3", i+1), d.Col3)
		a.Equal(fmt.Sprintf("data_%d_4", i+1), d.Col4)
	}
}

func createDataFileMigrations(testSchemaName string) []*Group {
	return []*Group{
		NewGroupWithActions(
			NewStep(
				SchemaNotExists(testSchemaName),
				Actions(
					// pq can't parameterize Create
					func(i context.Context, connection *db.Connection, tx *sql.Tx) error {
						err := connection.Exec(fmt.Sprintf("CREATE SCHEMA %s;", testSchemaName))
						if err != nil {
							return err
						}
						return nil
					},
					func(i context.Context, connection *db.Connection, tx *sql.Tx) error {
						// This is a hack to set the schema on the connection
						(&connection.Config).Schema = testSchemaName
						return nil
					},
				))),
		NewGroupWithActions(
			NewStep(
				TableNotExists("test_table_one"),
				Exec(fmt.Sprintf("CREATE TABLE %s.test_table_one (col_1 varchar(16) not null, col_2 varchar(16) not null, col_3 varchar(16) not null, col_4 varchar(16) not null);", testSchemaName)),
			)),
	}
}

func getSchemaConnection(schema string, a *assert.Assertions) *db.Connection {
	var conn *db.Connection
	cfg, err := db.NewConfigFromEnv()
	a.Nil(err)
	(&cfg).Schema = schema
	conn, err = db.New(db.OptConfig(cfg))
	a.Nil(err)
	err = conn.Open()
	a.Nil(err)
	return conn
}
