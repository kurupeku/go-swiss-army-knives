package client

import (
	"bytes"
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
	var httpClient HttpClient
	url, _ := url.ParseRequestURI(rawurl)

	httpClient.url = url
	httpClient.method = method

	switch method {
	case "GET", "DELETE":
		httpClient.requestBody = nil
		httpClient.requestHeader = map[string]string{}
		for _, customHeader := range customHeaders {
			cH := strings.Split(customHeader, ":")
			if cH[0] == "Content-Type" {
				continue
			}
			httpClient.requestHeader[cH[0]] = strings.TrimSpace(cH[1]);
		}
	case "POST", "PUT", "PATCH":
		if data == "" {
			return nil, errors.New("エラー")
		}
		httpClient.requestBody = &data
		httpClient.requestHeader = map[string]string{}
		for _, customHeader := range customHeaders {
			cH := strings.Split(customHeader, ":")
			httpClient.requestHeader[cH[0]] = strings.TrimSpace(cH[1]);
		}
		httpClient.requestHeader["Content-Type"] = "application/json";
	}

	return &httpClient, nil
}

func (c *HttpClient) Execute() (string, string, error) {
	req, res, err := c.SendRequest()
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	return CreateRequestText(req), CreateResponseText(res), nil
}

func createBody(c *HttpClient) (io.Reader) {
	if c.requestBody != nil { return strings.NewReader(*c.requestBody) }
	return nil;
}

func createHeader(c *HttpClient) (map[string][]string) {
	requestHeader := map[string][]string{}
	for key, value := range c.requestHeader {
		requestHeader[key] = append(requestHeader[key], value)
	}
	return requestHeader
}

// TODO:URL, HTTPメソッド, リクエストヘッダ, リクエストボディが適切に設定された*http.Requestを生成
// TODO:HTTPリクエストを実行後の*http.Request, *http.Responseを返却
// TODO:ただ単にオブジェクトを作るだけでなく、このメソッド内でリクエストの実行も完了させる
func (c *HttpClient) SendRequest() (*http.Request, *http.Response, error) {
	// リクエストの生成
	req, _ := http.NewRequest(c.method, c.url.String(), createBody(c))
	req.Header = createHeader(c)

	// httpリクエスト
	response, err := new(http.Client).Do(req)
	if err != nil {
		return nil, nil, err
	}

	return req, response, err
}

func createRequestTitle() string {
	return "\n===Request===\n"
}

func createRequestUrlText(url string) string {
	return "[URL] " + url + "\n"
}

func createRequestMethodText(method string) string {
	return "[Method] " + method + "\n"
}

func getSortedKeys(header map[string][]string) []string {
	var keys []string
	for key := range header {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func getSortedHeaderByKey(header map[string][]string) map[string][]string {
	sortedHeader := map[string][]string{}
	for _, key := range getSortedKeys(header) {
		sortedHeader[key] = header[key]
	}
	return sortedHeader
}

func createHeaderContentText(headerName string, values []string) string {
	headerContentText := ""
	for _, contents := range values {
		// ヘッダータイトル
		headerContentText += "  " + headerName + ": "
		// ヘッダー詳細
		contents_array := strings.Split(contents, ",")
		headerContentText += contents_array[0]
		for _, content := range contents_array[1:] {
			headerContentText += ";" + content
		}
		// 改行
		headerContentText += "\n"
	}
	return headerContentText
}

func createHeaderText(header map[string][]string) string {
	headerText := "[Headers]\n"
	for headerName, values := range getSortedHeaderByKey(header) {
		headerText += createHeaderContentText(headerName, values)
	}
	return headerText
}


// TODO:リクエストURL,HTTPメソッド,リクエストヘッダを所定のフォーマットで返却
func CreateRequestText(req *http.Request) string {
	var requestText string
	requestText = createRequestTitle()
	requestText += createRequestUrlText(req.URL.String())
	requestText += createRequestMethodText(req.Method)
	requestText += createHeaderText(req.Header)
	return requestText
}

func createResponseTitle() string {
	return "\n===Response===\n"
}

func createResponseStatusText(status string) string {
	return "[Status] " + status + "\n"
}

func createResponseBodyText(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return "[Body]\n" + buf.String() + "\n"
}

// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
func CreateResponseText(res *http.Response) string {
	var responseText string
	responseText = createResponseTitle()
	responseText += createResponseStatusText(strconv.Itoa(res.StatusCode))
	responseText += createHeaderText(res.Header)
	responseText += createResponseBodyText(res.Body)
	return responseText
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
