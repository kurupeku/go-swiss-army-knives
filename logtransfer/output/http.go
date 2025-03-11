package output

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	contentType = "text/plain"
)

// Forward は out から受け取ったデータを指定の url へ HTTP POST する関数
func Forward(ctx context.Context, out chan []byte, errc chan error, url string) {
	// TODO: Implement here

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-out:
			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				url,
				bytes.NewReader(data),
			)
			if err != nil {
				errc <- err
				continue
			}
			req.Header.Set("Content-Type", contentType)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				errc <- err
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					errc <- fmt.Errorf("HTTP request failed with status: %s", resp.Status)
					continue
				}
				errc <- fmt.Errorf("HTTP request failed with status: %s, body: %s", resp.Status, string(body))
			}
		}
	}
}
