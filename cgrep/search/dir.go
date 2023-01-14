package search

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"cgrep/result"
)

var currentDir string

type Dir interface {
	Search(wg *sync.WaitGroup)
}

type dir struct {
	path          string
	regexp        *regexp.Regexp
	subDirs       []Dir
	fileFullPaths []string
}

func New(fullPath string, re *regexp.Regexp) (Dir, error) {
	d := &dir{path: fullPath, regexp: re}
	d.scan()

	return d, nil
}

func (d *dir) scan() error {
	fs, err := ioutil.ReadDir(d.path)
	if err != nil {
		return err
	}

	for _, f := range fs {
		path := filepath.Join(d.path, f.Name())
		if f.IsDir() {
			subDir, err := New(path, d.regexp)
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

func (d *dir) Search(wg *sync.WaitGroup) {
	defer wg.Done()

	for _, subDir := range d.subDirs {
		wg.Add(1)
		go subDir.Search(wg)
	}

	if err := d.grepFiles(); err != nil {
		result.SetError(err)
	}
}

func (d *dir) grepFiles() error {
	for _, path := range d.fileFullPaths {
		if strings.Contains(path, ".git/") {
			continue
		}

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

			rel, err := relativePath(f.Name())
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

func relativePath(fullPath string) (string, error) {
	return filepath.Rel(currentDir, fullPath)
}

func init() {
	currentDir, _ = os.Getwd()
}
