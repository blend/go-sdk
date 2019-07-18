package main

import (
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/validate"
)

var (
	_ validate.Validated = (*Validated)(nil)
)

// Validated is a validated object.
type Validated struct {
	ID   uuid.UUID
	Name string
	Year int
}

// Validate implements validated.
func (v Validated) Validate() error {
	return validate.Any(
		validate.NotNil(v.ID),
		validate.String.Matches("foo$")(v.Name),
		validate.Int.Between(1997, 2019)(v.Year),
		validate.NotEquals(2001)(v.Year),
	)
}

func main() {
	fmt.Printf("%v\n", Validated{ID: uuid.V4(), Name: "my foo", Year: 2018}.Validate())
	fmt.Printf("%v\n", Validated{ID: uuid.V4(), Name: "my bar", Year: 2018}.Validate())
	fmt.Printf("%v\n", Validated{ID: uuid.V4(), Name: "my foo", Year: 1901}.Validate())
	fmt.Printf("%v\n", Validated{ID: nil, Name: "my foo", Year: 2018}.Validate())
	fmt.Printf("%v\n", Validated{ID: uuid.V4(), Name: "my foo", Year: 2001}.Validate())
}
