package input

import (
	"bufio"
	"context"
	"io"
)

func Monitor(ctx context.Context, ln chan []byte, errc chan error, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for {
		select {
		case <-ctx.Done():
			close(ln)
			return
		default:
			if scanner.Scan() {
				ln <- scanner.Bytes()
			}
			if err := scanner.Err(); err != nil {
				errc <- err
			}
		}
	}
}
