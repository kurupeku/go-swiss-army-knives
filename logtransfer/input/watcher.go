package input

import (
	"context"
	"io"
)

// r から受け取ったデータを 1行ずつ ln へ送信する関数
func Monitor(ctx context.Context, ln chan []byte, errc chan error, r io.Reader) {
	// TODO: Implement here
	// 以下の処理を実装する必要があります：
	// 1. 入力ストリームの監視
	//    - 継続的にデータを読み取る
	//    - 1行ずつ処理する
	//
	// 2. チャネルを介したデータの転送
	//    - 読み取ったデータを ln チャネルへ送信する
	//    - エラー発生時は errc チャネルへ通知する
	//
	// 3. コンテキストによる制御
	//    - コンテキストのキャンセルを検知する
	//    - キャンセル時は適切に処理を終了する
	//
	// ヒント：
	// - io.Reader からの効率的な読み取りには bufio パッケージが有用です
	// - コンテキストによる制御は select 文と組み合わせることで実現できます
}
