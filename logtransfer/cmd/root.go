/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
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

		var (
			ctx, stop     = signal.NotifyContext(context.Background(), os.Interrupt)
			ln, out, errc = make(chan []byte, 1), make(chan []byte, 1), make(chan error, 1)
		)
		defer stop()

		go input.Monitor(ctx, ln, errc, stdout)
		go storage.Listen(ctx, ln, errc)
		go storage.Load(ctx, out, errc, 5*time.Second)
		go output.Forward(ctx, out, errc, u.String())
		go logs.Error(ctx, errc)

		subCmd.Run()
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
