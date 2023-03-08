/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateCtx(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, got1 := NewCtx()
		assert.Implements(t, (*context.Context)(nil), got)
		assert.NotNil(t, got)
		assert.IsType(t, (context.CancelFunc)(nil), got1)
		assert.NotNil(t, got1)
	})
}

func TestStartBackgrounds(t *testing.T) {
	u, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		lines string
	}{
		{
			name:  "3 lines from STDOUT",
			lines: "line1\nline2\nline3\n",
		},
		{
			name:  "no lines from STDOUT",
			lines: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("POST", u.String(), func(req *http.Request) (*http.Response, error) {
				b, err := io.ReadAll(req.Body)
				if err != nil {
					return nil, err
				}

				if assert.Equal(t, tt.lines, string(b)) {
					return httpmock.NewStringResponse(400, ""), nil
				}
				return httpmock.NewStringResponse(200, ""), nil
			})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			buf := bytes.NewBuffer([]byte{})
			StartBackgrounds(ctx, u, buf)

			buf.WriteString(tt.lines)
			time.Sleep((timeSpan + 1) * time.Second)

			var expectCount int
			if tt.lines == "" {
				expectCount = 0
			} else {
				expectCount = 1
			}
			assert.Equal(t, expectCount, httpmock.GetTotalCallCount())
		})
	}
}
