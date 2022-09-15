package types

import (
	"fmt"
	"strings"

	"github.com/francoispqt/gojay"
)

// Interface to convert any type into types.Fields.
type Fielder interface {
	Fields() Fields
}

type Fields map[string]any

func (f Fields) MarshalJSONObject(enc *gojay.Encoder) {
	for k, v := range f {
		switch value := v.(type) {
		case Fielder:
			// Recursive expansion until there are no more items to expand.
			enc.AddObjectKey(k, value.Fields())
		case gojay.MarshalerJSONObject:
			enc.AddObjectKey(k, value)
		case gojay.MarshalerJSONArray:
			enc.AddArrayKey(k, value)
		case error:
			b, err := gojay.MarshalAny(value)
			if err != nil {
				enc.AddStringKey(k, err.Error())
				continue
			}
			if len(b) == 2 && b[0] == '{' && b[1] == '}' {
				enc.AddStringKey(k, value.Error())
				continue
			}
			raw := gojay.EmbeddedJSON(b)
			enc.AddEmbeddedJSONKey(k, &raw)
			enc.AddStringKey(k+"_summary", value.Error())
		default:
			b, err := gojay.MarshalAny(value)
			if err != nil {
				enc.AddStringKey(k, err.Error())
				continue
			}
			raw := gojay.EmbeddedJSON(b)
			enc.AddEmbeddedJSONKey(k, &raw)
		}
	}
}

func (f Fields) MarshalJSON() ([]byte, error) {
	return gojay.MarshalJSONObject(f)
}

func (f Fields) IsNil() bool {
	return f == nil
}

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

If len of Values is odd, an extra value of nil will be appended to the end.

If values are empty, an empty map will be returned.

Example:

	fields := types.NewFields(
		"foo", "bar",
		"baz", 1,
		"qux", 2,
	)
*/
func NewFields(values ...any) Fields {
	length := len(values)
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
