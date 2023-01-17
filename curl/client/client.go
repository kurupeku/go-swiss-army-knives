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

	http_client := HttpClient{}
	http_client.requestHeader = map[string]string{}
	http_client.method = method

	// urlフィールド は net/url パッケージの*url.URL で構築する
	u, _ := url.Parse(rawurl)
	http_client.url = u

	// リクエストボディ(requestBodyフィールド)はnil
	if method == "GET" || method == "DELETE" {
		http_client.requestBody = nil
	}
	// data の値をそのままレスポンスボディ(requestBodyフィールド)に設定
	// その際、data が空であればエラー
	if method == "POST" || method == "PUT" || method == "PATCH" {
		if data == "" {
			return nil, errors.New("no Data")
		}
		http_client.requestBody = &data
	}

	// customHeaders引数の要素を:で区切って、requestHeaderフィールドのキーと値に設定
	// HTTP メソッドが GET,DELETE の場合
	// リクエストヘッダに Content-Type が含まれている場合は削除
	// HTTP メソッドが POST,PUT,DELETE の場合
	// リクエストヘッダの Content-Type は"application/json"にする
	for i := range customHeaders {
		v := strings.Split(customHeaders[i], ":")
		if v[0] == "Content-Type" {
			continue
		}
		http_client.requestHeader[v[0]] = strings.TrimSpace(v[1])
	}
	if method == "POST" || method == "PUT" || method == "PATCH" {
		http_client.requestHeader["Content-Type"] = "application/json"
	}
	fmt.Println(http_client)
	return &http_client, nil
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
	return nil, nil, nil
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
