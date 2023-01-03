package client

import (
	"reflect"
	"testing"
)

func TestHttpClientBuilder_Validate(t *testing.T) {
	type fields struct {
		rawurl        string
		method        string
		data          string
		customHeaders []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &HttpClientBuilder{
				rawurl:        tt.fields.rawurl,
				method:        tt.fields.method,
				data:          tt.fields.data,
				customHeaders: tt.fields.customHeaders,
			}
			if err := b.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("HttpClientBuilder.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHttpClientBuilder_Build(t *testing.T) {
	type fields struct {
		rawurl        string
		method        string
		data          string
		customHeaders []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *HttpClient
		wantErr bool
	}{
		// TODO: Add test cases.
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
