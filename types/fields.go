package types

import (
	"fmt"
	"strings"
)

type Fields map[string]any

var (
	_ Display       = (*Fields)(nil)
	_ DisplayWriter = (*Fields)(nil)
)

// Human readable output.
func (f Fields) Display() string {
	s := &strings.Builder{}
	f.WriteDisplay(s)
	return s.String()
}

func (f Fields) WriteDisplay(s StringWriter) {
	i := 0
	for k, v := range f {
		if i > 0 {
			_, _ = s.WriteString("\n")
		}
		_, _ = s.WriteString(k)
		_, _ = s.WriteString(": ")
		if v == nil {
			_, _ = s.WriteString("null")
			i++
			continue
		}
		switch value := v.(type) {
		case fmt.Stringer:
			_, _ = s.WriteString(value.String())
		case error:
			_, _ = s.WriteString(value.Error())
		case string:
			if len(value) == 0 {
				_, _ = s.WriteString(`""`)
			} else {
				_, _ = s.WriteString(value)
			}
		case []byte:
			if len(value) == 0 {
				_, _ = s.WriteString(`""`)
				continue
			} else {
				_, _ = s.Write(value)
			}
		case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, complex64, complex128:
			_, _ = s.WriteString(fmt.Sprintf("%v", value))
		default:
			_, _ = s.WriteString("[object]")
		}
		i++
	}
}

/*
Values must be alternating between string and other values, with even index as string.
If value is not of type string, it will be coerced to string using fmt.Sprintf.

Values must be at least of Length of two, otherwise it will be ignored and returns nil.

If len of Values is odd, an extra value of nil will be appended to the end.
*/
func NewFields(values ...any) Fields {
	length := len(values)
	if length < 2 {
		return nil
	}

	if length%2 != 0 {
		values = append(values, nil)
		length += 1
	}

	f := make(Fields, length/2)

	for i := 0; i < length; i += 2 {
		if key, ok := values[i].(string); ok {
			f[key] = values[i+1]
		} else {
			key := fmt.Sprintf("%v", values[i])
			f[key] = values[i+1]
		}
	}
	return f
}
