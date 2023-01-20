package client

import (
	"bytes"
	"encoding/json"
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
	c := HttpClient{}
	c.url, _ = url.Parse(rawurl)
	c.method = method
	c.requestHeader = make(map[string]string)
	for _, v := range customHeaders {
		e := strings.Split(v, ":")
		c.requestHeader[e[0]] = strings.TrimSpace(e[1])
	}
	switch method {
	case "GET", "DELETE":
		c.requestBody = nil
		if c.requestHeader["Content-Type"] != "" {
			delete(c.requestHeader, "Content-Type")
		}
	case "POST", "PUT", "PATCH":
		if data == "" {
			return nil, errors.New("empty data")
		}
		c.requestBody = &data
		c.requestHeader["Content-Type"] = "application/json"
	default:
		return nil, errors.New("invalid method")
	}
	return &c, nil
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
	jsonString, _ := json.Marshal(c.requestBody)
	req, err := http.NewRequest(c.method, c.url.String(), bytes.NewBuffer(jsonString))
	if err != nil {
		return nil, nil, err
	}
	for k, _ := range c.requestHeader {
		req.Header.Set(k, c.requestHeader[k])
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	return req, res, nil
}

// TODO:リクエストURL,HTTPメソッド,リクエストヘッダを所定のフォーマットで返却
func CreateRequestText(req *http.Request) string {
	// TODO: 3 週目：HTTP 通信結果のテキストを構築
	reqstr := fmt.Sprintf("\n===Request===\n[URL] %s\n[Method] %s\n[Headers]\n", req.URL, req.Method)
	for _, k := range sortedKeys(req.Header) {
		reqstr = fmt.Sprintf(reqstr+"  %s: %s\n", k, req.Header[k][0])
	}
	return reqstr
}

// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func CreateResponseText(res *http.Response) string {
	// TODO: 3 週目：HTTP 通信結果のテキストを構築
	resstr := fmt.Sprintf("\n===Response===\n[Status] %d\n[Headers]", res.StatusCode)
	for _, k := range sortedKeys(res.Header) {
		resstr = fmt.Sprintf(resstr+"\n  %s: %s", k, strings.Join(res.Header[k], "; "))
	}
	resstr = fmt.Sprintf(resstr + "\n")
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	resstr = fmt.Sprintf(resstr+"[Body]\n%s\n", buf.String())
	return resstr
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
