/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"logtransfer/input"
	"logtransfer/output"
	"logtransfer/storage"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

const errorFilePath = "error.log"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "logtransfer",
	Short: "A log transfer application over HTTP#POST",
	Long: `A log transfer application over HTTP#POST.
The application consists of a distributed system
with multi-threaded safe transfers.

e.g ) logtransfer https://sample.com`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("request path is required")
		}

		u, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		var (
			ctx, stop     = signal.NotifyContext(context.Background(), os.Interrupt)
			ln, out, errc = make(chan []byte, 1), make(chan []byte, 1), make(chan error, 1)
			ef            *os.File
		)
		defer stop()

		go input.Monitor(ctx, ln, errc, os.Stdout)
		go storage.Listen(ctx, ln, errc)
		go storage.Load(ctx, out, errc, 5*time.Second)
		go output.Forward(ctx, out, errc, u.String())

		ef, err = os.Create(errorFilePath)
		if err != nil {
			return err
		}
		defer ef.Close()

		for {
			select {
			case <-ctx.Done():
				return nil
			case err := <-errc:
				ef.WriteString(err.Error())
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.logtransfer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
