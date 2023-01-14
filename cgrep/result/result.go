package result

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
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
	max  int
}

var r = &result{data: make(map[string][]line, 100)}

func Set(fileName, txt string, no int) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.data[fileName]; !ok {
		r.data[fileName] = make([]line, 0, 10)
	}
	r.data[fileName] = append(r.data[fileName], line{txt, no})

	noLen := len(strconv.Itoa(no))
	if noLen > r.max {
		r.max = noLen
	}
}

func Get() *result {
	return r
}

func RenderWithContent() {
	for i, fName := range r.files() {
		if i > 0 {
			fmt.Print("\n")
		}
		fmt.Println(fName)
		for _, l := range r.data[fName] {
			fmt.Printf(r.paddingTemplate(), l.no, l.txt)
		}
	}
}

func RenderFiles() {
	for _, fName := range r.files() {
		fmt.Println(fName)
	}
}

func (r *result) files() []string {
	files := make([]string, 0, len(r.data))
	for k := range r.data {
		files = append(files, k)
	}

	sort.Strings(files)
	return files
}

func (r *result) paddingTemplate() string {
	return strings.Join([]string{"%", strconv.Itoa(r.max), "d: %s\n"}, "")
}
