package output

import (
	"bytes"
	"context"
	"errors"
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
			if err := ctx.Err(); errors.Is(err, context.Canceled) {
				errc <- err
			} else if errors.Is(err, context.DeadlineExceeded) {
				errc <- err
			}
			return
		default:
			body := bytes.NewBuffer(<-out)
			req, err := http.NewRequest("POST", url, body)
			if err != nil {
				errc <- err
			}
			req.Header.Set("Content-Type", "plain/text")
			client := &http.Client{}
			_, err = client.Do(req)
			if err != nil {
				errc <- err
			}
		}
	}
}
