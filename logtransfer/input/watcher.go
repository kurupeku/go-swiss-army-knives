package input

import (
	"context"
	"io"
)

// TODO: 引数 r io.Reader とｓて標準出力が渡されてくるので、入力を待ち受ける
// TODO: 入力があった場合は 1 行だけ読み込み、その文字列を引数 ln chan []byte へ送信した後、待受状態に戻る
// TODO: ctx context.Context がキャンセルされた場合には ln を close し、速やかに関数を終了する
// TODO: エラーが発生した際には引数 errc chan error へエラーを送信する
func Monitor(ctx context.Context, ln chan []byte, errc chan error, r io.Reader) {
	// TODO: 1 週目：標準出力（`io.Reader` として受け取る）から出力内容を読み取る処理と、読み取った結果を内部のバッファに保存する処理
}
