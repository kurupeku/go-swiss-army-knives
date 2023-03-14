package cmd

import (
	"context"
	"errors"
	"io"
	"logtransfer/input"
	"logtransfer/logs"
	"logtransfer/output"
	"logtransfer/storage"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

const (
	timeSpan   = 5
	channelLen = 10
)

var rootCmd = &cobra.Command{
	Use:   "logtransfer",
	Short: "A log transfer application over HTTP#POST",
	Long: `A log transfer application over HTTP#POST.
The application consists of a distributed system
with multi-threaded safe transfers.

- Args1 : request url
- Args2~: executable command

e.g ) logtransfer https://sample.com sh ./sample.sh`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("request url and command are required")
		}

		u, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		subCmd := exec.Command(args[1], args[2:]...)
		stdout, err := subCmd.StdoutPipe()
		if err != nil {
			return err
		}

		ctx, cancel := NewCtx()
		defer cancel()

		StartBackgrounds(ctx, u, stdout)

		subCmd.Run()
		return nil
	},
}

// TODO: シグナル（SIGTERM など）が呼ばれた際に、それを検知してキャンセル処理が走る context.Context を用意する
// TODO: context.CancelFunc も同時に返す
func NewCtx() (context.Context, context.CancelFunc) {
	// TODO: 3 週目：1 ~ 2 週目の処理を別スレッドで実行しつつ、シグナルを受け取った際にそれらを安全に終了させるメイン処理
	// https://zenn.dev/nekoshita/articles/dba0a7139854bb#signal.notifycontext%E3%82%92%E4%BD%BF%E3%81%A3%E3%81%9F%E6%9B%B8%E3%81%8D%E6%96%B9
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	return ctx, stop
}

// TODO2: 渡す channel のサイズは定数 channelLen を使用して定義する
// TODO3: 各関数に渡す context.Context は引数 ctx context.Context を使用する
// TODO4: 標準出力は r io.Reader として渡される
// TODO5: storage.Load() の実行間隔は定数 timeSpan を利用して渡す
// TODO6: output.Forward() の送信先 URL は引数 u *url.URL を使用して渡す
func StartBackgrounds(ctx context.Context, u *url.URL, r io.Reader) {
	// TODO: 3 週目：1 ~ 2 週目の処理を別スレッドで実行しつつ、シグナルを受け取った際にそれらを安全に終了させるメイン処理

	ln := make(chan []byte)
	out := make(chan []byte)
	errc := make(chan error)

	// TODO1: すべての処理を goroutine にて発火させる
	// TODO: 1 週目：標準出力（`io.Reader` として受け取る）から出力内容を読み取る処理と、読み取った結果を内部のバッファに保存する処理
	go input.Monitor(ctx, ln, errc, r)
	// TODO: 1 週目：標準出力（`io.Reader` として受け取る）から出力内容を読み取る処理と、読み取った結果を内部のバッファに保存する処理
	go storage.Listen(ctx, ln, errc)
	// TODO: 2 週目：内部バッファに保存された内容を一定時間ごとに読み込む処理と、読み取った文字列を Body とした HTTP#POST リクエストを投げる処理
	go storage.Load(ctx, out, errc, time.Second*5)
	// TODO: 2 週目：内部バッファに保存された内容を一定時間ごとに読み込む処理と、読み取った文字列を Body とした HTTP#POST リクエストを投げる処理
	go output.Forward(ctx, out, errc, u.String())

	go logs.Error(ctx, errc)

}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
