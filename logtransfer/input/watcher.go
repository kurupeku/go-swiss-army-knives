package input

import (
	"bufio"
	"context"
	"io"
)

func Monitor(ctx context.Context, ln chan []byte, errc chan error, r io.Reader) {
	// TODO: 1 週目：標準出力（`io.Reader` として受け取る）から出力内容を読み取る処理と、読み取った結果を内部のバッファに保存する処理

	// NewScanner returns a new Scanner to read from r.
	scanner := bufio.NewScanner(r)

	// 無限ループ
	for {
		select {
		// TODO 3: ctx context.Context がキャンセルされた場合には ln を close し、速やかに関数を終了する
		case <-ctx.Done():
			close(ln)
			return
		// TODO 1: 引数 r io.Reader とｓて標準出力が渡されてくるので、入力を待ち受ける
		default:
			// TODO 2: 入力があった場合は (scanner.Scan == true)
			if scanner.Scan() {
				ln <- scanner.Bytes()
				// TODO 2: ... 待受状態に戻る (returnしない)
			}
			// TODO 4: エラーが発生した際には引数 errc chan error へエラーを送信する
			err := scanner.Err()
			if err != nil {
				errc <- err
			}
		}
	}

}
