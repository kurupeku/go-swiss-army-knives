package result

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

type Line struct {
	Text string
	No   int
}

type Result struct {
	sync.Mutex
	Data map[string][]Line
}

var GlobalResult = &Result{Data: make(map[string][]Line, 100)}

func Set(fileName, txt string, no int) {
	GlobalResult.Lock()
	defer GlobalResult.Unlock()

	if _, ok := GlobalResult.Data[fileName]; !ok {
		GlobalResult.Data[fileName] = make([]Line, 0, 10)
	}
	GlobalResult.Data[fileName] = append(GlobalResult.Data[fileName], Line{txt, no})
}

func Get() *Result {
	return GlobalResult
}

func RenderWithContent(w io.Writer) {
	for i, fName := range GlobalResult.Files() {
		if i > 0 {
			fmt.Fprintln(w, "")
		}
		fmt.Fprintln(w, fName)
		for _, l := range GlobalResult.Data[fName] {
			fmt.Fprintf(w, "%d: %s\n", l.No, l.Text)
		}
	}
}

func RenderFiles(w io.Writer) {
	for _, fName := range GlobalResult.Files() {
		fmt.Fprintln(w, fName)
	}
}

func (r *Result) Files() []string {
	files := make([]string, 0, len(r.Data))
	for k := range r.Data {
		files = append(files, k)
	}

	sort.Strings(files)
	return files
}

func Reset() {
	GlobalResult = &Result{Data: make(map[string][]Line, 100)}
}
