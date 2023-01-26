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
		{
			name:  "5 lines sended",
			lines: []string{"line1", "line2", "line3", "line4", "line5"},
			wantBuf: bytes.NewBufferString(
				"line1\nline2\nline3\nline4\nline5\n",
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
	span := 1 * time.Second
	tests := []struct {
		name   string
		dumped string
	}{
		{
			name:   "transferred every 1 seconds",
			dumped: "dumped\nstring\n",
		},
		{
			name:   "when buf is blank",
			dumped: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer buf.Reset()

			wait := 500 * time.Millisecond
			ctx, cancel := context.WithTimeout(context.Background(), span+wait)
			defer cancel()
			out, errc := make(chan []byte, 1), make(chan error, 1)
			go Load(ctx, out, errc, span)

			buf.WriteString(tt.dumped)
			if tt.dumped != "" {
				sended := <-out
				assert.Equal(t, []byte(tt.dumped), sended)
			} else {
				time.Sleep(wait)
				assert.Len(t, out, 0)
			}
			assert.Equal(t, buf.Len(), 0)
		})
	}
}
