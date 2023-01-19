package client

import (
	"errors"
	"fmt"
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
	// TODO: 1 週目：HTTP 通信用クライアントを構築

	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Can't parse url: %s", rawurl)
	}

	switch method {
	case "GET", "DELETE", "POST", "PUT", "PATCH":
	default:
		return nil, fmt.Errorf("[ERROR] Unauthorized method: %s", method)
	}

	rBody, err := setBody(method, data)
	if err != nil {
		return nil, err
	}
	rHeaders := setHeaders(method, customHeaders)

	client := HttpClient{
		url:           u,
		method:        method,
		requestBody:   rBody,
		requestHeader: rHeaders,
	}

	return &client, nil
}

func setBody(method, data string) (*string, error) {
	var rb *string
	switch method {
	case "GET", "DELETE":
		rb = nil
	case "POST", "PUT", "PATCH":
		if data == "" {
			return nil, errors.New("[ERROR] Empty data")
		}
		rb = &data
	}
	return rb, nil
}

func setHeaders(method string, ch []string) map[string]string {
	rh := make(map[string]string)
	switch method {
	case "GET", "DELETE":
		for _, h := range ch {
			if !strings.Contains(h, "Content-Type") {
				kv := strings.Split(h, ":")
				rh[kv[0]] = strings.TrimSpace(kv[1])
			}
		}
	case "POST", "PUT", "PATCH":
		for _, h := range ch {
			kv := strings.Split(h, ":")
			rh[kv[0]] = strings.TrimSpace(kv[1])
		}
		rh["Content-Type"] = "application/json"
	}
	return rh
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
// TODO:ただ単にオブジェクトを作るだけでなく、このメソッド内でリクエストの実行も完了させる
func (c *HttpClient) SendRequest() (*http.Request, *http.Response, error) {
	// TODO: 2 週目：HTTP 通信を実行
	req, err := http.NewRequest(c.method, c.url.String(), nil)
	if err != nil {
		return nil, nil, err
	}
	for k, v := range c.requestHeader {
		req.Header.Add(k, v)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	return req, resp, nil
}

// TODO:リクエストURL,HTTPメソッド,リクエストヘッダを所定のフォーマットで返却
func CreateRequestText(req *http.Request) string {
	// TODO: 3 週目：HTTP 通信結果のテキストを構築
	return ""
}

// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func CreateResponseText(res *http.Response) string {
	// TODO: 3 週目：HTTP 通信結果のテキストを構築
	return ""
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
