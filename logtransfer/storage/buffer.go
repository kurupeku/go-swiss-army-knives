package storage

import (
	"bytes"
	"context"
	"io/ioutil"
	"time"
)

var (
	buf = bytes.NewBuffer([]byte{})
)

func Listen(ctx context.Context, ln chan []byte, errc chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case b, ok := <-ln:
			if !ok {
				return
			}
			_, err := buf.Write([]byte(string(b) + "\n"))
			if err != nil {
				errc <- err
			}
		}
	}
}

func Load(ctx context.Context, out chan []byte, errc chan error, span time.Duration) {
	tick := time.NewTicker(span)
	for {
		select {
		case <-ctx.Done():
			close(out)
			return
		case <-tick.C:
			b, err := ioutil.ReadAll(buf)
			if err != nil {
				errc <- err
				continue
			}
			if len(b) == 0 {
				continue
			}
			out <- b
		}
	}
}
