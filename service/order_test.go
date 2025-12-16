package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLuhn(t *testing.T) {
	type args struct {
		number string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty string",
			args: args{number: ""},
			want: false,
		},
		{
			name: "valid luhn 79927398713",
			args: args{number: "79927398713"},
			want: true,
		},
		{
			name: "valid luhn 4532015112830366",
			args: args{number: "4532015112830366"},
			want: true,
		},
		{
			name: "invalid luhn 79927398710",
			args: args{number: "79927398710"},
			want: false,
		},
		{
			name: "invalid luhn letters",
			args: args{number: "1234abcd567"},
			want: false,
		},
		{
			name: "single digit 0",
			args: args{number: "0"},
			want: false,
		},
		{
			name: "single digit 8",
			args: args{number: "8"},
			want: false,
		},
		{
			name: "valid luhn long number",
			args: args{number: "6011000990139424"},
			want: true,
		},
		{
			name: "number with spaces",
			args: args{number: " 79927398713 "},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, CheckLuhn(tt.args.number), "CheckLuhn(%v)", tt.args.number)
		})
	}
}
