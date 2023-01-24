package client

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	newHC := HttpClient{}
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, errors.New("URLが不正なフォーマットです。")
	}
	newHC.url = u
	newHC.method = method
	newHC.requestHeader = make(map[string]string)
	for _, header := range customHeaders {
		parts := strings.SplitN(header, ": ", 2)
		key := parts[0]
		value := parts[1]
		newHC.requestHeader[key] = value
	}

	switch method {
	case "GET", "DELETE":
		newHC.requestBody = nil
		delete(newHC.requestHeader, "Content-Type")

	case "POST", "PUT", "PATCH":
		if data == "" {
			return nil, errors.New("dataが空です。")
		}
		newHC.requestBody = &data
		newHC.requestHeader["Content-Type"] = "application/json"
	}

	return &newHC, nil
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
	var httpReq *http.Request
	var err error

	if c.requestBody != nil {
		httpReq, err = http.NewRequest(c.method, c.url.String(), strings.NewReader(*c.requestBody))
	} else {
		httpReq, err = http.NewRequest(c.method, c.url.String(), nil)
	}
	if err != nil {
		return nil, nil, err
	}

	for key, value := range c.requestHeader {
		httpReq.Header.Add(key, value)
	}

	client := &http.Client{}
	httpRes, err := client.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}

	return httpReq, httpRes, err

}

// TODO:リクエストURL,HTTPメソッド,リクエストヘッダを所定のフォーマットで返却
func CreateRequestText(req *http.Request) string {
	// TODO: 3 週目：HTTP 通信結果のテキストを構築
	reqText := "\n===Request===\n"
	reqText += fmt.Sprintf("[URL] %s\n", req.URL.String())
	reqText += fmt.Sprintf("[Method] %s\n", req.Method)
	reqText += "[Headers]\n"

	sortHeader := sortedKeys(req.Header)
	for _, key := range sortHeader {
		for _, value := range req.Header[key] {
			reqText += fmt.Sprintf("  %s: %s\n", key, value)
		}
	}

	return reqText
}

// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func CreateResponseText(res *http.Response) string {
	// TODO: 3 週目：HTTP 通信結果のテキストを構築
	resText := "\n===Response===\n"
	resText += fmt.Sprintf("[Status] %d\n", res.StatusCode)
	resText += "[Headers]\n"

	sortHeader := sortedKeys(res.Header)
	for _, key := range sortHeader {
		for _, value := range res.Header[key] {
			resText += fmt.Sprintf("  %s: %s\n", key, value)
		}
	}

	resText += "[Body]\n"
	resBody, _ := ioutil.ReadAll(res.Body)
	resText += string(resBody) + "\n"

	return resText
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
