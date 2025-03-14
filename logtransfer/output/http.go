package output

import (
	"context"
)

const (
	contentType = "text/plain"
)

// Forward は out から受け取ったデータを指定の url へ HTTP POST する関数
func Forward(ctx context.Context, out chan []byte, errc chan error, url string) {
	// TODO: Implement here
	// 以下の処理を実装する必要があります：
	// 1. データの受信と HTTP リクエストの作成
	//    - out チャネルからデータを受信する
	//    - チャネルのクローズを適切に検知する
	//    - 受信したデータを使用して POST リクエストを作成する
	//
	// 2. HTTP リクエストの送信と制御
	//    - コンテキストを使用してリクエストをキャンセル可能にする
	//    - Content-Type ヘッダーを適切に設定する
	//
	// 3. レスポンスの処理
	//    - ステータスコードの確認
	//    - エラー時にはレスポンスボディを読み取る
	//
	// 4. エラーハンドリング
	//    - リクエスト作成時のエラー処理
	//    - 送信時のエラー処理
	//    - レスポンスステータスが 400 以上の場合の処理
	//
	// ヒント：
	// - http パッケージに context を渡すことのできる関数があります
	// - []byte を io.Reader に変換する便利なメソッドが bytes パッケージなどに存在します
}
