package search

import (
	"bufio"
	"context"
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
	defer d.wg.Done()

	// サブディレクトリの検索を非同期で実行
	for _, subDir := range d.subDirs {
		d.wg.Add(1)
		go subDir.Search(ctx)
	}

	// ファイル内容の検索
	if err := d.GrepFiles(); err != nil {
		errors.Set(err)
	}
}

// 配下のファイルの内容を読み取り、正規表現に一致するファイルを検索するメソッド
func (d *dir) GrepFiles() error {
	for _, path := range d.fileFullPaths {
		if err := func(path string) error {
			// TODO: Implement here

			// 対象のファイルを開く
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// 相対パスを取得
			relPath, err := relativePath(file)
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(file)
			scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
			lineNo := 0
			for scanner.Scan() {
				lineNo++
				line := scanner.Text()
				if d.regexp.MatchString(line) {
					result.Set(relPath, line, lineNo)
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}

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
