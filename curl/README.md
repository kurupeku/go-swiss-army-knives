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
  Vary: Accept-Encoding
  Expires: Wed, 11 Jan 2023 08:26:15 GMT
  Content-Type: text/html; charset=UTF-8
  Date: Wed, 04 Jan 2023 08:26:15 GMT
  Last-Modified: Thu, 17 Oct 2019 07:18:26 GMT
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
  Cache-Control: max-age=604800
  Last-Modified: Thu, 17 Oct 2019 07:18:26 GMT
  Content-Type: text/html; charset=UTF-8
  Expires: Wed, 11 Jan 2023 08:33:43 GMT
  Accept-Ranges: bytes
  Date: Wed, 04 Jan 2023 08:33:43 GMT
  Content-Length: 1256
[Body]
<!doctype html>
<html>省略</html>
```

## 実装課題

- コマンド引数・フラグを受け取る部分は実装済み
- 通信結果を出力する部分は実装済み
- 未実装となっている以下の 3 つの処理を実装する
  - [1 週目] コマンド引数・フラグの入力値の妥当性チェック
  - [2 週目] HTTP 通信用クライアントを構築
  - [3 週目] HTTP 通信を実行して通信結果のテキストを構築

### 1 週目：コマンド引数・フラグの入力値の妥当性チェック

- 対応ファイル：`curl/client/builder.go`
- 実装対象型：`HttpClientBuilder`
- 実装内容：`func (b *HttpClientBuilder) Validate() error`内部から呼ばれている 4 つの妥当性チェックメソッドを実装する
- 実装対象メソッド・実装条件
  - `func (b *HttpClientBuilder) validateRawURL() error`
    - `b.rawurl`について以下のチェックを行い、違反している場合は`error`を返却
      - 正しい URL のフォーマットになっている
        - `net/url`パッケージの`ParseRequestURI`でエラーが起きなければ OK
      - プロトコルが http または https になっている
  - `func (b *HttpClientBuilder) validateMethod() error`
    - `b.method`が許容されている HTTP メソッドに一致しなければ`error`を返却
      - 未指定の場合のデフォルト値`GET`も既に設定された状態なので、`b.method`が空文字になっている場合もエラーにしてください
  - `func (b *HttpClientBuilder) validateData() error`
    - `b.data`について以下の状態と一致するかのチェックを行い、いずれにも該当しない場合は`error`を返却
      - 値が設定されていない（空文字）
      - 正しい JSON 形式の文字列になっている
        - `encoding/json`パッケージを使うと正しい JSON 形式の文字列であることの確認が簡単になります
  - `func (b *HttpClientBuilder) validateHeader() error`
    - `b.customHeaders`の全ての要素が以下の条件を満たすことを確認し、違反している場合は`error`を返却
      - `:`が 1 つだけ含まれており、`:`の前後が空ではない

### 2 週目：HTTP 通信用クライアントを構築

- 対応ファイル：`curl/client/builder.go`
- 実装対象型：`HttpClientBuilder`
- 実装内容：`func (b *HttpClientBuilder) Build() (*HttpClient, error)`で`*HttpClient`型のインスタンスを構築して返却する
- 実装対象メソッド・実装条件
  - `func (b *HttpClientBuilder) Build() (*HttpClient, error)`
    - URL は net/url パッケージの\*url.URL で構築する
    - b.customHeaders をリクエストヘッダとして設定
    - HTTP メソッドが GET,DELETE の場合
      - リクエストボディは設定しない
      - リクエストヘッダに Content-Type が含まれている場合は削除
    - HTTP メソッドが POST,PUT,DELETE の場合
      - リクエストヘッダの Content-Type は"application/json"にする
      - b.data の値をそのままレスポンスボディに設定
        - その際、b.data が空であればエラー

### 3 週目：HTTP 通信を実行してリクエストとレスポンスの内容を返却

- 対応ファイル：`curl/client/client/go`
- 実装対象型：`HttpClient`
- 実装内容：`func (c *HttpClient) Execute() (string, string, error)`でリクエストを送信してリクエスト、レスポンスの内容を`string`で返却する
- 実装対象メソッド・実装条件
  - `func (c *HttpClient) newRequest() (*http.Request, error)`
    - URL, HTTP メソッド, リクエストヘッダ, リクエストボディが適切に設定された`*http.Request`を返却
  - `func sendRequest(req *http.Request) (request string, response string, e error)`
    - HTTP リクエストを実行してレスポンスを取得
    - リクエスト URL,HTTP メソッド,リクエストヘッダを所定のフォーマットで返却
    - レスポンスのステータスコード,レスポンスヘッダ,レスポンスボディを所定のフォーマットで返却
- リクエスト内容のフォーマット
  - 改行コードは`\n`
  - 最初に空行を 1 行入れる
  - 以下の形式で URL, Method, Headers を入れる
    - Headers はスペース 2 つでインデントをつける
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
