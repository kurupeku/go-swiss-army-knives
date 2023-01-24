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
	for {
		select {
		case <-ctx.Done():
			return
		case b, ok := <-out:
			if !ok {
				return
			}

			body := bytes.NewBuffer(b)
			res, err := http.Post(url, contentType, body)
			if err != nil {
				errc <- err
				continue
			}

			defer res.Body.Close()
			if res.StatusCode != http.StatusOK {
				b, err := io.ReadAll(res.Body)
				if err != nil {
					errc <- err
				} else {
					errc <- errors.New(string(b))
				}
			}
		}
	}
}
