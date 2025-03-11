/*
Copyright © 2023 kurupeku <22340645+kurupeku@users.noreply.github.com>
*/
package cmd

import (
	"context"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sync"

	"cgrep/errors"
	"cgrep/result"
	"cgrep/search"

	"github.com/spf13/cobra"
)

var (
	dir         string
	withContent bool
)

var rootCmd = &cobra.Command{
	Use:   "cgrep [flags] [args]",
	Short: "Search for file names containing a argument",
	Long: `Search file names contains argument.
Arguments are treated as regular expressions.

Args:
  A search string that can be compiled as a regular expression`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		fullPath, err := filepath.Abs(dir)
		if err != nil {
			return err
		}

		if err := ExecSearch(ctx, fullPath, args[0]); err != nil {
			return err
		}

		if err := errors.Error(); err != nil {
			return err
		}

		if ctx.Err() != nil {
			return nil
		}

		Render(os.Stdout)
		return nil
	},
}

// 検索処理を非同期で実行する関数
func ExecSearch(ctx context.Context, fullPath, regexpWord string) error {
	// TODO: Implement here

	re := regexp.MustCompile(regexpWord)
	var wg sync.WaitGroup
	wg.Add(1)

	dir, err := search.New(&wg, fullPath, re)
	if err != nil {
		return err
	}

	go dir.Search(ctx)
	wg.Wait()

	return nil
}

// 検索結果を出力する関数
func Render(w io.Writer) {
	// TODO: Implement here

	if withContent {
		result.RenderWithContent(w)
		return
	}
	result.RenderFiles(w)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&dir, "dir", "d", "./", "searching directory")
	rootCmd.Flags().BoolVarP(&withContent, "with-content", "c", false, "render with matched content lines")
}
