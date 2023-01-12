package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

type HttpClient struct {
	url           *url.URL
	method        string
	requestBody   *string
	requestHeader map[string]string
}

// TODO:URLはnet/urlパッケージの*url.URLで構築する
//
// TODO:customHeadersをリクエストヘッダとして設定
//
// TODO:HTTPメソッドがGET,DELETEの場合
// - リクエストボディは設定しない
// - リクエストヘッダにContent-Typeが含まれている場合は削除
//
// TODO:HTTPメソッドがPOST,PUT,DELETEの場合
// - リクエストヘッダのContent-Typeは"application/json"にする
// - dataの値をそのままレスポンスボディに設定
// - その際、dataが空であればエラー
func NewHttpClient(
	rawurl string,
	method string,
	data string,
	customHeaders []string,
) (*HttpClient, error) {
	var requestBody *string
	requestHeader := map[string]string{}
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}

	for _, v := range customHeaders {
		kv := strings.Split(v, ":")
		requestHeader[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}

	switch method {
	case http.MethodGet, http.MethodDelete:
		delete(requestHeader, "Content-Type")
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		if len(data) == 0 {
			return nil, errors.New("requires data json")
		}
		requestHeader["Content-Type"] = "application/json"
		requestBody = &data
	}

	return &HttpClient{u, method, requestBody, requestHeader}, nil
}

func (c *HttpClient) Execute() (string, string, error) {
	req, res, err := c.SendRequest()
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	return CreateRequestText(req), CreateResponseText(res), nil
}

// TODO:URL, HTTPメソッド, リクエストヘッダ, リクエストボディが適切に設定された*http.Requestを生成
// TODO:HTTPリクエストを実行後の*http.Request, *http.Responseを返却
func (c *HttpClient) SendRequest() (*http.Request, *http.Response, error) {
	var body io.Reader
	if c.requestBody != nil {
		body = bytes.NewBufferString(*c.requestBody)
	}
	req, err := http.NewRequest(c.method, c.url.String(), body)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range c.requestHeader {
		req.Header.Set(k, v)
	}

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	return req, res, nil
}

// TODO:リクエストURL,HTTPメソッド,リクエストヘッダを所定のフォーマットで返却
func CreateRequestText(req *http.Request) string {
	var sb strings.Builder
	sb.WriteString("\n===Request===\n")
	sb.WriteString(fmt.Sprintf("[URL] %s\n", req.URL.String()))
	sb.WriteString(fmt.Sprintf("[Method] %s\n", req.Method))
	sb.WriteString("[Headers]\n")
	for _, k := range sortedKeys(req.Header) {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(req.Header[k], "; ")))
	}
	return sb.String()
}

// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func CreateResponseText(res *http.Response) string {
	var sb strings.Builder
	sb.WriteString("\n===Response===\n")
	sb.WriteString(fmt.Sprintf("[Status] %d\n", res.StatusCode))
	sb.WriteString("[Headers]\n")
	for _, k := range sortedKeys(res.Header) {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(res.Header[k], "; ")))
	}
	sb.WriteString("[Body]\n")
	b, _ := io.ReadAll(res.Body)
	sb.WriteString(fmt.Sprintf("%s\n", string(b)))
	return sb.String()
}

// http.Request.Header と http.Response.Header を渡すと昇順にソートされた Key を返す関数
func sortedKeys(m map[string][]string) []string {
	s := make([]string, 0, len(m))
	for k := range m {
		s = append(s, k)
	}

	sort.Strings(s)
	return s
}
