package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	tests := []struct {
		name      string
		set       []error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "No errors",
			set:       []error{},
			assertion: assert.NoError,
		},
		{
			name: "1 error",
			set: []error{
				errors.New("error1"),
			},
			assertion: assert.Error,
		},
		{
			name: "2 errors",
			set: []error{
				errors.New("error1"),
				errors.New("error2"),
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Reset()

			Store.errs = tt.set
			tt.assertion(t, Error())
		})
	}
}

func TestSetError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *ErrorLogs
	}{
		{
			name: "Success",
			args: args{errors.New("error1")},
			want: &ErrorLogs{
				errs: []error{errors.New("error1")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Reset()

			Set(tt.args.err)
			assert.Equal(t, tt.want, Store)
		})
	}
}

func TestResetError(t *testing.T) {
	tests := []struct {
		name string
		set  []error
		want *ErrorLogs
	}{
		{
			name: "",
			set:  []error{errors.New("error1")},
			want: &ErrorLogs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Store.errs = tt.set
			Reset()

			assert.Equal(t, tt.want, Store)
		})
	}
}
