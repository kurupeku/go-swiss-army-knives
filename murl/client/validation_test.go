package client

import (
	"testing"
)

func TestValidateFlags(t *testing.T) {
	type args struct {
		rawurl        string
		method        string
		data          string
		customHeaders []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "get request",
			args: args{
				rawurl:        "https://hoge.example.com",
				method:        "GET",
				customHeaders: []string{"Connection: keep-alive"},
			},
			wantErr: false,
		},
		{
			name: "post request",
			args: args{
				rawurl:        "http://hoge.example.com",
				method:        "POST",
				data:          `{"hoge":"fuga"}`,
				customHeaders: []string{"Connection: keep-alive"},
			},
			wantErr: false,
		},
		{
			name: "put request",
			args: args{
				rawurl:        "http://hoge.example.com",
				method:        "PUT",
				data:          `{"hoge":"fuga"}`,
				customHeaders: []string{"Connection: keep-alive"},
			},
			wantErr: false,
		},
		{
			name: "patch request",
			args: args{
				rawurl:        "http://hoge.example.com",
				method:        "PATCH",
				data:          `{"hoge":"fuga"}`,
				customHeaders: []string{"Connection: keep-alive"},
			},
			wantErr: false,
		},
		{
			name: "delete request",
			args: args{
				rawurl:        "https://hoge.example.com",
				method:        "DELETE",
				customHeaders: []string{"Connection: keep-alive"},
			},
			wantErr: false,
		},
		{
			name: "invalid url",
			args: args{
				rawurl:        "https//hoge.example.com",
				method:        "GET",
				customHeaders: []string{"Connection: keep-alive"},
			},
			wantErr: true,
		},
		{
			name: "invalid schema",
			args: args{
				rawurl:        "ftp://hoge.example.com",
				method:        "GET",
				customHeaders: []string{"Connection: keep-alive"},
			},
			wantErr: true,
		},
		{
			name: "unsupported method",
			args: args{
				rawurl:        "https://hoge.example.com",
				method:        "HEAD",
				customHeaders: []string{"Connection: keep-alive"},
			},
			wantErr: true,
		},
		{
			name: "invalid data",
			args: args{
				rawurl:        "http://hoge.example.com",
				method:        "POST",
				data:          `{"hoge:"fuga"}`,
				customHeaders: []string{"Connection: keep-alive"},
			},
			wantErr: true,
		},
		{
			name: "empty header text",
			args: args{
				rawurl:        "http://hoge.example.com",
				method:        "POST",
				data:          `{"hoge":"fuga"}`,
				customHeaders: []string{""},
			},
			wantErr: true,
		},
		{
			name: "invalid header text",
			args: args{
				rawurl:        "http://hoge.example.com",
				method:        "POST",
				data:          `{"hoge":"fuga"}`,
				customHeaders: []string{"hoge:huga:hige"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateFlags(tt.args.rawurl, tt.args.method, tt.args.data, tt.args.customHeaders); (err != nil) != tt.wantErr {
				t.Errorf("HttpClientBuilder.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
