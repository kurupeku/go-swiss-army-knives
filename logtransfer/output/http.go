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

			req := newRequest(url, b, errc)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				errc <- err
				continue
			}

			checkError(res, errc)
		}
	}
}

func newRequest(url string, body []byte, errc chan error) *http.Request {
	buf := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		errc <- err
		return nil
	}

	req.Header.Add("Content-Type", contentType)
	return req
}

func checkError(res *http.Response, errc chan error) {
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
