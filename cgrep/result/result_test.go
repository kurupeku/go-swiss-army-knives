package result

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	type args struct {
		fileName string
		txt      string
		no       int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.args.fileName, tt.args.txt, tt.args.no)
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		want *Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Get())
		})
	}
}

func TestRenderWithContent(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RenderWithContent()
		})
	}
}

func TestRenderFiles(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RenderFiles()
		})
	}
}

func TestResult_Files(t *testing.T) {
	type fields struct {
		Mutex sync.Mutex
		Data  map[string][]Line
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{
				Mutex: tt.fields.Mutex,
				Data:  tt.fields.Data,
			}
			assert.Equal(t, tt.want, r.Files())
		})
	}
}

func TestReset(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
		})
	}
}
