package input

import (
	"bufio"
	"context"
	"io"
)

// r から受け取ったデータを 1行ずつ ln へ送信する関数
func Monitor(ctx context.Context, ln chan []byte, errc chan error, r io.Reader) {
	// TODO: Implement here
	defer close(ln)

	scanner := bufio.NewScanner(r)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					errc <- err
				}
				return
			}

			ln <- scanner.Bytes()
		}
	}
}
