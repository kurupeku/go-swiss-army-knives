package search

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"cgrep/errors"
	"cgrep/result"
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
// TODO: 非同期処理の開始を d.wg に知らせるようにする
// TODO: 自身も非同期で実行される想定なので d.wg に処理完了を知らせる
// TODO: 配下のファイル郡の内容一致検索用メソッド d.GrepFiles() を実行する
// TODO: エラーが発生したら errors.Set(err error) に投げる
func (d *dir) Search() {
	defer d.wg.Done()

	for _, subDir := range d.subDirs {
		d.wg.Add(1)
		go subDir.Search()
	}

	if err := d.GrepFiles(); err != nil {
		errors.Set(err)
	}
}

// TODO: 配下のファイルの内容を読み取り、正規表現に一致するファイルを検索する
// TODO: 配下のファイルは d.fileFullPaths にフルパスの []string として保存されている
// TODO: d.regexp に一致させたい正規表現が保存されているのでファイル内の文字列が一致するか検証する
// TODO: 一致した場合はファイル名、一致した行の内容、行番号を result.Set() に渡して保存する
// TODO: ファイル名は検索ルートからの相対パスを添えて保存する
// TODO: エラーが発生したら即時リターンする
func (d *dir) GrepFiles() error {
	for _, path := range d.fileFullPaths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			f.Close()
		}(file)

		var lineNo int
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lineNo++
			str := scanner.Text()
			if !d.regexp.MatchString(str) {
				continue
			}

			rel, err := relativePath(file)
			if err != nil {
				return err
			}

			result.Set(rel, str, lineNo)
		}
		if err := scanner.Err(); err != nil {
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
