package client

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const testURL = "https://example.com?hoge=fuga&foo=var"

func TestNewHttpClient(t *testing.T) {
	type args struct {
		rawurl        string
		method        string
		data          string
		customHeaders []string
	}
	rawURL := "https://hoge.example.com"
	url, err := url.ParseRequestURI(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	data := `{"hoge":"fuga"}`
	tests := []struct {
		name    string
		args    args
		want    *HttpClient
		wantErr bool
	}{
		{
			name: "get request",
			args: args{
				rawurl: rawURL,
				method: "GET",
				data:   data,
				customHeaders: []string{
					"Connection: keep-alive",
					"Content-Type: application/json",
				},
			},
			want: &HttpClient{
				url:         url,
				method:      "GET",
				requestBody: nil,
				requestHeader: map[string]string{
					"Connection": "keep-alive",
				},
			},
			wantErr: false,
		},
		{
			name: "post request",
			args: args{
				rawurl: rawURL,
				method: "POST",
				data:   data,
				customHeaders: []string{
					"Connection: keep-alive",
				},
			},
			want: &HttpClient{
				url:         url,
				method:      "POST",
				requestBody: &data,
				requestHeader: map[string]string{
					"Connection":   "keep-alive",
					"Content-Type": "application/json",
				},
			},
			wantErr: false,
		},
		{
			name: "put request",
			args: args{
				rawurl: rawURL,
				method: "PUT",
				data:   data,
				customHeaders: []string{
					"Connection: keep-alive",
				},
			},
			want: &HttpClient{
				url:         url,
				method:      "PUT",
				requestBody: &data,
				requestHeader: map[string]string{
					"Connection":   "keep-alive",
					"Content-Type": "application/json",
				},
			},
			wantErr: false,
		},
		{
			name: "patch request",
			args: args{
				rawurl: rawURL,
				method: "PATCH",
				data:   data,
				customHeaders: []string{
					"Connection: keep-alive",
				},
			},
			want: &HttpClient{
				url:         url,
				method:      "PATCH",
				requestBody: &data,
				requestHeader: map[string]string{
					"Connection":   "keep-alive",
					"Content-Type": "application/json",
				},
			},
			wantErr: false,
		},
		{
			name: "delete request",
			args: args{
				rawurl: rawURL,
				method: "DELETE",
				data:   data,
				customHeaders: []string{
					"Connection: keep-alive",
					"Content-Type: application/json",
				},
			},
			want: &HttpClient{
				url:         url,
				method:      "DELETE",
				requestBody: nil,
				requestHeader: map[string]string{
					"Connection": "keep-alive",
				},
			},
			wantErr: false,
		},
		{
			name: "empty data json",
			args: args{
				rawurl: rawURL,
				method: "POST",
				data:   "",
				customHeaders: []string{
					"Connection: keep-alive",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHttpClient(tt.args.rawurl, tt.args.method, tt.args.data, tt.args.customHeaders)
			if (err != nil) != tt.wantErr {
				t.Errorf("HttpClientBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HttpClientBuilder.Build() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
[URL] https://example.com?hoge=fuga&foo=var
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

func TestHttpClient_SendRequest(t *testing.T) {
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
		wantReq  bool
		wantRes  bool
		wantErr  bool
	}{
		{
			name:     "success",
			mockFunc: success,
			wantErr:  false,
		},
		{
			name:     "error",
			mockFunc: errf,
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
				method:        http.MethodGet,
				requestBody:   nil,
				requestHeader: map[string]string{"Connection": "keep-alive"},
			}
			got, got1, err := c.SendRequest()
			defer func() {
				if got1 != nil {
					got1.Body.Close()
				}
			}()

			if (err != nil) != tt.wantErr {
				t.Errorf("HttpClient.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Nil(t, got1)
			} else {
				assert.NotNil(t, got)
				assert.Equal(t, url.String(), got.URL.String())
				assert.Equal(t, http.MethodGet, got.Method)
				assert.Equal(t, "keep-alive", got.Header.Values("Connection")[0])
				assert.Equal(t, 200, got1.StatusCode, 200)
				assert.Equal(t, "application/json", got1.Header.Values("Content-Type")[0])
				b, _ := io.ReadAll(got1.Body)
				assert.Equal(t, "\"{\\\"status\\\":\\\"ok\\\"}\"", string(b))
			}
		})
	}
}

func TestCreateRequestText(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		headers map[string]string
		want    string
	}{
		{
			name:   "with_header",
			method: http.MethodPost,
			headers: map[string]string{
				"Content-Type": "application/json",
				"Connection":   "keep-alive",
			},
			want: `
===Request===
[URL] https://example.com?hoge=fuga&foo=var
[Method] POST
[Headers]
  Connection: keep-alive
  Content-Type: application/json
`,
		},
		{
			name:    "without_header",
			method:  http.MethodGet,
			headers: nil,
			want: `
===Request===
[URL] https://example.com?hoge=fuga&foo=var
[Method] GET
[Headers]
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, testURL, nil)
			if err != nil {
				t.Fatal(err)
			}
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			assert.Equal(t, tt.want, CreateRequestText(req))
		})
	}
}

func TestCreateResponseText(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		body    string
		want    string
	}{
		{
			name: "with_header_and_body",
			headers: map[string]string{
				"Content-Type": "application/json",
				"Server":       "Apache",
			},
			body: `"{\"status\":\"ok\"}"`,
			want: `
===Response===
[Status] 200
[Headers]
  Content-Type: application/json
  Server: Apache
[Body]
"{\"status\":\"ok\"}"
`,
		},
		{
			name:    "without_header_and_body",
			headers: nil,
			body:    "",
			want: `
===Response===
[Status] 200
[Headers]
[Body]

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(tt.body)),
				Header:     make(http.Header),
			}
			defer res.Body.Close()
			for k, v := range tt.headers {
				res.Header.Set(k, v)
			}
			assert.Equal(t, tt.want, CreateResponseText(res))
		})
	}
}
