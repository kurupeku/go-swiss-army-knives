package client

import (
	"net/url"
	"testing"
)

func TestHttpClient_Execute(t *testing.T) {
	type fields struct {
		url           *url.URL
		method        string
		requestBody   *string
		requestHeader map[string]string
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
			c := &HttpClient{
				url:           tt.fields.url,
				method:        tt.fields.method,
				requestBody:   tt.fields.requestBody,
				requestHeader: tt.fields.requestHeader,
			}
			if err := c.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("HttpClient.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
