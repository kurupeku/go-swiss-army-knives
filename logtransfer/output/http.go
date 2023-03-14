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

func Forward(ctx context.Context, out chan []byte, errc chan error, url string) {
	// TODO: 2 週目：内部バッファに保存された内容を一定時間ごとに読み込む処理と、読み取った文字列を Body とした HTTP#POST リクエストを投げる処理
	for {
		select {
		// TODO3: ctx context.Context がキャンセルされた場合には速やかに関数を終了する
		case <-ctx.Done():
			return
		// TODO1: 引数 `out chan []byte` で文字列を受信した際に
		case b, ok := <-out:
			if !ok {
				return
			}

			// TODO1: その内容 Body として引数 `url string` への HTTP#POST リクエストを行う
			body := bytes.NewBuffer(b)
			// TODO2: `Content-Type: plain/text` を Header に添えて送信を行う
			res, err := http.Post(url, contentType, body)
			// TODO4: エラーが発生した際には errc chan error へエラーを送信する
			if err != nil {
				errc <- err
				continue
			}

			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				b, err := io.ReadAll(res.Body)
				// TODO4: エラーが発生した際には errc chan error へエラーを送信する
				if err != nil {
					errc <- err
				} else {
					errc <- errors.New(string(b))
				}
			}
		}
	}
}
