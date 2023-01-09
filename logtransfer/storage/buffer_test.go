package storage

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListen(t *testing.T) {
	tests := []struct {
		name    string
		lines   []string
		wantBuf *bytes.Buffer
	}{
		{
			name:  "3 lines sended",
			lines: []string{"line1", "line2", "line3"},
			wantBuf: bytes.NewBufferString(
				"line1\nline2\nline3\n",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ln, errc := make(chan []byte, 1), make(chan error, 1)
			go Listen(ctx, ln, errc)
			for _, line := range tt.lines {
				ln <- []byte(line)
			}
			time.Sleep(500 * time.Millisecond)
			assert.Equal(t, tt.wantBuf.Bytes(), buf.Bytes())
			buf.Reset()
		})
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		span time.Duration
	}
	tests := []struct {
		name   string
		args   args
		dumped string
	}{
		{
			name:   "transferred every 1 seconds",
			args:   args{span: 1 * time.Second},
			dumped: "dumped\nstring\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.args.span+500*time.Millisecond)
			defer cancel()
			out, errc := make(chan []byte, 1), make(chan error, 1)
			buf.WriteString(tt.dumped)
			go Load(ctx, out, errc, tt.args.span)
			sended := <-out
			assert.Equal(t, []byte(tt.dumped), sended)
			assert.Equal(t, buf.Len(), 0)
			buf.Reset()
		})
	}
}
