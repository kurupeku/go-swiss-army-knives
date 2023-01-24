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
	return signal.NotifyContext(context.Background(), os.Interrupt)
}

// TODO: すべての処理を goroutine にて発火させる
// TODO: 渡す channel のサイズは定数 channelLen を使用して定義する
// TODO: 各関数に渡す context.Context は引数 ctx context.Context を使用する
// TODO: 標準出力は r io.Reader として渡される
// TODO: storage.Load() の実行間隔は定数 timeSpan を利用して渡す
// TODO: output.Forward() の送信先 URL は引数 u *url.URL を使用して渡す
func StartBackgrounds(ctx context.Context, u *url.URL, r io.Reader) {
	var ln, out, errc = make(chan []byte, channelLen), make(chan []byte, channelLen), make(chan error, channelLen)

	go input.Monitor(ctx, ln, errc, r)
	go storage.Listen(ctx, ln, errc)
	go storage.Load(ctx, out, errc, timeSpan*time.Second)
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
