// Package sterr tries to implement ultimate standard error,
// its very simple but power full for testing mainly
package sterr

import (
	"errors"
	"fmt"
)

// Err is standard error with nesting capability
type Err struct {
	message string
	args    []interface{}
	err     error
}

// New returns error instance with given message
func New(message string) Err {
	return Err{message: message}
}

// Args sets error args that should be used when formatting, copy is returned
func (e Err) Args(args ...interface{}) Err {

	e.args = args
	return e
}

// Wrap wraps error into caller and returns copy or nil if err == nil
func (e Err) Wrap(err error) error {
	if err == nil {
		return nil
	}

	e.err = err
	return e
}

//Unwrap unwraps the error if it is holding any
func (e *Err) Unwrap() error {
	return e.err
}

// Error performs error formating
func (e Err) Error() string {
	message := fmt.Sprintf(e.message, e.args...)
	if e.err == nil {
		return message
	}
	return message + ": " + e.err.Error()
}

// SameSurface compares only first error, does not check recursively
func (e Err) SameSurface(err error) bool {
	if err == nil && e.message == "" {
		return true
	}

	if val, ok := err.(Err); ok {
		return val.message == e.message
	}

	return false
}

// Is returns whether this Err is equal to other Err instance. It will
// check recursively. It ignores arguments and only compares messages.
//
// will panic if e is nil
func (e Err) Is(err error) bool {
	if val, ok := err.(Err); ok && val.message == e.message {
		return errors.Is(e.err, val.err)
	}

	return false
}
