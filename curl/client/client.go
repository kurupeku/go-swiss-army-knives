package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HttpClient struct {
	url           *url.URL
	method        string
	requestBody   *string
	requestHeader map[string]string
}

func (c *HttpClient) Execute() error {
	req, err := c.newRequest()
	if err != nil {
		return err
	}

	return sendRequest(req)
}

// TODO:URL, HTTPメソッド, リクエストヘッダ, リクエストボディが適切に設定された*http.Requestを返却
func (c *HttpClient) newRequest() (*http.Request, error) {
	var body io.Reader
	if c.requestBody != nil {
		body = bytes.NewBufferString(*c.requestBody)
	}
	req, err := http.NewRequest(c.method, c.url.String(), body)
	if err != nil {
		return nil, err
	}
	for k, v := range c.requestHeader {
		req.Header.Set(k, v)
	}

	return req, nil
}

// TODO:HTTPリクエストを実行してレスポンスを取得
// TODO:リクエストURL,HTTPメソッド,リクエストヘッダを標準出力
// TODO:レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを標準出力
func sendRequest(req *http.Request) error {
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	outputRequest(req)
	outputResponse(res)

	return nil
}

func outputRequest(req *http.Request) {
	fmt.Println("")
	fmt.Println("===Request===")
	fmt.Printf("[URL] %s\n", req.URL.String())
	fmt.Printf("[Method] %s\n", req.Method)
	fmt.Println("[Headers]")
	for k, v := range req.Header {
		fmt.Printf("  %s: %s\n", k, strings.Join(v, "; "))
	}
}

func outputResponse(res *http.Response) {
	fmt.Println("")
	fmt.Println("===Response===")
	fmt.Printf("[Status] %d\n", res.StatusCode)
	fmt.Println("[Headers]")
	for k, v := range res.Header {
		fmt.Printf("  %s: %s\n", k, strings.Join(v, "; "))
	}
	fmt.Println("[Body]")
	b, _ := io.ReadAll(res.Body)
	fmt.Printf("%s\n", string(b))
}
