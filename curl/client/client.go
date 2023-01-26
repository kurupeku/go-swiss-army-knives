package client

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
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

	if method == "POST" || method == "PUT" || method == "PATCH" {
		if data == "" {
			return nil, errors.New("no Data")
		}
		http_client.requestBody = &data
	}

	http_client.method = method

	u, _ := url.Parse(rawurl)
	http_client.url = u

	if method == "GET" || method == "DELETE" {
		http_client.requestBody = nil
	}

	http_client.requestHeader = map[string]string{}
	for _, val := range customHeaders {
		v := strings.Split(val, ":")
		if v[0] == "Content-Type" {
			continue
		}
		http_client.requestHeader[v[0]] = strings.TrimSpace(v[1])
	}
	if method == "POST" || method == "PUT" || method == "PATCH" {
		http_client.requestHeader["Content-Type"] = "application/json"
	}
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

	//http.Requestを生成。bodyは無い場合があるのでチェック
	var b io.Reader //初期値nil
	if c.requestBody != nil {
		b = strings.NewReader(*c.requestBody)
	}

	req, _ := http.NewRequest(c.method, c.url.String(), b)

	//リクエストヘダーは複数あるはずなので繰り返しセット
	for k, v := range c.requestHeader {
		req.Header.Set(k, v)
	}

	cli := new(http.Client) //送信する発射台 newを使わないこともできる new初期化した上でポインターになる関数
	res, err := cli.Do(req)

	//エラーならerr以外何も返さない
	if err != nil {
		return nil, nil, err
	}

	return req, res, nil
}

// TODO:リクエストURL,HTTPメソッド,リクエストヘッダを所定のフォーマットで返却
func CreateRequestText(req *http.Request) string {
	// TODO: 3 週目：HTTP 通信結果のテキストを構築
	rs := "\n===Request===\n"
	rs += "[URL] " + req.URL.String() + "\n"
	rs += "[Method] " + req.Method + "\n"
	rs += "[Headers]"
	keys := sortedKeys(req.Header)
	var sameKey string
	for _, key := range keys {
		if sameKey == key {
			rs += "; " + req.Header.Get(key)
		} else {
			rs += "\n  " + key + ": " + req.Header.Get(key)
		}
	}
	rs += "\n"
	return rs
}

// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func CreateResponseText(res *http.Response) string {
	// TODO: 3 週目：HTTP 通信結果のテキストを構築
	rs := "\n===Response===\n"
	rs += "[Status] " + strconv.Itoa(res.StatusCode) + "\n"
	rs += "[Headers]"
	keys := sortedKeys(res.Header)
	var sameKey string
	for _, key := range keys {
		if sameKey == key {
			rs += "; " + res.Header.Get(key)
		} else {
			rs += "\n  " + key + ": " + res.Header.Get(key)
		}
	}

	rs += "\n[Body]\n"
	s := bufio.NewScanner(res.Body)
	for s.Scan() {
		rs += s.Text()
	}
	rs += "\n"
	return rs
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
