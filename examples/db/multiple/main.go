package main

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/db"
	"github.com/blend/go-sdk/db/migration"
	"github.com/blend/go-sdk/logger"
)

type book struct {
	ID   int    `db:"id,pk,auto"`
	Name string `db:"name"`
}

type person struct {
	ID   int    `db:"id,pk,auto"`
	Name string `db:"name"`
}

type ledger struct {
	BookID   int `db:"book_id,pk"`
	PersonID int `db:"person_id,pk"`
}

func createSchema(log logger.Log, conn *db.Connection) error {
	books := migration.NewStep(
		migration.TableNotExists("book"),
		migration.Statements(
			"CREATE TABLE book (id serial not null primary key, name varchar(255))",
		),
	)
	people := migration.NewStep(
		migration.TableNotExists("person"),
		migration.Statements(
			"CREATE TABLE person(id serial not null primary key, name varchar(255))",
		),
	)
	ledger := migration.NewStep(
		migration.TableNotExists("ledger"),
		migration.Statements(
			`CREATE TABLE ledger(
				book_id int not null references book(id),
				person_id int not null references person(id)
			)`,
			`ALTER TABLE ledger ADD PRIMARY KEY (book_id, person_id)`,
		),
	)

	suite := migration.NewWithGroups(
		migration.NewGroupWithActions(
			books,
			people,
			ledger,
		),
	)
	suite.Log = log
	return suite.Apply(context.TODO(), conn)
}

func seedData(log logger.Log, conn *db.Connection) error {
	// seed books
	if err := conn.Invoke().CreateMany([]book{
		{Name: "Old Man and the Sea"},
		{Name: "Romeo & Juliet"},
		{Name: "The Prince"},
		{Name: "1984"},
		{Name: "A Brave New World"},
	}); err != nil {
		return err
	}

	// seed books
	if err := conn.Invoke().CreateMany([]person{
		{Name: "Will"},
		{Name: "Bailey"},
		{Name: "Mike"},
		{Name: "Ayman"},
		{Name: "Madhu"},
	}); err != nil {
		return err
	}

	return conn.Invoke().CreateMany([]ledger{
		{BookID: 1, PersonID: 1},
		{BookID: 2, PersonID: 1},
		{BookID: 3, PersonID: 1},
		{BookID: 1, PersonID: 2},
		{BookID: 4, PersonID: 2},
		{BookID: 1, PersonID: 3},
		{BookID: 5, PersonID: 3},
		{BookID: 1, PersonID: 4},
		{BookID: 2, PersonID: 4},
		{BookID: 3, PersonID: 4},
		{BookID: 4, PersonID: 4},
		{BookID: 5, PersonID: 4},
	})
}

func dropSchema(log logger.Log, conn *db.Connection) error {
	ledger := migration.NewStep(
		migration.TableExists("ledger"),
		migration.Statements(
			"DROP TABLE ledger",
		),
	)
	people := migration.NewStep(
		migration.TableExists("person"),
		migration.Statements(
			"DROP TABLE person",
		),
	)
	books := migration.NewStep(
		migration.TableExists("book"),
		migration.Statements(
			"DROP TABLE book",
		),
	)

	suite := migration.NewWithGroups(
		migration.NewGroupWithActions(
			ledger,
			people,
			books,
		),
	)

	suite.Log = log
	return suite.Apply(context.TODO(), conn)
}

func main() {
	cfg := db.Config{
		Database: "postgres",
		SSLMode:  db.SSLModeDisable,
	}
	conn, err := db.New(db.OptConfig(cfg))
	if err != nil {
		logger.FatalExit(err)
	}
	if err := conn.Open(); err != nil {
		logger.FatalExit(err)
	}

	log := logger.Prod()
	if err := createSchema(log, conn); err != nil {
		logger.FatalExit(err)
	}
	defer dropSchema(log, conn)

	if err := seedData(log, conn); err != nil {
		log.Fatal(err)
		return
	}

	results, err := conn.Connection.Query("select * from book; select * from person; select * from ledger")
	if err != nil {
		log.Fatal(err)
		return
	}

	var b book
	bookColumns := db.Columns(b)
	for results.Next() {
		if err = db.PopulateInOrder(&b, results, bookColumns); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("%#v\n", b)
	}

	if !results.NextResultSet() {
		log.Fatalf("no person result set, cannot continue")
		return
	}

	var p person
	personColumns := db.Columns(p)
	for results.Next() {
		if err = db.PopulateInOrder(&p, results, personColumns); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("%#v\n", p)
	}

	if !results.NextResultSet() {
		log.Fatalf("no ledger result set, cannot continue")
		return
	}

	var l ledger
	ledgerColumns := db.Columns(l)
	for results.Next() {
		if err = db.PopulateInOrder(&l, results, ledgerColumns); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("%#v\n", l)
	}
}
