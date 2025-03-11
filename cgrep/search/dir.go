package search

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

var (
	currentDir string
	gitRegExp  = regexp.MustCompile(`\.git$`)
)

type Dir interface {
	Search(ctx context.Context)
}

type dir struct {
	wg            *sync.WaitGroup
	path          string
	regexp        *regexp.Regexp
	subDirs       []Dir
	fileFullPaths []string
}

// ディレクトリごとに検索用オブジェクトを生成するファクトリ関数
func New(wg *sync.WaitGroup, fullPath string, re *regexp.Regexp) (Dir, error) {
	d := &dir{wg: wg, path: fullPath, regexp: re}
	if d.isGitDri() {
		return d, nil
	}

	err := d.Scan()
	if err != nil {
		return nil, err
	}

	return d, nil
}

// func New() を実行した際、自身のサブディレクトリとファイル郡をスキャンする処理
func (d *dir) Scan() error {
	fs, err := os.ReadDir(d.path)
	if err != nil {
		return err
	}

	for _, f := range fs {
		path := filepath.Join(d.path, f.Name())
		if f.IsDir() {
			subDir, err := New(d.wg, path, d.regexp)
			if err != nil {
				return err
			}
			d.subDirs = append(d.subDirs, subDir)
			continue
		}

		d.fileFullPaths = append(d.fileFullPaths, path)
	}

	return nil
}

// 対象ディレクトリ内のファイルの内容を正規表現で検索し、サブディレクトリに対して再帰的に検索を行うメソッド
func (d *dir) Search(ctx context.Context) {
	// TODO: Implement here
	// 以下の処理を実装する必要があります：
	// 1. サブディレクトリの検索
	//    - 配下のディレクトリに対して再帰的に検索を実行する
	//    - 各ディレクトリの検索は非同期で行う
	//
	// 2. ファイル検索の実行
	//    - 現在のディレクトリ内のファイルに対して検索を実行する
	//    - エラーが発生した場合は errors パッケージを使用してエラーを設定する
	//
	// ヒント：
	// - このメソッド自体は goroutine として実行される前提で設計しています
	// . - すべての非同期処理が完了するまで待ち合わせるためにはどうするか考えてみましょう
	// . - ctx がキャンセルされた場合に非同期処理もキャンセル可能な形で実装する必要があります
	//   - dir.wg の型である sync.WaitGroup の使い方を調べてみましょう
	// - エラーは errors パッケージを使用して処理します
}

// 配下のファイルの内容を読み取り、正規表現に一致するファイルを検索するメソッド
func (d *dir) GrepFiles(ctx context.Context) error {
	for _, path := range d.fileFullPaths {
		if ctx.Err() != nil {
			return nil
		}

		if err := func(path string) error {
			// TODO: Implement here
			// 以下の処理を実装する必要があります：
			// 1. ファイルの準備
			//    - 指定されたパスのファイルを開く
			//
			// 2. ファイル内容の検索
			//    - ファイルの内容を読み込み、正規表現パターンにマッチする行を検出する
			//    - 検索結果を適切に保存する
			//
			// ヒント：
			// - ファイルを Open した後は必ず Close するようにしましょう
			// . - エラーで return された場合でも必ずです
			// - path は絶対パスで渡されます
			// - 余裕があれば、読み込むファイルのサイズが大きい場合なども考慮できると素晴らしいです
			return nil
		}(path); err != nil {
			return err
		}
	}
	return nil
}

// 自身が .git ディレクトリであるかを検証するメソッド
func (d *dir) isGitDri() bool {
	return gitRegExp.MatchString(d.path)
}

// *os.File を渡すと、ファイル名にカレントディレクトリからそのファイルまでのフルパスを添えて返す関数
func relativePath(file *os.File) (string, error) {
	return filepath.Rel(currentDir, file.Name())
}

// 処理の開始時に実行される関数
func init() {
	currentDir, _ = os.Getwd()
}
