package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	req, err := c.newRequest()
	if err != nil {
		return "", "", err
	}

	return sendRequest(req)
}

// TODO:URL, HTTPメソッド, リクエストヘッダ, リクエストボディが適切に設定された*http.Requestを返却
func (c *HttpClient) newRequest() (*http.Request, error) {
	var body io.Reader
	if c.requestBody != nil {
		body = bytes.NewBufferString(*c.requestBody)
	}
	req, err := http.NewRequest(c.method, c.url.String(), body)
	if err != nil {
		return nil, err
	}
	for k, v := range c.requestHeader {
		req.Header.Set(k, v)
	}

	return req, nil
}

// TODO:HTTPリクエストを実行してレスポンスを取得
// TODO:リクエストURL,HTTPメソッド,リクエストヘッダを所定のフォーマットで返却
// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func sendRequest(req *http.Request) (request string, response string, e error) {
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		e = err
		return
	}
	defer res.Body.Close()

	request = createRequestText(req)
	response = createResponseText(res)

	return
}

func createRequestText(req *http.Request) string {
	var sb strings.Builder
	sb.WriteString("\n===Request===\n")
	sb.WriteString(fmt.Sprintf("[URL] %s\n", req.URL.String()))
	sb.WriteString(fmt.Sprintf("[Method] %s\n", req.Method))
	sb.WriteString("[Headers]\n")
	for k, v := range req.Header {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, "; ")))
	}
	return sb.String()
}

func createResponseText(res *http.Response) string {
	var sb strings.Builder
	sb.WriteString("\n===Response===\n")
	sb.WriteString(fmt.Sprintf("[Status] %d\n", res.StatusCode))
	sb.WriteString("[Headers]\n")
	for k, v := range res.Header {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, "; ")))
	}
	sb.WriteString("[Body]\n")
	b, _ := io.ReadAll(res.Body)
	sb.WriteString(fmt.Sprintf("%s\n", string(b)))
	return sb.String()
}
