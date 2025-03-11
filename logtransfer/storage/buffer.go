package storage

import (
	"bytes"
	"context"
	"sync"
	"time"
)

// Buffer は同期制御されたバッファを提供します
type Buffer struct {
	sync.Mutex
	buf *bytes.Buffer
}

// Write はバッファにデータを書き込みます
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

func (b *Buffer) WriteString(s string) (n int, err error) {
	return b.Write([]byte(s))
}

// Read はバッファから内容を読み出し、バッファをクリアします
func (b *Buffer) Read() []byte {
	b.Lock()
	defer b.Unlock()

	data := make([]byte, b.buf.Len())
	copy(data, b.buf.Bytes())
	b.buf.Reset()
	return data
}

// Len はバッファの長さを返します
func (b *Buffer) Len() int {
	b.Lock()
	defer b.Unlock()
	return b.buf.Len()
}

// Reset はバッファをクリアします
func (b *Buffer) Reset() {
	b.Lock()
	defer b.Unlock()
	b.buf.Reset()
}

var buf = &Buffer{
	sync.Mutex{},
	bytes.NewBuffer([]byte{}),
}

// TODO: 引数 ln chan []byte で文字列を受信した際に、グローバル変数 buf *bytes.Buffer へ書き込む
// TODO: ctx context.Context がキャンセルされた場合には速やかに関数を終了する
// TODO: エラーが発生した際には errc chan error へエラーを送信する
func Listen(ctx context.Context, ln chan []byte, errc chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case line, ok := <-ln:
			if !ok {
				return
			}
			if _, err := buf.Write(line); err != nil {
				errc <- err
				return
			}
		}
	}
}

// TODO: グローバル変数 buf *bytes.Buffer から一定時間ごとに内容を読み込み、内容を引数 out chan []byte へ送信する
// TODO: 読み込む間隔は引数 span time.Duration を利用して制御する
// TODO: buf に何も保存されていなければ内容の送信は行わない
// TODO: 一度に保存された内容すべてを読み取り、 buf にはなにも保存されていない状態にリセットする
// TODO: ctx context.Context がキャンセルされた場合には速やかに関数を終了する
// TODO: エラーが発生した際には errc chan error へエラーを送信する
func Load(ctx context.Context, out chan []byte, errc chan error, span time.Duration) {
	ticker := time.NewTicker(span)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if buf.Len() == 0 {
				continue
			}

			content := buf.Read()
			// 空行や空白文字のみの場合はスキップ
			if len(bytes.TrimSpace(content)) > 0 {
				out <- content
			}
		}
	}
}
