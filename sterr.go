// Package sterr tries to implement ultimate standard error,
// its very simple but power full for testing mainly
package sterr

import "fmt"

// Err is standard error with nesting capability
type Err struct {
	message string
	args    []interface{}
	err     error
}

// New returns error instance with given message
func New(message string) *Err {
	return &Err{message: message}
}

// Args sets error args that should be used when formatting, copy is returned
func (e *Err) Args(args ...interface{}) *Err {
	err := *e
	err.args = args
	return &err
}

// Wrap wraps error into caller and returns copy or nil if err == nil
func (e *Err) Wrap(err error) *Err {
	if err == nil {
		return nil
	}

	er := *e
	er.err = err
	return &er
}

//Unwrap unwraps the error if it is holding any
func (e *Err) Unwrap() error {
	return e.err
}

// Error performs error formating
func (e *Err) Error() string {
	message := fmt.Sprintf(e.message, e.args...)
	if e.err == nil {
		return message
	}
	return message + ": " + e.err.Error()
}

// Is returns whether this Err is equal to other Err instance. It will
// check recursively. It ignores arguments and only compares messages
func (e *Err) Is(err error) bool {
	if val, ok := err.(*Err); ok && val.message == e.message {
		if e.err == nil && val.err == nil {
			return true
		}

		if e, ok := e.err.(*Err); ok {
			return e.Is(val.err)
		}
	}

	return false
}
