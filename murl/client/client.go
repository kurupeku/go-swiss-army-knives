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
	requestHeader map[string][]string
}

// TODO: NewHttpClientで必要な実装
// 1. URLの構築
//   - net/urlパッケージのParseRequestURIでURLを構築する
//
// 2. リクエストヘッダの設定
//
//   - customHeadersの内容をrequestHeaderマップに設定する
//
//   - 形式：map[string]string{"Header-Name": "header-value"}
//
//     3. HTTPメソッドとボディの設定
//     GET/DELETEの場合:
//
//   - requestBodyはnil
//
//   - Content-Typeヘッダーは含めない（存在する場合は削除）
//     POST/PUT/PATCHの場合:
//
//   - Content-Typeヘッダーを"application/json"に設定
//
//   - requestBodyにdataを設定
//
//   - dataが空文字列の場合はエラーを返す
func NewHttpClient(
	rawurl string,
	method string,
	data string,
	customHeaders []string,
) (*HttpClient, error) {
	// 各パラメータのバリデーション
	if err := validateRawURL(rawurl); err != nil {
		return nil, err
	}
	if err := validateMethod(method); err != nil {
		return nil, err
	}
	if err := validateHeader(customHeaders); err != nil {
		return nil, err
	}

	// URLの構築
	parsedURL, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, err
	}

	// リクエストヘッダの初期化
	headers := make(map[string][]string)

	// カスタムヘッダーの設定
	for _, header := range customHeaders {
		kv := strings.Split(strings.TrimSpace(header), ":")
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		headers[key] = append(headers[key], value)
	}

	var requestBody *string

	// メソッドに応じた処理
	switch method {
	case http.MethodGet, http.MethodDelete:
		// GET/DELETEの場合はボディなしでContent-Typeヘッダーを削除
		delete(headers, "Content-Type")
		requestBody = nil
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		// POST/PUT/PATCHの場合
		if data == "" {
			return nil, errors.New("request body is required")
		}
		// JSONデータのバリデーション
		if err := validateData(data); err != nil {
			return nil, err
		}
		headers["Content-Type"] = []string{"application/json"}
		requestBody = &data
	}

	return &HttpClient{
		url:           parsedURL,
		method:        method,
		requestBody:   requestBody,
		requestHeader: headers,
	}, nil
}

func (c *HttpClient) Execute() (string, string, error) {
	// リクエストの生成
	req, err := c.BuildRequest()
	if err != nil {
		return "", "", err
	}

	// http.DefaultClient を使ってリクエストを送信
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	// http.Response.Body は defer で必ず Close する
	defer res.Body.Close()

	return CreateRequestText(req), CreateResponseText(res), nil
}

func (c *HttpClient) BuildRequest() (*http.Request, error) {
	// URLの検証
	if c.url == nil || c.url.String() == "" {
		return nil, errors.New("invalid URL: empty URL")
	}

	// リクエストボディの作成
	var body io.Reader
	if c.requestBody != nil {
		body = strings.NewReader(*c.requestBody)
	}

	// リクエストの生成
	req, err := http.NewRequest(c.method, c.url.String(), body)
	if err != nil {
		return nil, err
	}

	// ヘッダーの設定
	for key, values := range c.requestHeader {
		req.Header[key] = values
	}

	return req, nil
}

// TODO:リクエストURL,HTTPメソッド,リクエストヘッダを所定のフォーマットで返却
func CreateRequestText(req *http.Request) string {
	var b strings.Builder

	// 最初に空行
	b.WriteString("\n===Request===\n")

	// URLを表示
	b.WriteString(fmt.Sprintf("[URL] %s\n", req.URL.String()))

	// メソッドを表示
	b.WriteString(fmt.Sprintf("[Method] %s\n", req.Method))

	// ヘッダーを表示
	b.WriteString("[Headers]\n")

	// ヘッダーキーを昇順でソート
	keys := sortedKeys(req.Header)

	// ヘッダーを表示
	for _, key := range keys {
		b.WriteString(fmt.Sprintf("  %s: %s\n", key, strings.Join(req.Header[key], "; ")))
	}

	return b.String()
}

// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func CreateResponseText(res *http.Response) string {
	var b strings.Builder

	// 最初に空行とタイトル
	b.WriteString("\n===Response===\n")

	// ステータスコードを表示
	b.WriteString(fmt.Sprintf("[Status] %d\n", res.StatusCode))

	// ヘッダーを表示
	b.WriteString("[Headers]\n")

	// ヘッダーキーを昇順でソート
	keys := sortedKeys(res.Header)

	// ヘッダーを表示
	for _, key := range keys {
		b.WriteString(fmt.Sprintf("  %s: %s\n", key, strings.Join(res.Header[key], "; ")))
	}

	// ボディを表示
	b.WriteString("[Body]\n")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	b.WriteString(string(body) + "\n")

	return b.String()
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
