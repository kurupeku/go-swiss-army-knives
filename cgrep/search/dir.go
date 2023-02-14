package search

import (
	"bufio"
	"cgrep/errors"
	"cgrep/result"
	"fmt"
	"io/ioutil"
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
	Search()
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

	d.Scan()
	return d, nil
}

// func New() を実行した際、自身のサブディレクトリとファイル郡をスキャンする処理
func (d *dir) Scan() error {
	fs, err := ioutil.ReadDir(d.path)
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

// TODO: サブディレクトリの検索を非同期で行う
func (d *dir) Search() {
	// TODO: 1 週目：配下のディレクトリ・ファイル検索機能の実装
	// レシーバー関数 dには、初期値(カレントディレクトリのスキャン結果)が入っている
	// 自身も非同期で実行される想定なので Done状態の時のみ実行可能とする
	defer d.wg.Done()

	// サブディレクトリの並列処理
	for _, s := range d.subDirs {

		// 非同期処理の開始を d.wg に知らせるようにする
		d.wg.Add(1)
		fmt.Println(s)
		go s.Search()
	}

	// 配下のファイル郡の内容一致検索用メソッド d.GrepFiles() を実行する
	err := d.GrepFiles()
	if err != nil {
		// エラーが発生したら errors.Set(err error) に投げる
		errors.Set(err)
	}

}

func (d *dir) GrepFiles() error {
	// TODO: 1 週目：配下のディレクトリ・ファイル検索機能の実装

	// 配下のファイルは d.fileFullPaths にフルパスの []string として保存されている
	for _, path := range d.fileFullPaths {
		println("ファイルパスは ", path)

		// 配下のファイルの内容を読み取り、正規表現に一致するファイルを検索する
		file, err := os.Open(path) // For read access.
		if err != nil {
			// エラーが発生したら即時リターンする
			return err
		}
		defer func() {
			file.Close()
		}()

		scanner := bufio.NewScanner(file)
		var lineCount int = 0
		for scanner.Scan() {
			lineCount++
			line := scanner.Text()

			// d.regexp に一致させたい正規表現が保存されているのでファイル内の文字列が一致するか検証する
			if d.regexp.MatchString(line) {
				// ファイル名は検索ルートからの相対パスを添えて保存する
				relPath, err := relativePath(file)
				if err != nil {
					return err
				}
				// 一致した場合はファイル名、一致した行の内容、行番号を result.Set() に渡して保存する
				// result.Set(fileName, txt string, no int)
				result.Set(relPath, line, lineCount)
			}
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
