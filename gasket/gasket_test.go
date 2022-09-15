// gasket is a module for handling errors with meaningful contexts.

package gasket

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tigorlazuardi/gears-go/types"
)

func TestHasCode(t *testing.T) {
	type args struct {
		err  error
		code int
	}
	tests := []struct {
		name    string
		args    args
		want    *Error
		wantNil bool
	}{
		{
			name: "success - get the correct error",
			args: args{
				err: Wrap(
					Wrap(
						Wrap(
							Wrap(errors.New("hi")).SetCode(400),
						).
							SetCode(404).
							SetMessage("foo"),
					).SetCode(500),
				),
				code: 404,
			},
			want: &Error{
				code:    404,
				message: "foo",
			},
		},
		{
			name: "success - get the correct error - even when wrapped by other types",
			args: args{
				err: fmt.Errorf("ooo %w",
					Wrap(
						Wrap(
							Wrap(
								Wrap(
									Wrap(errors.New("hi")).SetCode(400).SetMessage("ah"),
								).
									SetCode(404).
									SetMessage("foo"),
							).SetCode(500).SetMessage("wtf"),
						),
					),
				),
				code: 404,
			},
			want: &Error{
				code:    404,
				message: "foo",
			},
		},
		{
			name: "return nil when given error is nil",
			args: args{
				err:  nil,
				code: 404,
			},
			want:    nil,
			wantNil: true,
		},
		{
			name: "return nil when there is no gasket error on the chain",
			args: args{
				err:  fmt.Errorf("foo %w", errors.New("bar")),
				code: 404,
			},
			want:    nil,
			wantNil: true,
		},
		{
			name: "return nil when there is no error with the given code even when there is gasket error on the chain",
			args: args{
				err:  Wrap(Wrap(fmt.Errorf("foo %w", errors.New("bar"))).SetCode(400)).SetCode(500),
				code: 404,
			},
			want:    nil,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasCode(tt.args.err, tt.args.code)
			if tt.wantNil {
				require.Nil(t, got)
				return
			}
			require.NotNil(t, got)

			assert.Equal(t, tt.want.message, got.Message())
			assert.Equal(t, tt.want.code, got.Code())
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		code    int
		message string
		source  error
		context types.Fields
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "expected message",
			fields: fields{
				code:    0,
				message: "top",
				source: Wrap(
					Wrap(
						Wrap(errors.New("hi")),
					).SetCode(400),
				),
				context: map[string]any{},
			},
			want: "top => hi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Error{
				code:    tt.fields.code,
				message: tt.fields.message,
				source:  tt.fields.source,
				context: tt.fields.context,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "expected default error code and message on wrapping error",
			args: args{
				err: Wrap(errors.New("hi")).SetCode(400),
			},
			want: &Error{
				code:    400,
				message: "hi",
				source:  Wrap(errors.New("hi")).SetCode(400),
				context: nil,
			},
		},
		{
			name: "new wrap should only take top level code",
			args: args{
				err: Wrap(Wrap(errors.New("hi")).SetCode(400)).SetCode(404),
			},
			want: &Error{
				code:    404,
				message: "hi",
				source:  Wrap(Wrap(errors.New("hi")).SetCode(400)).SetCode(404),
				context: nil,
			},
		},
		{
			name: "support for nil",
			args: args{
				err: nil,
			},
			want: &Error{
				code:    500,
				message: NilText,
				source:  Nil,
				context: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.args.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
