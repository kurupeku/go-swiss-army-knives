package result

import (
	"bytes"
	"io"
	"sort"
	"strconv"
	"strings"
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
	/* 最終的な標準出力は []byte で渡すため、最初からbyteで処理。
	   配列文字列の単純な連結は join の方が早い？ */
	
	w.Write([]byte(strings.Join(Store.Files(), "\n")))
	w.Write([]byte("\n"))
}

// TODO: ファイル名と一致した行番号、一致した行の標準出力に出力する
// TODO: ファイル名は昇順で出力する
// TODO: 出力フォーマットは README.md を参照
// TODO: 標準出力は引数 w io.Writer として渡される想定
func RenderWithContent(w io.Writer) {
	// TODO: 2 週目：検索結果のレンダリング & コマンド実行時のメイン処理の実装

	byteTmp := &bytes.Buffer{}
	for _, fileName := range Store.Files() {
		byteTmp.Write([]byte(fileName))
		for _, m := range Store.Data[fileName] {
			byteTmp.WriteString("\n")
			byteTmp.WriteString(strconv.Itoa(m.No))
			byteTmp.WriteString(": ")
			byteTmp.WriteString(m.Text)
		}
		byteTmp.WriteString("\n\n")
	}
	w.Write(byteTmp.Bytes()[:len(byteTmp.Bytes())-1]) //普通に作ると、末尾の改行1つが邪魔になるため削除
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
