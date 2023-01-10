package client

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
)

const testURL = "https://example.com"

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
