package storage

import (
	"bytes"
	"context"
	"sync"
	"time"
)

// 同期制御を前提としたカスタムバッファ
type Buffer struct {
	sync.Mutex
	buf *bytes.Buffer
}

// バッファに対してスレッドセーフな書き込みを行う
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.Lock()
	defer b.Unlock()

	// データが空の場合は書き込みをスキップ
	if len(p) == 0 {
		return 0, nil
	}

	got := make([]byte, len(p))
	// 最後に改行がなければ追加
	if p[len(p)-1] != '\n' {
		copy(got, p)
		got = append(got, '\n')
	} else {
		copy(got, p)
	}

	return b.buf.Write(got)
}

// バッファに対してスレッドセーフな文字列書き込みを行う
func (b *Buffer) WriteString(s string) (n int, err error) {
	return b.Write([]byte(s))
}

// バッファからスレッドセーフな読み込みを行う
func (b *Buffer) Read() []byte {
	b.Lock()
	defer b.Unlock()

	data := make([]byte, b.buf.Len())
	copy(data, b.buf.Bytes())
	b.buf.Reset()
	return data
}

// 保存されたデータの長さを返す
func (b *Buffer) Len() int {
	b.Lock()
	defer b.Unlock()
	return b.buf.Len()
}

// バッファをリセットする
func (b *Buffer) Reset() {
	b.Lock()
	defer b.Unlock()
	b.buf.Reset()
}

var buf = &Buffer{
	sync.Mutex{},
	bytes.NewBuffer([]byte{}),
}

// ln から受け取ったデータを buf へ書き込む関数
func Listen(ctx context.Context, ln chan []byte, errc chan error) {
	// TODO: Implement here
	// 以下の処理を実装する必要があります：
	// 1. チャネルからのデータ受信
	//    - ln チャネルからデータを継続的に受信する
	//    - チャネルのクローズを適切に検知する
	//
	// 2. バッファへの書き込み
	//    - 受信したデータをスレッドセーフにバッファへ書き込む
	//    - 書き込みエラーを適切に処理する
	//
	// 3. コンテキストによる制御
	//    - キャンセル時は処理を適切に終了する
	//
	// ヒント：
	// - Buffer 構造体のメソッドはスレッドセーフな操作を提供します
	// - チャネルがクロースされたかどうかを検知する方法を調べてみましょう
}

// グローバル変数 buf から一定間隔（spanで指定される）でデータを読み込み、空行や空白文字のみの行をスキップして out へ送信する関数
func Load(ctx context.Context, out chan []byte, errc chan error, span time.Duration) {
	// TODO: Implement here
	// 以下の処理を実装する必要があります：
	// 1. 定期的なデータの読み取り
	//    - 指定された間隔でバッファからデータを読み取る
	//    - バッファが空の場合は適切にスキップする
	//
	// 2. データの検証
	//    - 空行や空白文字のみの行を除外する
	//    - 有効なデータのみを送信対象とする
	//
	// 3. チャネルへの送信
	//    - 検証済みデータを out チャネルへ送信する
	//    - チャネルのクローズを適切に処理する
	//
	// 4. コンテキストによる制御
	//    - 定期的な処理をキャンセル可能にする
	//    - リソースを適切に解放する
	//
	// ヒント：
	// - time パッケージに一定間隔で処理を行うための機能があります
	// - bytes パッケージには文字列の検証に便利な関数があります
}
