package result

import (
	"fmt"
	"sort"
	"sync"
)

type line struct {
	txt string
	no  int
}

type Result struct {
	sync.Mutex
	Data map[string][]line
}

var GlobalResult = &Result{Data: make(map[string][]line, 100)}

func Set(fileName, txt string, no int) {
	GlobalResult.Lock()
	defer GlobalResult.Unlock()

	if _, ok := GlobalResult.Data[fileName]; !ok {
		GlobalResult.Data[fileName] = make([]line, 0, 10)
	}
	GlobalResult.Data[fileName] = append(GlobalResult.Data[fileName], line{txt, no})
}

func Get() *Result {
	return GlobalResult
}

func RenderWithContent() {
	for i, fName := range GlobalResult.Files() {
		if i > 0 {
			fmt.Print("\n")
		}
		fmt.Println(fName)
		for _, l := range GlobalResult.Data[fName] {
			fmt.Printf("%d: %s\n", l.no, l.txt)
		}
	}
}

func RenderFiles() {
	for _, fName := range GlobalResult.Files() {
		fmt.Println(fName)
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
	GlobalResult = &Result{Data: make(map[string][]line, 100)}
}
