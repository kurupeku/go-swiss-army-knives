package storage

import (
	"bytes"
	"context"
	"time"
)

var (
	buf = bytes.NewBuffer([]byte{})
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
