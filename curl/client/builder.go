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
	rawurl              string
	method              string
	data                []string
	withQueryParamsFlag bool
	customHeaders       []string
}

func NewHttpClientBuilder(
	rawurl string,
	method string,
	data []string,
	withQueryParamsFlag bool,
	customHeaders []string,
) *HttpClientBuilder {
	return &HttpClientBuilder{rawurl, method, data, withQueryParamsFlag, customHeaders}
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
// - b.withQueryParamsFlagがtrueの場合にHTTPメソッドはGETかDELETEになっているか
func (b *HttpClientBuilder) validateMethod() error {
	switch b.method {
	case http.MethodGet, http.MethodDelete:
		return nil
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if b.withQueryParamsFlag {
			return fmt.Errorf("HTTP method '%s' is not allowed to use with quary params flag(--get)", b.method)
		}
		return nil
	}

	return fmt.Errorf("HTTP method '%s' is not supported", b.method)
}

// TODO:チェックをして問題があればerrorを返却
// b.dataのすべての要素が以下の条件を満たす
// - 空文字ではない
// - '{'で始まる場合は正しいJSON形式になっている
// - '{'で始まらない場合は'='区切りの文字列、または'='区切りの文字列が'&'区切りで連結されていること
//   - ok:a=b
//   - ok:a=b&c=d
//   - ng:ab('='区切りになっていない)
//   - ng:ab&c=d('&'で連結されたすべての文字列が'='区切りになっていない)
func (b *HttpClientBuilder) validateData() error {
	for _, v := range b.data {
		s := strings.TrimSpace(v)
		if len(s) == 0 {
			return errors.New("empty data")
		}
		if strings.HasPrefix(s, "{") {
			m := make(map[string]interface{})
			if err := json.Unmarshal([]byte(s), &m); err != nil {
				return fmt.Errorf("json parse error: %s", err.Error())
			}
		} else {
			kvs := strings.Split(s, "&")
			for _, kv := range kvs {
				if len(strings.Split(kv, "=")) != 2 {
					return fmt.Errorf("invalid format data: %s", kv)
				}
			}
		}
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
			return errors.New("empty data")
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
// - b.withQueryParamsFlagがtrueの場合は、b.dataをURLのクエリパラメータに設定する('&'で連結)
// - その際、b.dataはすべてパーセントエンコーディング(クエリパラメータ用のパーセントエンコーディング)
// - リクエストヘッダにContent-Typeが含まれている場合は削除
//
// TODO:HTTPメソッドがPOST,PUT,DELETEの場合
// - b.customHeadersにContent-Typeが設定されていない場合はContent-Typeを"application/x-www-form-urlencoded"としてリクエストヘッダを設定
// - Content-Typeは"application/x-www-form-urlencoded"と"application/json"以外許容しない
// - "application/x-www-form-urlencoded"の場合はすべてのb.dataを'&'で連結してレスポンスボディに設定
// - その際、b.dataはすべてパーセントエンコーディング(クエリパラメータ用のパーセントエンコーディング)
// - "application/json"の場合はb.dataの先頭をそのままレスポンスボディに設定
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
		if b.withQueryParamsFlag {
			query := u.Query()
			err := b.setValues(&query)
			if err != nil {
				return nil, err
			}
			u.RawQuery = query.Encode()
		}
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		contentType, exists := requestHeader["Content-Type"]
		if !exists {
			contentType = "application/x-www-form-urlencoded"
			requestHeader["Content-Type"] = contentType
		}

		switch contentType {
		case "application/x-www-form-urlencoded":
			values := url.Values{}
			err := b.setValues(&values)
			if err != nil {
				return nil, err
			}
			bodyStr := values.Encode()
			requestBody = &bodyStr
		case "application/json":
			if len(b.data) == 0 {
				return nil, errors.New("requires data")
			}
			bodyStr := b.data[0]
			requestBody = &bodyStr
		default:
			return nil, errors.New("invalid content-type")
		}
	}

	return &HttpClient{u, b.method, requestBody, requestHeader}, nil
}

func (b *HttpClientBuilder) setValues(values *url.Values) error {
	for _, v := range b.data {
		qs := strings.Split(v, "&")
		for _, q := range qs {
			kv := strings.Split(q, "=")
			if len(kv) != 2 {
				return errors.New("invalid format data")
			}
			values.Set(kv[0], kv[1])
		}
	}

	return nil
}
