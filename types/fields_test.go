package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFields(t *testing.T) {
	type args struct {
		values []any
	}
	tests := []struct {
		name string
		args args
		want Fields
	}{
		{
			name: "Ok",
			args: args{
				values: []any{
					"foo", "bar",
					"baz", 123,
					7777, "5000",
					"watata",
				},
			},
			want: map[string]any{
				"foo":    "bar",
				"baz":    123,
				"7777":   "5000",
				"watata": nil,
			},
		},
		{
			name: "Support empty",
			args: args{
				values: []any{},
			},
			want: map[string]any{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFields(tt.args.values...)
			assert.Equal(t, got, tt.want)
		})
	}
}
