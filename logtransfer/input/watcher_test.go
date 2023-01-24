package input

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitor(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name          string
		args          args
		wantSendCount int
	}{
		{
			name: "3 lines outputted",
			args: args{
				r: bytes.NewBufferString(
					`line1
					line2
					line3`,
				),
			},
			wantSendCount: 3,
		},
		{
			name: "5 lines outputted",
			args: args{
				r: bytes.NewBufferString(
					`line1
					line2
					line3
					line4
					line5`,
				),
			},
			wantSendCount: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			ln, errc := make(chan []byte, 1), make(chan error, 1)
			go Monitor(ctx, ln, errc, tt.args.r)
			for i := 0; i < tt.wantSendCount; i++ {
				b := <-ln
				assert.Greater(t, len(b), 0)
			}
			assert.Len(t, ln, 0)
		})
	}
}
