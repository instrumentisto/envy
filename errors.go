package envigo

import (
	"errors"
	"fmt"
)

// ErrNotStructPtr occurs when an object passed to Parse() method is not
// a pointer to struct, but something else.
var ErrNotStructPtr = errors.New("envigo: expected a pointer to a struct")

// EmptyVarNameError occurs when struct field is tagged with an empty `env` tag.
type EmptyVarNameError struct {
	Field string
}

// Error returns string representation of empty env var name error.
func (e EmptyVarNameError) Error() string {
	return fmt.Sprintf(
		"envigo: env car name cannot be empty on field '%s'", e.Field)
}

// UnparsableTypeError occurs when struct field is tagged with `env` tag,
// but there is no parser for struct field type.
type UnparsableTypeError struct {
	Field string
}

// Error returns string representation of unparsable struct field error.
func (e UnparsableTypeError) Error() string {
	return fmt.Sprintf(
		"envigo: type of field '%s' is not parsable from string", e.Field)
}

// ParseError occurs when parsing from env var value fails.
type ParseError struct {
	Field  string
	EnvVar string
	reason string
}

// Error returns string representation of parsing error.
func (e ParseError) Error() string {
	return fmt.Sprintf(
		"envigo: field '%s' failed to parse from '%s' env var: %s",
		e.Field, e.EnvVar, e.reason)
}
