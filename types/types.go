package types

import (
	"io"
)

type Display interface {
	// Human readable output.
	Display() string
}

// A Writer interface that supports both Write and WriteString.
type StringWriter interface {
	io.Writer
	io.StringWriter
}

type DisplayWriter interface {
	// Collect accumulated Human Readable output to w.
	// Ignore must not be returned, nor does it should panic, because this is just intended for human read.
	//
	// Interface implementer will not return how much bytes written or error.
	//
	// So ensure the Writer is reliable and can consume everything without error, e.g. bytes.Buffer or strings.Builder.
	WriteDisplay(w StringWriter)
}
