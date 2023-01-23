package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// structは funcの外に書ける。つまりfunc間の値受け渡しに使える
type HttpClient struct {
	url           *url.URL
	method        string
	requestBody   *string
	requestHeader map[string]string
}

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

	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Can't parse url: %s", rawurl)
	}

	rHeader, err := setHeader(method, customHeaders)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Can't parse url: %s", rawurl)
	}

	var rBody *string
	// method 判別
	switch method {
	case "GET", "DELETE":
		fmt.Print("method = GET,DELETE \n")
		rBody = nil
	case "POST", "PUT", "PATCH":
		fmt.Print("method = POST,PUT,PATCH \n")
		if data == "" {
			return nil, fmt.Errorf("[ERROR] rBody is empty with %s", method)
		} else {
			rBody = &data
		}
	default:
		return nil, fmt.Errorf("[ERROR] Unauthorized method: %s", method)
	}

	client := HttpClient{
		url:           u,
		method:        method,
		requestBody:   rBody,
		requestHeader: rHeader,
	}

	fmt.Printf("client.url = %v \n", client.url)
	fmt.Printf("client.method = %v \n", client.method)
	fmt.Printf("client.requestBody = %v \n", client.requestBody)
	fmt.Printf("client.requestHeader = %v \n", client.requestHeader)

	// HttpClientのポインターで戻り値を返す ,
	// errorはnil(正常終了)で戻り値を返す
	return &client, nil
}

func setHeader(method string, customHeaders []string) (map[string]string, error) {
	m := make(map[string]string)
	for _, l := range customHeaders {
		p := strings.Split(l, ":")
		// mapにkey , valueを追加。valueはSpaceを除去
		m[p[0]] = strings.TrimSpace(p[1])
	}

	switch method {
	case "GET", "DELETE":
		// "Content-Type" があったら消す
		_, ok := m["Content-Type"]
		if ok {
			delete(m, "Content-Type")
		}
	case "POST", "PUT", "PATCH":
		// "Content-Type": "application/json" 必須
		m["Content-Type"] = "application/json"
	default:
		return nil, fmt.Errorf("[ERROR] Unauthorized method: %s", method)
	}

	return m, nil
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

	// https://pkg.go.dev/net/http#NewRequest
	// https://pkg.go.dev/net/http#Request
	// https://golangstart.com/go_post/
	//   # bodyをreqが扱える形に変換
	data, _ := json.Marshal(c.requestBody)
	//   # http.NewRequest で struct http.Requestの初期化(値の代入)
	req, err := http.NewRequest(c.method, c.url.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, nil, err
	}
	//   # headerを順繰りに追加
	for k, v := range c.requestHeader {
		req.Header.Add(k, v)
	}

	// https://pkg.go.dev/net/http#Client.Do
	// https://golangstart.com/go_http_get/
	// # respは (*Client)Doで受ける。ので、*Clientの初期化(生成)が必要らしい (よくわからない。。)
	client := new(http.Client)
	// #
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("Request: \n  %T\n  %v\n\n", req, req)
	fmt.Printf("Response:\n  %T\n  %v\n\n", resp, resp)
	return req, resp, nil
}

// https://pkg.go.dev/net/http#Request

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
