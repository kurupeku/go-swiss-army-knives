package result

import (
	"bytes"
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
		want *Result
	}{
		{
			name: "Success",
			args: args{
				fileName: "filename",
				txt:      "text",
				no:       1,
			},
			want: &Result{
				Mutex: sync.Mutex{},
				Data: map[string][]Line{
					"filename": {
						{
							Text: "text",
							No:   1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Reset()

			Set(tt.args.fileName, tt.args.txt, tt.args.no)
			assert.Equal(t, tt.want, Store)
		})
	}
}

func TestRenderWithContent(t *testing.T) {
	tests := []struct {
		name string
		set  *Result
		want string
	}{
		{
			name: "Success1",
			set: &Result{
				Mutex: sync.Mutex{},
				Data: map[string][]Line{
					"filename": {
						{
							Text: "text",
							No:   1,
						},
					},
				},
			},
			want: "filename\n1: text\n",
		},
		{
			name: "Success2",
			set: &Result{
				Mutex: sync.Mutex{},
				Data: map[string][]Line{
					"filename1": {
						{
							Text: "text1",
							No:   1,
						},
						{
							Text: "  text2",
							No:   2,
						},
					},
					"dir/filename2": {
						{
							Text: "text3",
							No:   3,
						},
						{
							Text: "text4",
							No:   4,
						},
					},
				},
			},
			want: "dir/filename2\n3: text3\n4: text4\n\nfilename1\n1: text1\n2:   text2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Reset()

			buf := bytes.NewBuffer([]byte{})
			Store = tt.set
			RenderWithContent(buf)

			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestRenderFiles(t *testing.T) {
	tests := []struct {
		name string
		set  *Result
		want string
	}{
		{
			name: "Success1",
			set: &Result{
				Mutex: sync.Mutex{},
				Data: map[string][]Line{
					"filename": {},
				},
			},
			want: "filename\n",
		},
		{
			name: "Success2",
			set: &Result{
				Mutex: sync.Mutex{},
				Data: map[string][]Line{
					"filename1": {},
					"filename2": {},
					"filename3": {},
				},
			},
			want: "filename1\nfilename2\nfilename3\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer Reset()

			buf := bytes.NewBuffer([]byte{})
			Store = tt.set
			RenderFiles(buf)

			assert.Equal(t, tt.want, buf.String())

		})
	}
}

func TestResult_Files(t *testing.T) {
	type fields struct {
		Data map[string][]Line
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "Success1",
			fields: fields{
				Data: map[string][]Line{},
			},
			want: []string{},
		},
		{
			name: "Success2",
			fields: fields{
				Data: map[string][]Line{
					"c": {},
					"2": {},
				},
			},
			want: []string{"2", "c"},
		},
		{
			name: "Success2",
			fields: fields{
				Data: map[string][]Line{
					"x": {},
					"d": {},
					"o": {},
					"l": {},
				},
			},
			want: []string{"d", "l", "o", "x"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{
				Mutex: sync.Mutex{},
				Data:  tt.fields.Data,
			}
			assert.Equal(t, tt.want, r.Files())
		})
	}
}

func TestReset(t *testing.T) {
	tests := []struct {
		name string
		set  *Result
	}{
		{
			name: "Success",
			set: &Result{
				Mutex: sync.Mutex{},
				Data: map[string][]Line{
					"filename": {
						{
							Text: "text",
							No:   1,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Store = tt.set
			Reset()

			Store = &Result{Data: make(map[string][]Line, 100)}
		})
	}
}
