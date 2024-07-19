/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

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
