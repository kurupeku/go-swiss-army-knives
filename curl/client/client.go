package client

import (
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
	// TODO: 1 週目：HTTP 通信用クライアントを構築
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	var getclient HttpClient
	getclient.url = u
	getclient.method = method
	getclient.requestHeader = make(map[string]string)

	for _, v := range customHeaders {
		s := strings.Split(v, ":")
		getclient.requestHeader[s[0]] = strings.TrimSpace(s[1])
	}

	switch getclient.method {
	case "GET", "DELETE":
		getclient.requestBody = nil
		_, ok := getclient.requestHeader["Content-Type"]
		if ok {
			delete(getclient.requestHeader, "Content-Type")
		}

	case "POST", "PUT", "PATCH":
		getclient.requestHeader["Content-Type"] = "application/json"

		if data == "" {
			return nil, errors.New("error!")
		}
		getclient.requestBody = &data
	}

	return &getclient, nil
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

	var r io.Reader
	if c.requestBody != nil {
		r = strings.NewReader(*c.requestBody)
	}

	req, err := http.NewRequest(c.method, c.url.String(), r)
	if err != nil {
		return nil, nil, err
	}

	for s, v := range c.requestHeader {
		req.Header.Set(s, v)
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
	msg := "\n"
	msg = msg + "===Request===\n"
	msg = msg + "[URL] " + req.URL.String() + "\n"
	msg = msg + "[Method] " + req.Method + "\n"
	msg = msg + "[Headers]" + "\n"

	for _, s := range sortedKeys(req.Header) {
		msg = msg + "  " + s + ": " + strings.Join(req.Header[s], "; ") + "\n"
	}
	fmt.Println(msg)
	return msg
}

// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func CreateResponseText(res *http.Response) string {
	// TODO: 3 週目：HTTP 通信結果のテキストを構築
	msg := ""
	for _, s := range sortedKeys(res.Header) {
		msg = msg + "  " + s + ": " + strings.Join(res.Header[s], "; ") + "\n"
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "err"
	}
	m := fmt.Sprintf("\n===Response===\n[Status] %d\n[Headers]\n%s[Body]\n%s\n", res.StatusCode, msg, b)
	//	res.Body

	return m
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
