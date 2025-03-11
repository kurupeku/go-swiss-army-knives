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

// 検索結果はこのグローバル変数に保存される
var Store = &Result{Data: make(map[string][]Line, 100)}

// ファイル名、一致した行の内容、行番号を渡すと var Store に保存する関数
func Set(fileName, txt string, no int) {
	Store.Lock()
	defer Store.Unlock()

	if _, ok := Store.Data[fileName]; !ok {
		Store.Data[fileName] = make([]Line, 0, 10)
	}
	Store.Data[fileName] = append(Store.Data[fileName], Line{txt, no})
}

// TODO: ファイル名のみを標準出力に出力する
// TODO: ファイル名は昇順で出力する
// TODO: 標準出力は引数 w io.Writer として渡される想定
func RenderFiles(w io.Writer) {
	files := Store.Files()
	for _, file := range files {
		fmt.Fprintln(w, file)
	}
}

// TODO: ファイル名と一致した行番号、一致した行の標準出力に出力する
// TODO: ファイル名は昇順で出力する
// TODO: 出力フォーマットは README.md を参照
// TODO: 標準出力は引数 w io.Writer として渡される想定
func RenderWithContent(w io.Writer) {
	files := Store.Files()
	for i, file := range files {
		// ファイル名を出力
		fmt.Fprintln(w, file)

		// 一致した行を出力
		lines := Store.Data[file]
		for _, line := range lines {
			fmt.Fprintf(w, "%d: %s\n", line.No, line.Text)
		}

		// 最後のファイル以外は空行を入れる
		if i < len(files)-1 {
			fmt.Fprintln(w)
		}
	}
}

// 保存されているファイル名を昇順でソートした上で []string として返す関数
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
