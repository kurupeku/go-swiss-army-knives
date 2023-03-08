# 課題 2:ログ転送アプリの実装

## コマンド仕様

### 説明

- 標準出力を内部バッファに一時保存し、5 秒毎にバッファ内の文字列を HTTP#POST で転送する CLI アプリ
- 5 秒ごとに転送がかかる（その間ログが出ていない場合はリクエストしない）
- POST 時の Body は出力内容をそのまま `plain/text` で送信する
- アプリケーションの動作は与えられたコマンドと別スレッドの goroutine として実行され、 channel を通して分散動作する

### コマンド引数・フラグ

- 第一引数は転送先の URL を 設定する
- 第二引数以降はログを吐く何かしらのコマンドを設定する
- フラグはデフォルトの `--help` のみ対応

### 使用例

```bash
logtransfer https://sample.com sh ./sample.sh

# 標準出力がインタラプトされなにも表示されない
```

## 実装課題

- コマンド引数を受け取り、その実行結果の標準出力をインタラプトする部分は実装済み
- goroutine 内で起こったエラーをログファイルに書き出す処理は実装済み
- 未実装となっている以下の 6 つの処理を実装する
  - [1 週目] 標準出力（`io.Reader` として受け取る）から出力内容を読み取る処理と、読み取った結果を内部のバッファに保存する処理
  - [2 週目] 内部バッファに保存された内容を一定時間ごとに読み込む処理と、読み取った文字列を Body とした HTTP#POST リクエストを投げる処理
  - [3 週目] 1 ~ 2 週目の処理を別スレッドで実行しつつ、シグナルを受け取った際にそれらを安全に終了させるメイン処理

### 1 週目：標準出力（`io.Reader` として受け取る）から出力内容を読み取る処理と、読み取った結果を内部のバッファに保存する処理

- 対応ファイル： `logtransfer/input/watcher.go` および `logtransfer/storage/buffer.go`
- 実装内容： `func Monitor(ctx context.Context, ln chan []byte, errc chan error, r io.Reader)`
  で標準出力を転送し `func Listen(ctx context.Context, ln chan []byte, errc chan error)` で内部のバッファに保存
- 実装対象メソッド・実装条件
  - `func Monitor(ctx context.Context, ln chan []byte, errc chan error, r io.Reader)`
    - 引数 `r io.Reader` とｓて標準出力が渡されてくるので、入力を待ち受ける
    - 入力があった場合は 1 行だけ読み込み、その文字列を引数 `ln chan []byte` へ送信した後、待受状態に戻る
    - `ctx context.Context` がキャンセルされた場合には `ln` を close し、速やかに関数を終了する
    - エラーが発生した際には引数 `errc chan error` へエラーを送信する
  - `func Listen(ctx context.Context, ln chan []byte, errc chan error)`
    - 引数 `ln chan []byte` で文字列を受信した際に、グローバル変数 `buf *bytes.Buffer` へ書き込む
    - `ctx context.Context` がキャンセルされた場合には速やかに関数を終了する
    - エラーが発生した際には `errc chan error` へエラーを送信する

### 2 週目：内部バッファに保存された内容を一定時間ごとに読み込む処理と、読み取った文字列を Body とした HTTP#POST リクエストを投げる処理

- 対応ファイル： `logtransfer/storage/buffer.go` および `logtransfer/output/http.go`
- 実装内容： `func Load(ctx context.Context, out chan []byte, errc chan error, span time.Duration)`
  で一定時間ごとにバッファを読み込み `func Forward(ctx context.Context, out chan []byte, errc chan error, url string)` リクエストとして送信する
- 実装対象メソッド・実装条件
  - `func Load(ctx context.Context, out chan []byte, errc chan error, span time.Duration)`
    - グローバル変数 `buf *bytes.Buffer` から一定時間ごとに内容を読み込み、内容を引数 `out chan []byte` へ送信する
    - 読み込む間隔は引数 `span time.Duration` を利用して制御する
    - `buf` に何も保存されていなければ内容の送信は行わない
    - 一度に保存された内容すべてを読み取り、 `buf` にはなにも保存されていない状態にリセットする
    - `ctx context.Context` がキャンセルされた場合には `out` を close し、速やかに関数を終了する
    - エラーが発生した際には `errc chan error` へエラーを送信する
  - `func Forward(ctx context.Context, out chan []byte, errc chan error, url string)`
    - 引数 `out chan []byte` で文字列を受信した際に、その内容 Body として引数 `url string` への HTTP#POST リクエストを行う
    - `Content-Type: plain/text` を Header に添えて送信を行う
    - `ctx context.Context` がキャンセルされた場合には速やかに関数を終了する
    - エラーが発生した際には `errc chan error` へエラーを送信する

### 3 週目：1 ~ 2 週目の処理を別スレッドで実行しつつ、シグナルを受け取った際にそれらを安全に終了させるメイン処理

- 対応ファイル： `logtransfer/cmd/root.go`
- 実装内容: `func NewCtx() (context.Context, context.CancelFunc)`
  でシグナルを検知した際にキャンセルが伝播するコンテキストを用意し, `func StartBackgrounds(ctx context.Context, u *url.URL, r io.Reader)` にて 1 ~
  2 週目に実装したメソッド郡を goroutine で起動する
- 実装対象メソッド・実装条件
  - `func NewCtx() (context.Context, context.CancelFunc)`
    - シグナル（`SIGTERM` など）が呼ばれた際に、それを検知してキャンセル処理が走る `context.Context` を用意する
    - `context.CancelFunc` も同時に返す
  - `func StartBackgrounds(ctx context.Context, u *url.URL, r io.Reader)`
    - すべての処理を goroutine にて発火させる
    - 渡す channel のサイズは定数 `channelLen` を使用して定義する
    - 各関数に渡す `context.Context` は引数 `ctx context.Context` を使用する
    - 標準出力は `r io.Reader` として渡される
    - `storage.Load()` の実行間隔は定数 `timeSpan` を利用して渡す
    - `output.Forward()` の送信先 URL は引数 `u *url.URL` を使用して渡す

## 動作プレビュー

CLI 完成後、以下の手順で動作を確認できます。

1. `docker compose up -d` でコンテナを立ち上げる
2. [http://localhost:3000](http://localhost:3000) にアクセスする
3. ローカル or コンテナ内 or DevContainer 内のプロジェクトルートで以下のコマンドを実行する

- ローカル
  - `task preview_lt`
- コンテナ内 or DevContainer 内
  - `task preview_lt_docker`