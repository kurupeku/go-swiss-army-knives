/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"cgrep/result"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDirPath, _ = filepath.Abs("../testdata")
)

func TestExecSearch(t *testing.T) {
	type args struct {
		fullPath   string
		regexpWord string
	}
	dir = "../testdata"
	defer func() {
		dir = ""
	}()

	tests := []struct {
		name      string
		args      args
		want      *result.Result
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Matched",
			args: args{
				fullPath:   testDirPath,
				regexpWord: `_1\-\d`,
			},
			want: &result.Result{
				Mutex: sync.Mutex{},
				Data: map[string][]result.Line{
					"../testdata/text.txt": {
						{Text: "sample_text_1-1", No: 1},
						{Text: "  sample_text_1-2", No: 2},
						{Text: "sample_text_1-3", No: 3},
					},
				},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer result.Reset()

			tt.assertion(t, ExecSearch(tt.args.fullPath, tt.args.regexpWord))
			assert.Equal(t, tt.want, result.Store)
		})
	}

}

func TestRender(t *testing.T) {
	tests := []struct {
		name        string
		set         *result.Result
		withContent bool
		wantW       string
	}{
		{
			name: "filename only",
			set: &result.Result{
				Mutex: sync.Mutex{},
				Data: map[string][]result.Line{
					"filename1":     {},
					"dir/filename2": {},
				},
			},
			withContent: false,
			wantW:       "dir/filename2\nfilename1\n",
		},
		{
			name: "With content",
			set: &result.Result{
				Mutex: sync.Mutex{},
				Data: map[string][]result.Line{
					"dir/filename1": {
						{Text: "sample_text_1-1", No: 1},
						{Text: "  sample_text_1-2", No: 2},
						{Text: "sample_text_1-3", No: 3},
					},
					"filename2": {
						{Text: "sample_text_2-1", No: 1},
						{Text: "  sample_text_2-2", No: 2},
						{Text: "sample_text_2-3", No: 3},
					},
				},
			},
			withContent: true,
			wantW:       "dir/filename1\n1: sample_text_1-1\n2:   sample_text_1-2\n3: sample_text_1-3\n\nfilename2\n1: sample_text_2-1\n2:   sample_text_2-2\n3: sample_text_2-3\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer result.Reset()
			defer func() {
				withContent = false
			}()

			result.Store = tt.set
			withContent = tt.withContent
			w := &bytes.Buffer{}
			Render(w)
			assert.Equal(t, tt.wantW, w.String())
		})
	}
}
