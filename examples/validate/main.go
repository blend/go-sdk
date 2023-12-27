/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package main

import (
	"fmt"
	"time"

	"github.com/blend/go-sdk/ref"
	"github.com/blend/go-sdk/uuid"

	// if you're feeling evil.
	joi "github.com/blend/go-sdk/validate"
)

// Validated is a validated object.
type Validated struct {
	ID       uuid.UUID
	Name     string
	Count    int
	Created  time.Time
	Optional *string
}

// Validate implements validated.
func (v Validated) Validate() error {
	return joi.ReturnFirst(
		joi.Any(v.ID).NotNil(),
		joi.String(&v.Name).Required(),
		joi.String(&v.Name).Matches("foo$"),
		joi.Int(&v.Count).Between(0, 99),
		joi.Any(&v.Count).NotEquals(81),
		joi.Time(&v.Created).BeforeNowUTC(),
		joi.WhenElse(
			func() bool { return v.ID != nil && v.ID.IsV4() },
			joi.String(v.Optional).IsURI(), // not sure why
			joi.String(v.Optional).IsIP(),  // still not sure why
		),
	)
}

func main() {
	objects := []Validated{
		{ID: uuid.V4(), Name: "foo", Count: 55, Created: time.Now().UTC(), Optional: ref.String("https://google.com")},
		{ID: uuid.Zero, Name: "foo", Count: 55, Created: time.Now().UTC(), Optional: ref.String("127.0.0.1")},
	}

	for index, obj := range objects {
		if err := obj.Validate(); err != nil {
			fmt.Printf("object %d fails validation: %v\n", index, joi.ErrFormat(err))
		}
	}
}
