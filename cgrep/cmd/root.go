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

	"cgrep/errors"

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
	// 以下の処理を実装する必要があります：
	// 1. 検索文字列のバリデーション
	//    - 引数として渡された文字列を正規表現としてコンパイルする
	//    - 不正な正規表現の場合はエラーを返す
	//
	// 2. 非同期検索の準備
	//    - 複数の検索処理を並行して実行できるよう、同期の仕組みを用意する
	//    - 全ての検索処理の完了を待ち合わせられるようにする
	//
	// 3. 検索の実行
	//    - 指定されたディレクトリに対して検索処理を開始する
	//    - 検索処理は非同期で実行し、完了を待ち合わせる
	//    - ctx がキャンセルされた場合に非同期処理もキャンセル可能な形で実装する
	// ヒント：
	// - 非同期処理は終了を待たないと期待通りに動作しません
	// - キャンセル処理には context.Context を使用します
	//   - context.Context は非常によく使われるので、これを期に理解を深めましょう
	// - 正規表現のコンパイルには regexp パッケージを使用します
	return nil
}

// 検索結果を出力する関数
func Render(w io.Writer) {
	// TODO: Implement here
	// 以下の処理を実装する必要があります：
	// 1. フラグに基づく出力内容の決定
	//    - フラグの状態を確認し、ユーザーが要求した出力形式を判断する
	//
	// 2. 適切な形式での出力
	//    - マッチしたファイルの一覧のみを表示するか
	//    - マッチした行の内容も含めて表示するか
	//
	// ヒント：
	// - フラグがどこに格納されているかは CLI ライブラリのドキュメントを参照しましょう
	//   - https://github.com/spf13/cobra
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
