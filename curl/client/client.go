package client

import (
	"bytes"
	"errors"
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
	d := data

	// URL生成
	u, err := url.Parse(rawurl)

	// リクエストヘッダ生成
	mp := make(map[string]string)
	for _, v := range customHeaders {
		array := strings.Split(v, ":")
		mp[strings.TrimSpace(array[0])] = strings.TrimSpace(array[1])
		array = nil
	}

	// 返却する構造体を作成
	var new_client HttpClient
	new_client.url = u
	new_client.method = method
	new_client.requestHeader = mp
	switch method {
	case "GET", "DELETE":
		new_client.requestBody = nil
		delete(new_client.requestHeader, "Content-Type")
	case "POST", "PUT", "PATCH":
		new_client.requestHeader["Content-Type"] = "application/json"
		if d == "" {
			err = errors.New("dataが空")
		}
		new_client.requestBody = &d
	}

	if err != nil {
		return nil, err
	}

	return &new_client, err
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
	//HttpClient生成
	cli := new(http.Client)
	//リクエストボディ生成
	var reqBody io.Reader
	if c.requestBody == nil {
		//reqBody = nil
	} else {
		reqBody = bytes.NewBufferString(*c.requestBody)
	}
	//リクエスト生成
	req, err := http.NewRequest(c.method, c.url.String(), reqBody)
	if err != nil {
		return nil, nil, err
	}
	//リクエストヘッダをセット
	for k, v := range c.requestHeader {
		req.Header.Set(k, v)
	}
	//HTTPリクエスト実行
	res, err := cli.Do(req)
	if err != nil {
		return nil, nil, err
	}
	return req, res, err
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
