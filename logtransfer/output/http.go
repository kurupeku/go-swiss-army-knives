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
