// Package sterr tries to implement ultimate standard error,
// its very simple but power full for testing mainly
package sterr

import (
	"errors"
	"fmt"
	"runtime"
)

// Err is standard error with nesting capability
type Err struct {
	message, file, line string
	args                []interface{}
	trace               StackTrace
	err                 error
}

// New returns error instance with given message
func New(message string) Err {
	return Err{message: message}
}

// Trace builds a stacktrace for the error, length specifies length of stacktrace
// from Trace call
func (e Err) Trace(length int) Err {
	e.trace = NStackTrace(length)
	return e
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

// ReadTrace retrieves trace from Err is given error is instance of it
func ReadTrace(err error) StackTrace {
	if val, ok := err.(Err); ok {
		return val.trace
	}

	return nil
}

// StackTrace stores filenames and lines of stack trace
type StackTrace []StackFrame

// NStackTrace creates new stacktrace, giving info about length of stackframes
func NStackTrace(length int) StackTrace {
	l := length
	if length == -1 {
		l = 20
	} else {
		length += 2
	}
	t := make(StackTrace, 0, l)
	for i := 2; i < length || length == -1; i++ {
		var st StackFrame
		var ok bool
		_, st.Filename, st.Line, ok = runtime.Caller(i)
		if !ok {
			break
		}
		t = append(t, st)
	}

	return t
}

func (s StackTrace) String() (res string) {
	if len(s) == 0 {
		return "no trace tracked"
	}

	for _, s := range s {
		res += s.String()
	}
	return
}

// StackFrame is segment of stacktrace
type StackFrame struct {
	Filename string
	Line     int
}

func (s StackFrame) String() string {
	return fmt.Sprintf("%s:%d\n", s.Filename, s.Line)
}
