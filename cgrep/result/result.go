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

var Store = &Result{Data: make(map[string][]Line, 100)}

func Set(fileName, txt string, no int) {
	Store.Lock()
	defer Store.Unlock()

	if _, ok := Store.Data[fileName]; !ok {
		Store.Data[fileName] = make([]Line, 0, 10)
	}
	Store.Data[fileName] = append(Store.Data[fileName], Line{txt, no})
}

func Get() *Result {
	return Store
}

func RenderWithContent(w io.Writer) {
	for i, fName := range Store.Files() {
		if i > 0 {
			fmt.Fprintln(w, "")
		}
		fmt.Fprintln(w, fName)
		for _, l := range Store.Data[fName] {
			fmt.Fprintf(w, "%d: %s\n", l.No, l.Text)
		}
	}
}

func RenderFiles(w io.Writer) {
	for _, fName := range Store.Files() {
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
	Store = &Result{Data: make(map[string][]Line, 100)}
}
