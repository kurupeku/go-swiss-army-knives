package result

import (
	"fmt"
	"sort"
	"sync"
)

type Result interface {
	Files() string
	WithContent() string
}

type line struct {
	txt string
	no  int
}

type result struct {
	sync.Mutex
	data map[string][]line
}

var r = &result{data: make(map[string][]line, 100)}

func Set(fileName, txt string, no int) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.data[fileName]; !ok {
		r.data[fileName] = make([]line, 0, 10)
	}
	r.data[fileName] = append(r.data[fileName], line{txt, no})
}

func Get() *result {
	return r
}

func RenderWithContent() {
	for i, fName := range r.Files() {
		if i > 0 {
			fmt.Print("\n")
		}
		fmt.Println(fName)
		for _, l := range r.data[fName] {
			fmt.Printf("%d: %s\n", l.no, l.txt)
		}
	}
}

func RenderFiles() {
	for _, fName := range r.Files() {
		fmt.Println(fName)
	}
}

func (r *result) Files() []string {
	files := make([]string, 0, len(r.data))
	for k := range r.data {
		files = append(files, k)
	}

	sort.Strings(files)
	return files
}

func reset() {
	r = &result{data: make(map[string][]line, 100)}
}
