package cmd

import (
	"context"
	"errors"
	"io"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"logtransfer/input"
	"logtransfer/logs"
	"logtransfer/output"
	"logtransfer/storage"

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

		ctx, cancel := NewCtx()
		defer cancel()

		subCmd := exec.CommandContext(ctx, args[1], args[2:]...)
		stdout, err := subCmd.StdoutPipe()
		if err != nil {
			return err
		}

		StartBackgrounds(ctx, u, stdout)

		if err := subCmd.Run(); err != nil {
			// コンテキストがキャンセルされた場合（Ctrl+Cが押された場合）は
			// エラーとして扱わない
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		return nil
	},
}

// シグナル（SIGTERM など）が呼ばれた際に、それを検知してキャンセル処理が走る context.Context を用意する関数
func NewCtx() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt)
}

// コアロジック群をバックグラウンドでの処理を開始する関数
func StartBackgrounds(ctx context.Context, u *url.URL, r io.Reader) {
	var (
		ln   = make(chan []byte, channelLen)
		out  = make(chan []byte, channelLen)
		errc = make(chan error, channelLen)
	)

	go logs.Error(ctx, errc)
	go input.Monitor(ctx, ln, errc, r)
	go storage.Listen(ctx, ln, errc)
	go storage.Load(ctx, out, errc, timeSpan*time.Second)
	go output.Forward(ctx, out, errc, u.String())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
