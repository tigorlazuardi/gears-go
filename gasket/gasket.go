// gasket is a module for handling errors with meaningful contexts.

package gasket

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tigorlazuardi/gears-go/types"
)

type ErrorWriter interface {
	WriteError(sw io.StringWriter)
}

var (
	_ types.DisplayWriter = (*Error)(nil)
	_ types.Display       = (*Error)(nil)
	_ fmt.Stringer        = (*Error)(nil)
	_ error               = (*Error)(nil)
	_ ErrorWriter         = (*Error)(nil)
)

type Error struct {
	message string
	source  error
	context types.Fields
}

func Wrap(err error, message string) *Error {
	return &Error{
		message: message,
		source:  err,
	}
}

// Human readable output.
func (e Error) Display() string {
	s := &strings.Builder{}
	e.WriteDisplay(s)
	return s.String()
}

// Collect accumulated Human Readable output to s.
// Ignore must not be returned, nor does it should panic, because this is just intended for human read.
//
// Interface implementer will not return how much bytes written or error.
//
// So ensure the Writer is reliable and can consume everything without error, e.g. bytes.Buffer or strings.Builder.
func (e Error) WriteDisplay(s types.StringWriter) {
	if e.source == nil {
		e.source = errors.New("[nil]")
	}
	_, _ = s.WriteString(e.message)
	_, _ = s.WriteString("\n\t")
	switch display := e.source.(type) {
	case types.DisplayWriter:
		display.WriteDisplay(s)
	case types.Display:
		_, _ = s.WriteString(display.Display())
	default:
		_, _ = s.WriteString(e.source.Error())
	}
}

/*
Sets the error context.

Values must be alternating between string and other values, with even index as string.
If value is not of type string, it will be coerced to string using fmt.Sprintf.

Values must be at least of Length of two, otherwise it will be ignored.

Example:

	err.SetContext("key", value, "key2", value2)
*/
func (e *Error) SetContext(values ...any) *Error {
	e.context = types.NewFields(values...)
	return e
}

// Returns the error's message.
func (e Error) Message() string {
	return e.message
}

// Returns the context of the error.
func (e Error) Context() map[string]any {
	return e.context
}

func (e Error) Error() string {
	s := &strings.Builder{}
	e.WriteError(s)
	return s.String()
}

func (e Error) WriteError(sw io.StringWriter) {
	if e.source == nil {
		e.source = errors.New("[nil]")
	}
	_, _ = sw.WriteString(e.message)
	_, _ = sw.WriteString(" => ")
	if we, ok := e.source.(ErrorWriter); ok {
		we.WriteError(sw)
		return
	}
	_, _ = sw.WriteString(e.source.Error())
}

func (e Error) String() string {
	return e.Error()
}

func (e Error) Is(err error) bool {
	return errors.Is(err, e.source)
}

func (e Error) As(err any) bool {
	return errors.As(e.source, err)
}

// Returns the wrapped error.
func (e Error) Unwrap() error {
	return e.source
}
