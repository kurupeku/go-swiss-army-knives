# 課題 1:簡易 curl コマンドを実装する

## コマンド仕様

### 説明

- http または https リクエストを投げてレスポンスを標準出力する
- 対応しているメソッドは GET, POST, PUT, DELETE, PATCH の 5 種
- GET, DELETE の場合は URL のみ指定可能でリクエストボディによるデータ送信は行わない
- POST, PUT, PATCH の場合は Content-Type を`application/json`固定とし、リクエストボディに JSON データを設定する
- リクエストヘッダはどのメソッドの場合も指定可能

### コマンド引数・フラグ

- コマンド引数はリクエストを送信する URL を 1 つだけ設定する
- フラグ（`-a`または`--aaa`という形式で設定するコマンドオプション）は以下の通り
  - `-X`(`--request`): HTTP メソッドを指定(無指定の場合のデフォルトは"GET")
  - `-d`(`--data`): リクエストボディを指定(POST,PUT,PATCH の場合だけ送信される。JSON 形式のみ許容)
  - `-H`(`--header`): "ヘッダ名:値"の形式で記述するリクエストヘッダ。複数個指定可能
- 以下はコマンドのヘルプ表示

```bash
$ go run main.go -h
curl is http/https client command.
- Args: URL
- Available HTTP Methods: GET, POST, PUT, DELETE, PATCH
- Available Content-Type: application/json(only for POST, PUT, PATCH)

Usage:
  curl [URL] [flags]

Flags:
  -d, --data string          HTTP Post, Put, Patch Data
  -H, --header stringArray   Pass custom header(s) to server
  -h, --help                 help for curl
  -X, --request string       HTTP method (default "GET")
```

### 使用例

#### GET(リクエストヘッダーを付与)

```bash
$ go run main.go http://example.com -H Sample:test

===Request===
[URL] http://example.com
[Method] GET
[Headers]
  Sample: test


===Response===
[Status] 200
[Headers]
  Age: 374145
  Cache-Control: max-age=604800
  Content-Type: text/html; charset=UTF-8
  Date: Wed, 04 Jan 2023 08:26:15 GMT
  Expires: Wed, 11 Jan 2023 08:26:15 GMT
  Last-Modified: Thu, 17 Oct 2019 07:18:26 GMT
  Vary: Accept-Encoding
[Body]
<!doctype html>
<html>省略</html>
```

#### POST(リクエストボディを付与)

```bash
$ go run main.go http://example.com -X POST -d '{"id":1}'

===Request===
[URL] http://example.com
[Method] POST
[Headers]
  Content-Type: application/json


===Response===
[Status] 200
[Headers]
  Accept-Ranges: bytes
  Cache-Control: max-age=604800
  Content-Length: 1256
  Content-Type: text/html; charset=UTF-8
  Date: Wed, 04 Jan 2023 08:33:43 GMT
  Expires: Wed, 11 Jan 2023 08:33:43 GMT
  Last-Modified: Thu, 17 Oct 2019 07:18:26 GMT
[Body]
<!doctype html>
<html>省略</html>
```

## 実装課題

- コマンド引数・フラグを受け取る部分は実装済み
- 通信結果を出力する部分は実装済み
- 未実装となっている以下の 3 つの処理を実装する
  - [1 週目] HTTP 通信用クライアントを構築
  - [2 週目] HTTP 通信を実行
  - [3 週目] HTTP 通信結果のテキストを構築

### 1 週目：HTTP 通信用クライアントを構築

- 対応ファイル：`curl/client/client.go`
- 実装内容：`func NewHttpClient(rawurl string, method string, data string, customHeaders []string) (*HttpClient, error)`で `*HttpClient`のインスタンス生成して返却
- 実装対象メソッド・実装条件
  - `func NewHttpClient(rawurl string, method string, data string, customHeaders []string) (*HttpClient, error)`
    - `url`フィールド は net/url パッケージの\*url.URL で構築する
    - `customHeaders`引数の要素を`:`で区切って、`requestHeader`フィールドのキーと値に設定
    - HTTP メソッドが GET,DELETE の場合
      - リクエストボディ(`requestBody`フィールド)は`nil`
      - リクエストヘッダに `Content-Type` が含まれている場合は削除
    - HTTP メソッドが POST,PUT,PATCH の場合
      - リクエストヘッダの `Content-Type` は"application/json"にする
      - data の値をそのままレスポンスボディ(`requestBody`フィールド)に設定
        - その際、data が空であればエラー

### 2 週目：HTTP 通信を実行

- 対応ファイル：`curl/client/client/go`
- 実装対象型：`HttpClient`
- 実装内容：`func (c *HttpClient) SendRequest() (*http.Request, *http.Response, error)`でリクエストを送信して`*http.Request`, `*http.Response`を返却する
- 実装対象メソッド・実装条件
  - `func (c *HttpClient) SendRequest() (*http.Request, *http.Response, error)`
    - URL, HTTP メソッド, リクエストヘッダ, リクエストボディが適切に設定された`*http.Request`を構築
    - HTTP リクエストを実行してレスポンスを取得
    - `*http.Request`, `*http.Response`を返却
    - `呼び出し元でCloseしているので、レスポンスボディのクローズは不要`

### 3 週目：HTTP 通信結果のテキストを構築

- 対応ファイル：`curl/client/client/go`
- 実装内容:`*http.Request`, `*http.Response`の内容を適切なフォーマットの`string`に整形する
- 実装対象メソッド・実装条件
  - `func CreateRequestText(req *http.Request) string`
    - リクエスト URL,HTTP メソッド,リクエストヘッダを所定のフォーマットで返却
  - `func CreateResponseText(res *http.Response) string`
    - レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
- リクエスト内容のフォーマット
  - 改行コードは`\n`
  - 最初に空行を 1 行入れる
  - 以下の形式で URL, Method, Headers を入れる
    - Headers はスペース 2 つでインデントをつける
    - Headers が複数ある場合は Key が昇順にソートされた状態で表示する
    - 一つの Key に対して値が複数ある場合は `;<半角スペース>` で区切る
      - e.g.) `Content-Type: text/html; charset=UTF-8`
    - Headers で表示するリクエストヘッダがなくても、`[Headers]`という行は必ず入れる

```bash

===Request===
[URL] https://example.com
[Method] GET
[Headers]
  Connection: keep-alive
```

- レスポンス内容のフォーマット
  - 改行コードは`\n`
  - 最初に空行を 1 行入れる
  - 以下の形式で Status, Headers, Body を入れる
    - Headers はスペース 2 つでインデントをつける
    - Headers が複数ある場合は Key が昇順にソートされた状態で表示する
    - 一つの Key に対して値が複数ある場合は `;<半角スペース>` で区切る
      - e.g.) `Content-Type: text/html; charset=UTF-8`
    - Headers で表示するリクエストヘッダがなくても、`[Headers]`という行は必ず入れる
  - Body はインデントなしで出力し、最後に改行を入れる
    - レスポンスボディが空の場合は、空行を出力する

```bash

===Response===
[Status] 200
[Headers]
  Content-Type: application/json
[Body]
"{\"status\":\"ok\"}"
```
