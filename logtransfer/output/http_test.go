package output

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestForward(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name         string
		args         args
		receivedBody string
	}{
		{
			name: "received buffered string",
			args: args{
				url: "https://example.com",
			},
			receivedBody: "body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("POST", tt.args.url, func(req *http.Request) (*http.Response, error) {
				b, err := io.ReadAll(req.Body)
				if err != nil {
					return nil, err
				}

				if assert.Equal(t, tt.receivedBody, string(b)) {
					return httpmock.NewStringResponse(400, ""), nil
				}

				return httpmock.NewStringResponse(200, ""), nil
			})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			out, errc := make(chan []byte, 1), make(chan error, 1)
			go Forward(ctx, out, errc, tt.args.url)
			out <- []byte(tt.receivedBody)

			time.Sleep(1 * time.Second)
			assert.Equal(t, 1, httpmock.GetTotalCallCount())
		})
	}
}
