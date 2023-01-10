package client

import (
	"net/url"
	"reflect"
	"testing"
)

func TestHttpClientBuilder_Build(t *testing.T) {
	type fields struct {
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
		fields  fields
		want    *HttpClient
		wantErr bool
	}{
		{
			name: "get request",
			fields: fields{
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
			fields: fields{
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
			fields: fields{
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
			fields: fields{
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
			fields: fields{
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
			fields: fields{
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
			b := &HttpClientBuilder{
				rawurl:        tt.fields.rawurl,
				method:        tt.fields.method,
				data:          tt.fields.data,
				customHeaders: tt.fields.customHeaders,
			}
			got, err := b.Build()
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
