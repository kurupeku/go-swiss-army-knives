package storage

import (
	"bytes"
	"context"
	"sync"
	"time"
)

// Buffer は同期制御されたバッファを提供します
type Buffer struct {
	mu  sync.Mutex
	buf *bytes.Buffer
}

// Write はバッファにデータを書き込みます
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

// Bytes はバッファ内のバイト列を返します
func (b *Buffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	data := make([]byte, b.buf.Len())
	copy(data, b.buf.Bytes())
	return data
}

// WriteString は文字列をバッファに書き込みます
func (b *Buffer) WriteString(s string) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.WriteString(s)
}

// Len はバッファの長さを返します
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Len()
}

// Reset はバッファをクリアします
func (b *Buffer) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf.Reset()
}

var (
	buf = &Buffer{
		mu:  sync.Mutex{},
		buf: bytes.NewBuffer([]byte{}),
	}
)

// TODO: 引数 ln chan []byte で文字列を受信した際に、グローバル変数 buf *bytes.Buffer へ書き込む
// TODO: ctx context.Context がキャンセルされた場合には速やかに関数を終了する
// TODO: エラーが発生した際には errc chan error へエラーを送信する
func Listen(ctx context.Context, ln chan []byte, errc chan error) {
	// TODO: 1 週目：標準出力（`io.Reader` として受け取る）から出力内容を読み取る処理と、読み取った結果を内部のバッファに保存する処理
}

// TODO: グローバル変数 buf *bytes.Buffer から一定時間ごとに内容を読み込み、内容を引数 out chan []byte へ送信する
// TODO: 読み込む間隔は引数 span time.Duration を利用して制御する
// TODO: buf に何も保存されていなければ内容の送信は行わない
// TODO: 一度に保存された内容すべてを読み取り、 buf にはなにも保存されていない状態にリセットする
// TODO: ctx context.Context がキャンセルされた場合には速やかに関数を終了する
// TODO: エラーが発生した際には errc chan error へエラーを送信する
func Load(ctx context.Context, out chan []byte, errc chan error, span time.Duration) {
	// TODO: 2 週目：内部バッファに保存された内容を一定時間ごとに読み込む処理と、読み取った文字列を Body とした HTTP#POST リクエストを投げる処理
}
