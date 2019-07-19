package main

import (
	"fmt"
	"time"

	"github.com/blend/go-sdk/uuid"
	// if you're feeling evil.
	joi "github.com/blend/go-sdk/validate"
)

var (
	_ joi.Validated = (*Validated)(nil)
)

// Validated is a validated object.
type Validated struct {
	ID      uuid.UUID
	Name    string
	Count   int
	Created time.Time
}

// Validate implements validated.
func (v Validated) Validate() error {
	return joi.First(
		joi.Any(&v.ID).NotNil(),
		joi.String(&v.Name).Matches("foo$"),
		joi.Int(&v.Count).Between(0, 99),
		joi.Any(&v.Count).NotEquals(81),
		joi.Time(&v.Created).BeforeNowUTC(),
	)
}

func main() {
	fmt.Println(joi.Format(Validated{ID: uuid.V4(), Name: "my foo", Count: 10, Created: time.Now().UTC()}.Validate()))
	fmt.Println(joi.Format(Validated{ID: uuid.V4(), Name: "my bar", Count: 10, Created: time.Now().UTC()}.Validate()))
	fmt.Println(joi.Format(Validated{ID: uuid.V4(), Name: "my foo", Count: 102, Created: time.Now().UTC()}.Validate()))
	fmt.Println(joi.Format(Validated{ID: nil, Name: "my foo", Count: 10, Created: time.Now().UTC()}.Validate()))
	fmt.Println(joi.Format(Validated{ID: uuid.V4(), Name: "my foo", Count: 10, Created: time.Now().UTC()}.Validate()))
}
