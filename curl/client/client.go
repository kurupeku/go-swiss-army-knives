package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
)

type HttpClient struct {
	url           *url.URL
	method        string
	requestBody   *string
	requestHeader map[string]string
}

const (
	GET    = "GET"
	DELETE = "DELETE"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
)

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
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	var hc HttpClient
	if method == GET || method == DELETE {
		hc.requestHeader = map[string]string{"Connection": "keep-alive"}
	}
	if method == POST || method == PUT || method == PATCH {
		if data == "" {
			return nil, errors.New("データが空です。")
		}
		hc.requestBody = &data
		hc.requestHeader = map[string]string{
			"Connection":   "keep-alive",
			"Content-Type": "application/json",
		}
	}

	return &HttpClient{
		url:           u,
		method:        method,
		requestBody:   hc.requestBody,
		requestHeader: hc.requestHeader,
	}, nil
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
	body, err := json.Marshal(c.requestBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest(c.method, c.url.String(), bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}

	c.setRequestHeader(req)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	return req, res, nil
}

func (c *HttpClient) setRequestHeader(req *http.Request) {
	for k, v := range c.requestHeader {
		req.Header.Set(k, v)
	}
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
