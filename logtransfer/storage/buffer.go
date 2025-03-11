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

// グローバル変数 buf から一定間隔（spanで指定される）でデータを読み込み、空行や空白文字のみの行をスキップして out へ送信する関数
func Load(ctx context.Context, out chan []byte, errc chan error, span time.Duration) {
	// TODO: Implement here

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
