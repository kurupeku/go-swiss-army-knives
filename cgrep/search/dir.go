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

func New(wg *sync.WaitGroup, fullPath string, re *regexp.Regexp) (Dir, error) {
	d := &dir{wg: wg, path: fullPath, regexp: re}
	if d.isGitDri() {
		return d, nil
	}

	d.Scan()
	return d, nil
}

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

func (d *dir) GrepFiles() error {
	for _, path := range d.fileFullPaths {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			f.Close()
		}()

		var lineNo int
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lineNo++
			str := scanner.Text()
			if !d.regexp.MatchString(str) {
				continue
			}

			rel, err := relativePath(f)
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

func (d *dir) isGitDri() bool {
	return gitRegExp.MatchString(d.path)
}

func relativePath(file *os.File) (string, error) {
	return filepath.Rel(currentDir, file.Name())
}

func init() {
	currentDir, _ = os.Getwd()
}
