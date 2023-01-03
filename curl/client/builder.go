package client

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (b *HttpClientBuilder) Validate() error {
	//rawurlのフォーマットをチェック
	if err := b.validateRawURL(); err != nil {
		return err
	}

	//methodの整合性をチェック
	if err := b.validateMethod(); err != nil {
		return err
	}

	//dataのフォーマットをチェック
	if err := b.validateData(); err != nil {
		return err
	}

	//customHeadersのフォーマットをチェック
	if err := b.validateHeader(); err != nil {
		return err
	}

	return nil
}

// TODO: チェックをして問題があればerrorを返却
// - b.urlのフォーマットがURLとして正しいかをチェック
// - プロトコルがhttp, httpsになっているかを確認する
func (b *HttpClientBuilder) validateRawURL() error {
	url, err := url.ParseRequestURI(b.rawurl)
	if err != nil {
		return err
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return fmt.Errorf("url schema '%s' is not supported", url.Scheme)
	}

	return nil
}

// TODO: チェックをして問題があればerrorを返却
// - b.methodが許容されているHTTPメソッド(GET, POST, PUT, DELETE, PATCH)になっているか
func (b *HttpClientBuilder) validateMethod() error {
	switch b.method {
	case http.MethodGet:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodDelete:
	case http.MethodPatch:
		return nil
	}
	return fmt.Errorf("HTTP method '%s' is not supported", b.method)
}

// TODO:チェックをして問題があればerrorを返却
// b.dataが以下のいずれかの条件を満たす
// - 空文字
// - 正しいJSON形式の文字列
func (b *HttpClientBuilder) validateData() error {
	s := strings.TrimSpace(b.data)
	if s == "" {
		return nil
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return fmt.Errorf("json parse error: %s", err.Error())
	}
	return nil
}

// TODO:チェックをして問題があればerrorを返却
// b.customHeadersのすべての要素が以下の条件を満たす
// - 空文字ではない
// - ':'が1つだけ含まれており、':'の前後が空ではない
func (b *HttpClientBuilder) validateHeader() error {
	for _, v := range b.customHeaders {
		s := strings.TrimSpace(v)
		if len(s) == 0 {
			return errors.New("headers include empty string")
		}
		kv := strings.Split(s, ":")
		if len(kv) != 2 {
			return fmt.Errorf("invalid format header: %s", kv)
		}
	}
	return nil
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
