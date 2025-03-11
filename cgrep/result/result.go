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

// Store に保存されているファイル名のみを出力する関数
func RenderFiles(w io.Writer) {
	// TODO: Implement here
	// 以下の処理を実装する必要があります：
	// 1. ファイル名の取得
	//    - 検索結果として保存されているファイル名の一覧を取得する
	//
	// 2. 出力処理
	//    - 取得したファイル名を表示する
	//    - ファイル名は昇順で表示する
	//    - フォーマットの詳細は README.md を参照
	//
	// ヒント：
	// - ファイル名をソート済みのスライスで取得できるメソッドが用意されています

	files := Store.Files()
	for _, file := range files {
		fmt.Fprintln(w, file)
	}
}

// Store に保存されているファイル名と一致した行の内容、行番号を出力する関数
func RenderWithContent(w io.Writer) {
	// TODO: Implement here
	// 以下の処理を実装する必要があります：
	// 1. ファイル名の取得
	//    - 検索結果として保存されているファイル名の一覧を取得する
	//
	// 2. ファイルごとの出力処理
	//    - ファイル名に加えて、そのファイルでマッチした各行の内容も出力する
	//    - ファイル名は昇順で表示する
	//    - フォーマットの詳細は README.md を参照
	//
	// ヒント：
	// - ファイル名をソート済みのスライスで取得できるメソッドが用意されています
	// - 結果をどのように保存しているかを他の実装から読み解きましょう

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

// 保存されている検索結果をリセットする関数
func Reset() {
	Store = &Result{Data: make(map[string][]Line, 100)}
}
