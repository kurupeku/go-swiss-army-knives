package client

import (
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

// TODO: NewHttpClientで必要な実装
// 1. URLの構築
//   - net/urlパッケージのParseRequestURIでURLを構築する
//
// 2. リクエストヘッダの設定
//   - customHeadersの内容をrequestHeaderマップに設定する
//   - 形式：map[string]string{"Header-Name": "header-value"}
//
// 3. HTTPメソッドとボディの設定
//   GET/DELETEの場合:
//   - requestBodyはnil
//   - Content-Typeヘッダーは含めない（存在する場合は削除）
//   POST/PUT/PATCHの場合:
//   - Content-Typeヘッダーを"application/json"に設定
//   - requestBodyにdataを設定
//   - dataが空文字列の場合はエラーを返す
func NewHttpClient(
	rawurl string,
	method string,
	data string,
	customHeaders []string,
) (*HttpClient, error) {
	// TODO: 1 週目：HTTP 通信用クライアントを構築
	return nil, nil
}

func (c *HttpClient) Execute() (string, string, error) {
	req, res, err := c.SendRequest()
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	return CreateRequestText(req), CreateResponseText(res), nil
}

// TODO: SendRequestで必要な実装
// 1. http.NewRequestでリクエストを生成
//   - URL、メソッド、ボディを使用してリクエストを作成する
//
// 2. リクエストヘッダの設定
//   - requestHeaderの内容をhttp.Request.Headerに設定する
//
// 3. リクエストの実行
//   - http.DefaultClientを使用してリクエストを送信し、結果を取得する
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
