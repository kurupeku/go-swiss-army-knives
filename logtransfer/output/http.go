package output

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

const (
	contentType = "text/plain"
)

// TODO: 引数 `out chan []byte` で文字列を受信した際に、その内容 Body として引数 `url string` への HTTP#POST リクエストを行う
// TODO: `Content-Type: plain/text` を Header に添えて送信を行う
// TODO: ctx context.Context がキャンセルされた場合には速やかに関数を終了する
// TODO: エラーが発生した際には errc chan error へエラーを送信する
func Forward(ctx context.Context, out chan []byte, errc chan error, url string) {
	// TODO: 2 週目：内部バッファに保存された内容を一定時間ごとに読み込む処理と、読み取った文字列を Body とした HTTP#POST リクエストを投げる処理
	for {
		select {
		case <-ctx.Done():
			close(errc)
			return
		case b, ok := <-out:
			if ok {
				res, err := http.Post(url, contentType, bytes.NewReader(b))
				if err != nil {
					errc <- err
					return
				}
				defer res.Body.Close()
				if res.StatusCode != http.StatusOK {
					b, err := io.ReadAll(res.Body)
					if err != nil {
						errc <- err
						return
					} else {
						st := string(b)
						errc <- errors.New(st)
					}

				}
			}
		}
	}
}
