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

// HTTP リクエストを送信するためのクライアントを生成する関数
func NewHttpClient(
	rawurl string,
	method string,
	data string,
	customHeaders []string,
) (*HttpClient, error) {
	// TODO: Implement here
	// HTTPクライアントを初期化する処理を実装してください：
	// 1. URLの構築
	//    - url パッケージを利用して `url.URL` オブジェクトを生成
	//
	// 2. リクエストヘッダの構築
	//    - カスタムヘッダー []string をパースして map[string][]string に変換
	// .  - map[string][]string のキーはヘッダー名、値はヘッダー値のリスト
	//
	// 3. メソッドに応じたリクエストボディとヘッダーの設定
	//    - GET, DELETE の場合:
	//      - リクエストボディは不要 (nil)
	//      - Content-Type ヘッダーを削除
	//    - POST, PUT, PATCH の場合:
	//      - リクエストボディが必須（空文字列の場合はエラー）
	//      - Content-Type ヘッダーを "application/json" のみに上書き
	//
	// ヒント：
	// - http パッケージには標準的な HTTP メソッドが定数として定義されています
	// - switch 文を使用してメソッドごとの処理を分岐させることができます

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
		// GET / DELETE の場合はボディなしでContent-Typeヘッダーを削除
		delete(headers, "Content-Type")
		requestBody = nil
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		// POST / PUT / PATCH の場合は Content-Type ヘッダーとリクエストボディが必要
		if data == "" {
			return nil, errors.New("request body is required")
		}
		headers["Content-Type"] = []string{"application/json"}
		requestBody = &data
	default:
		return nil, errors.New("unsupported HTTP method")
	}

	return &HttpClient{
		url:           parsedURL,
		method:        method,
		requestBody:   requestBody,
		requestHeader: headers,
	}, nil
}

func (c *HttpClient) Execute() (string, string, error) {
	// TODO: Implement here
	// HTTPリクエストを実行し、リクエストとレスポンスの文字列表現を返却する処理を実装してください：
	// 1. BuildRequest() を使用してリクエストを生成
	// 2. リクエスト実行とエラーハンドリング
	//    - http.DefaultClient.Do(req) でリクエストを送信
	//    - エラーの場合は空文字列とエラーを返却
	// 3. リソースのクリーンアップ
	//    - レスポンスボディは必ず Close すること（defer を使用）
	// 4. レスポンスの文字列化
	//    - CreateRequestText() でリクエスト情報を文字列化
	//    - CreateResponseText() でレスポンス情報を文字列化
	//
	// ヒント：
	// - http.Response.Body は io.ReadCloser インターフェースを実装しています
	// - レスポンスボディの Close 漏れは重大なリソースリークの原因となります

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

// リクエストURL,HTTPメソッド,リクエストヘッダを所定のフォーマットで返却
func CreateRequestText(req *http.Request) string {
	// TODO: Implement here
	// 以下の順序で情報を文字列化してください：
	// 1. リクエストセクションの開始
	//    - "\n===Request===\n" で開始（前後に改行）
	// 2. URL
	//    - [URL] の後に req.URL.String() の値を表示
	// 3. HTTPメソッド
	//    - [Method] の後に req.Method の値を表示
	// 4. ヘッダー
	//    - [Headers] の次行にヘッダー情報を一覧表示
	//    - ヘッダーキーを昇順でソート（sortedKeys()関数を使用）
	//    - 各ヘッダーは "  {key}: {value}" の形式で表示（先頭に2つの空白）
	//    - ヘッダー値が複数ある場合は "; " で結合
	//
	// ヒント：
	// - 文字列の操作には strings パッケージなどを利用すると便利です

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

// レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func CreateResponseText(res *http.Response) string {
	// TODO: Implement here
	// 以下の順序で情報を文字列化してください：
	// 1. レスポンスセクションの開始
	//    - "\n===Response===\n" で開始（前後に改行）
	// 2. ステータスコード
	//    - [Status] の後に res.StatusCode を表示
	// 3. ヘッダー
	//    - [Headers] の次行にヘッダー情報を一覧表示
	//    - ヘッダーキーを昇順でソート（sortedKeys()関数を使用）
	//    - 各ヘッダーは "  {key}: {value}" の形式で表示（先頭に2つの空白）
	//    - ヘッダー値が複数ある場合は "; " で結合
	// 4. ボディ
	//    - [Body] の次行にレスポンスボディを表示
	// .  - ボディが空の場合は改行のみを表示
	//
	// ヒント：
	// - 文字列の操作には strings パッケージなどを利用すると便利です

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
