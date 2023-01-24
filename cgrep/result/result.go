package result

import (
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
	// TODO: 2 週目：検索結果のレンダリング & コマンド実行時のメイン処理の実装
}

// TODO: ファイル名と一致した行番号、一致した行の標準出力に出力する
// TODO: ファイル名は昇順で出力する
// TODO: 出力フォーマットは README.md を参照
// TODO: 標準出力は引数 w io.Writer として渡される想定
func RenderWithContent(w io.Writer) {
	// TODO: 2 週目：検索結果のレンダリング & コマンド実行時のメイン処理の実装
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
