package client

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type HttpClientBuilder struct {
	rawurl        string
	method        string
	data          string
	customHeaders []string
}

func NewHttpClientBuilder(
	rawurl string,
	method string,
	data string,
	customHeaders []string,
) *HttpClientBuilder {
	return &HttpClientBuilder{rawurl, method, data, customHeaders}
}

// TODO:URLはnet/urlパッケージの*url.URLで構築する
//
// TODO:b.customHeadersをリクエストヘッダとして設定
//
// TODO:HTTPメソッドがGET,DELETEの場合
// - リクエストボディは設定しない
// - リクエストヘッダにContent-Typeが含まれている場合は削除
//
// TODO:HTTPメソッドがPOST,PUT,DELETEの場合
// - リクエストヘッダのContent-Typeは"application/json"にする
// - b.dataの値をそのままレスポンスボディに設定
// - その際、b.dataが空であればエラー
func (b *HttpClientBuilder) Build() (*HttpClient, error) {
	var requestBody *string
	requestHeader := map[string]string{}
	u, err := url.ParseRequestURI(b.rawurl)
	if err != nil {
		return nil, err
	}

	for _, v := range b.customHeaders {
		kv := strings.Split(v, ":")
		requestHeader[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	switch b.method {
	case http.MethodGet, http.MethodDelete:
		delete(requestHeader, "Content-Type")
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if len(b.data) == 0 {
			return nil, errors.New("requires data json")
		}
		requestHeader["Content-Type"] = "application/json"
		requestBody = &b.data
	}

	return &HttpClient{u, b.method, requestBody, requestHeader}, nil
}
