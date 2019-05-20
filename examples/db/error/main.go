package main

import (
	"fmt"
	"log"

	"github.com/blend/go-sdk/db"
)

func main() {

	conn, err := db.New(db.OptConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	conn.Open()

	_, err = conn.Connection.Query("select * from foo")
	fmt.Printf("error: %#v\n", err)
	fmt.Printf("parsed: %#v\n", db.Error(err))
}
