validate
========

`validate` adds basic validator functions that can be composed to form full validation suites.

It takes heavy inspiration from `JOI`, specifically value evaluation (and not so much type enforcement).

Validation faults are returned as exceptions with an outer exception class of `validate.ErrValidation` and a descriptive inner exception for the specific fault.

# Examples

See `/examples/validate/main.go` for a full example.