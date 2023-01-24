/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cgrep/errors"
	"cgrep/result"
	"cgrep/search"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/spf13/cobra"
)

var dir string
var withContent bool

var rootCmd = &cobra.Command{
	Use:   "cgrep [flags] [args]",
	Short: "Search for file names containing a argument",
	Long: `Search file names contains argument.
Arguments are treated as regular expressions.

Args:
  A search string that can be compiled as a regular expression`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fullPath, err := filepath.Abs(dir)
		if err != nil {
			return err
		}

		if err := ExecSearch(fullPath, args[0]); err != nil {
			return err
		}

		if err := errors.Error(); err != nil {
			return err
		}

		Render(os.Stdout)
		return nil
	},
}

// TODO: 検索処理を非同期で実行する関数
// TODO: sync.WaitGroup、検索ルート、正規表現オブジェクトを search.New() に渡して検索オブジェクトを作成する
// TODO: 検索オブジェクト生成後に非同期で Dir.Search() を実行する
// TODO: すべての検索処理が終わるまで処理をブロックして完了を待つ
// TODO: エラー発生時は即時リターンする
func ExecSearch(fullPath, regexpWord string) error {
	re, err := regexp.Compile(regexpWord)
	if err != nil {
		return err
	}

	wg := new(sync.WaitGroup)

	s, err := search.New(wg, fullPath, re)
	if err != nil {
		return err
	}

	wg.Add(1)
	go s.Search()
	wg.Wait()

	return nil
}

// TODO: 検索結果を標準出力に出力する
// TODO: 標準出力は引数 w io.Writer として渡される想定
// TODO: グローバル変数 withContent が false の場合はファイル名のみ、 true の場合は内容も出力する
func Render(w io.Writer) {
	if withContent {
		result.RenderWithContent(w)
	} else {
		result.RenderFiles(w)
	}
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
