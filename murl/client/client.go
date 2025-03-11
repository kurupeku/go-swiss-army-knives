package client

import (
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
	return nil, nil
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
	return "", "", nil
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
	return ""
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
