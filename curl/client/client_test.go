package client

import (
	"errors"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
)

const testURL = "https://example.com"

func TestHttpClient_Execute(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	url, err := url.ParseRequestURI(testURL)
	if err != nil {
		t.Fatal(err)
	}
	success := func() {
		res, err := httpmock.NewJsonResponse(200, `{"status":"ok"}`)
		if err != nil {
			t.Fatal(err)
		}
		httpmock.RegisterResponder("GET", testURL, httpmock.ResponderFromResponse(res))
	}
	errf := func() {
		httpmock.RegisterResponder("GET", testURL, httpmock.NewErrorResponder(errors.New("error!!")))
	}

	tests := []struct {
		name     string
		mockFunc func()
		want     string
		want1    string
		wantErr  bool
	}{
		{
			name:     "success",
			mockFunc: success,
			want: `
===Request===
[URL] https://example.com
[Method] GET
[Headers]
  Connection: keep-alive
`,
			want1: `
===Response===
[Status] 200
[Headers]
  Content-Type: application/json
[Body]
"{\"status\":\"ok\"}"
`,
			wantErr: false,
		},
		{
			name:     "error",
			mockFunc: errf,
			want:     "",
			want1:    "",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc()
			}
			c := &HttpClient{
				url:           url,
				method:        "GET",
				requestBody:   nil,
				requestHeader: map[string]string{"Connection": "keep-alive"},
			}
			got, got1, err := c.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpClient.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HttpClient.Execute() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("HttpClient.Execute() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
