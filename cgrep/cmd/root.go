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
// TODO: 検索オブジェクト生成後に非同期で Dir.Search() を実行する
// TODO: すべての検索処理が終わるまで処理をブロックして完了を待つ
func ExecSearch(fullPath, regexpWord string) error {
	// TODO: 2 週目：検索結果のレンダリング & コマンド実行時のメイン処理の実装

	re, err := regexp.Compile(regexpWord)
	// エラー発生時は即時リターンする
	if err != nil {
		return err
	}
	// newは指定した型のポインタ型を生成する関数
	wg := new(sync.WaitGroup)

	// sync.WaitGroup、検索ルート、正規表現オブジェクトを search.New() に渡して検索オブジェクトを作成する
	dir, err := search.New(wg, fullPath, re)
	if err != nil {
		return err
	}

	// goroutineの作法。非同期を始める前に WaitGroupを +1。defer Doneはgoroutine内部で
	wg.Add(1)
	go dir.Search()

	// メインのgoroutineはサブgoroutine が全て完了するのを待つ
	wg.Wait()

	return nil
}

func Render(w io.Writer) {
	// TODO: 2 週目：検索結果のレンダリング & コマンド実行時のメイン処理の実装
	// グローバル変数 withContent が false の場合はファイル名のみ、 true の場合は内容も出力する
	if withContent {
		//標準出力は引数 w io.Writer として渡される想定
		result.RenderWithContent(w)
	} else {
		result.RenderFiles(w)
	}

	// 検索結果を標準出力に出力する
	// ⇒　別途最終的に Render(os.Stdout)として別途出力される
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
