package client

import (
	"fmt"
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

	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Can't parse url: %s", rawurl)
	} else {
		fmt.Printf("[OK]URI= %v \n", url)
	}

	var rBody *string
	// method 判別
	switch method {
	case "GET", "DELETE":
		fmt.Print("method = GET,DELETE")
		rBody = nil
	case "POST", "PUT", "PATCH":
		fmt.Print("method = POST,PUT,PATCH")
		rBody = &data
	default:
		return nil, fmt.Errorf("[ERROR] Unauthorized method: %s", method)
	}

	client := HttpClient{
		url:         url,
		method:      method,
		requestBody: rBody,
		//requestHeader: rHeader, #TODO
	}

	// HttpClientのポインターで戻り値を返す ,
	// errorはnil(正常終了)で戻り値を返す
	return &client, nil
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
