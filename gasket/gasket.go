// gasket is a module for handling errors with meaningful contexts.

package gasket

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/francoispqt/gojay"
	"github.com/tigorlazuardi/gears-go/types"
)

const NilText = "gasket: nil"

var Nil = errors.New(NilText)

type ErrorWriter interface {
	WriteError(sw io.StringWriter)
}

type CodeHinter interface {
	// Gets the type's code.
	Code() int
}

var (
	_ types.DisplayWriter = (*Error)(nil)
	_ types.Display       = (*Error)(nil)
	_ fmt.Stringer        = (*Error)(nil)
	_ error               = (*Error)(nil)
	_ ErrorWriter         = (*Error)(nil)
	_ CodeHinter          = (*Error)(nil)
)

type Error struct {
	code    int
	message string
	source  error
	context types.Fields
}

func (e Error) MarshalJSONObject(enc *gojay.Encoder) {
	enc.AddIntKey("code", e.code)
	enc.AddStringKey("message", e.message)
	enc.AddInterfaceKey("source", e.source)
	enc.AddObjectKey("context", e.context)
}

func (e Error) IsNil() bool {
	return e.source == nil
}

/*
Wraps error with extra informations.

If error is a type of gasket.Error (more superficially implements types.CodeHinter), it will take the code from the type.
*/
func Wrap(err error) *Error {
	if err == nil {
		err = Nil
	}
	code := 500
	if hint, ok := err.(CodeHinter); ok {
		code = hint.Code()
	}
	return &Error{
		message: err.Error(),
		source:  err,
		code:    code,
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

// Set the code for error (and thus what types.CodeHinter interface will return).
func (e *Error) SetCode(code int) *Error {
	e.code = code
	return e
}

// Set message for current error.
func (e *Error) SetMessage(msg string) *Error {
	e.message = msg
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

// Implements CodeHinter interface.
//
// Returns the Error Code.
func (e Error) Code() int {
	return e.code
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
	inner := e.source.Error()
	if inner != e.message {
		_, _ = sw.WriteString(e.message)
		_, _ = sw.WriteString(" => ")
	}
	if we, ok := e.source.(ErrorWriter); ok {
		we.WriteError(sw)
		return
	}
	_, _ = sw.WriteString(inner)
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

// If somewhere on the error chain, gasket.(*Error) exists,
// and have the given code, returns that gasket.(*Error),
// otherwise returns nil.
func HasCode(err error, code int) *Error {
	if err == nil {
		return nil
	}

	if hint, ok := err.(CodeHinter); ok {
		if hint.Code() != code {
			return HasCode(errors.Unwrap(err), code)
		}
	}

	var ret = &Error{}

	// search for *Error in the chain, and also to ensure that error is gasket.(*Error).
	if errors.As(err, &ret) {
		if ret.Code() == code {
			return ret
		}
		// There is *Error, but does not have the same code.
		//
		// This happens when the initial error does not implements CodeHinter.
		//
		// We look deeper.
		return HasCode(ret.Unwrap(), code)
	}
	return nil
}
