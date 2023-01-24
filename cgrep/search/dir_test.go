package search

import (
	"cgrep/result"
	"path/filepath"
	"regexp"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDirPath, _     = filepath.Abs("../testdata")
	testFilePath, _    = filepath.Abs("../testdata/text.txt")
	testSubDirPath, _  = filepath.Abs("../testdata/dir")
	testSubFilePath, _ = filepath.Abs("../testdata/dir/text.txt")
	testRegExp1        = regexp.MustCompile("_1")
	testRegExp2        = regexp.MustCompile("_2")
)

type testMockDir struct {
	wg     *sync.WaitGroup
	called bool
}

func (d *testMockDir) Search() {
	defer d.wg.Done()
	d.called = true
}

func TestNew(t *testing.T) {
	type args struct {
		wg       *sync.WaitGroup
		fullPath string
		re       *regexp.Regexp
	}

	wg := new(sync.WaitGroup)

	tests := []struct {
		name      string
		args      args
		want      Dir
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				wg:       wg,
				fullPath: testDirPath,
				re:       testRegExp1,
			},
			want: &dir{
				wg:     wg,
				path:   testDirPath,
				regexp: testRegExp1,
				subDirs: []Dir{
					&dir{
						wg:            wg,
						path:          testSubDirPath,
						regexp:        testRegExp1,
						subDirs:       nil,
						fileFullPaths: []string{testSubFilePath},
					},
				},
				fileFullPaths: []string{testFilePath},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.wg, tt.args.fullPath, tt.args.re)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_dir_Scan(t *testing.T) {
	type fields struct {
		wg            *sync.WaitGroup
		path          string
		regexp        *regexp.Regexp
		subDirs       []Dir
		fileFullPaths []string
	}
	wg := new(sync.WaitGroup)
	tests := []struct {
		name      string
		fields    fields
		want      Dir
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			fields: fields{
				wg:     wg,
				path:   testDirPath,
				regexp: testRegExp1,
			},
			want: &dir{
				wg:     wg,
				path:   testDirPath,
				regexp: testRegExp1,
				subDirs: []Dir{
					&dir{
						wg:            wg,
						path:          testSubDirPath,
						regexp:        testRegExp1,
						subDirs:       nil,
						fileFullPaths: []string{testSubFilePath},
					},
				},
				fileFullPaths: []string{testFilePath},
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dir{
				wg:            wg,
				path:          tt.fields.path,
				regexp:        tt.fields.regexp,
				subDirs:       tt.fields.subDirs,
				fileFullPaths: tt.fields.fileFullPaths,
			}
			tt.assertion(t, d.Scan())
			assert.Equal(t, tt.want, d)
		})
	}
}

func Test_dir_Search(t *testing.T) {
	type fields struct {
		path   string
		regexp *regexp.Regexp
	}
	tests := []struct {
		name   string
		fields fields
		subDir *testMockDir
		setup  func(w *result.Result)
		want   *result.Result
	}{
		{
			name: "matched in root dir",
			fields: fields{
				path:   testDirPath,
				regexp: testRegExp1,
			},
			setup: func(r *result.Result) {
				r.Data = make(map[string][]result.Line, 1)
				r.Data["../testdata/text.txt"] = []result.Line{
					{Text: "sample_text_1-1", No: 1},
					{Text: "  sample_text_1-2", No: 2},
					{Text: "sample_text_1-3", No: 3},
				}
			},
			want: &result.Result{},
		},
		{
			name: "Matched in sub dir",
			fields: fields{
				path:   testDirPath,
				regexp: testRegExp2,
			},
			setup: func(r *result.Result) {
				r.Data = make(map[string][]result.Line, 1)
				r.Data["../testdata/dir/text.txt"] = []result.Line{
					{Text: "sample_text_2-1", No: 1},
					{Text: "sample_text_2-2", No: 2},
				}
			},
			want: &result.Result{},
		},
		{
			name:   "No matches",
			fields: fields{path: testDirPath, regexp: regexp.MustCompile("_3")},
			want:   &result.Result{Data: make(map[string][]result.Line)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer result.Reset()

			wg := new(sync.WaitGroup)
			d, err := New(wg, tt.fields.path, tt.fields.regexp)
			if err != nil {
				t.Fatal(err)
			}

			if tt.setup != nil {
				tt.setup(tt.want)
			}
			mockedSubDir := &testMockDir{wg: wg}
			if dd, ok := d.(*dir); ok {
				dd.subDirs = append(dd.subDirs, mockedSubDir)
			}

			wg.Add(1)
			go d.Search()
			wg.Wait()

			assert.Equal(t, tt.want, result.Store)
			if tt.subDir != nil {
				assert.True(t, tt.subDir.called)
			}
		})
	}
}

func Test_dir_GrepFiles(t *testing.T) {
	type fields struct {
		path          string
		regexp        *regexp.Regexp
		subDirs       []Dir
		fileFullPaths []string
	}
	tests := []struct {
		name      string
		fields    fields
		setup     func(w *result.Result)
		want      *result.Result
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:   "With correct word",
			fields: fields{path: testDirPath, regexp: testRegExp1, subDirs: nil, fileFullPaths: []string{testFilePath}},
			setup: func(r *result.Result) {
				r.Data = make(map[string][]result.Line, 1)
				r.Data["../testdata/text.txt"] = []result.Line{
					{Text: "sample_text_1-1", No: 1},
					{Text: "  sample_text_1-2", No: 2},
					{Text: "sample_text_1-3", No: 3},
				}
			},
			want:      &result.Result{},
			assertion: assert.NoError,
		},
		{
			name:      "With incorrect word",
			fields:    fields{path: testDirPath, regexp: regexp.MustCompile("_2"), subDirs: nil, fileFullPaths: []string{testFilePath}},
			want:      &result.Result{Data: make(map[string][]result.Line)},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer result.Reset()

			d := &dir{
				path:          tt.fields.path,
				regexp:        tt.fields.regexp,
				subDirs:       tt.fields.subDirs,
				fileFullPaths: tt.fields.fileFullPaths,
			}
			if tt.setup != nil {
				tt.setup(tt.want)
			}

			tt.assertion(t, d.GrepFiles())
			assert.Equal(t, tt.want, result.Store)
			result.Reset()
		})
	}
}
