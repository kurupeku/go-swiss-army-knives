/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cgrep/result"
	"cgrep/search"
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
		if err := ExecSearch(args[0]); err != nil {
			return err
		}

		if err := result.Error(); err != nil {
			return err
		}

		Render()
		return nil
	},
}

func ExecSearch(regexpWord string) error {
	fullPath, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	re, err := regexp.Compile(regexpWord)
	if err != nil {
		return err
	}

	s, err := search.New(fullPath, re)
	if err != nil {
		return err
	}

	var wg = new(sync.WaitGroup)
	wg.Add(1)
	go s.Search(wg)
	wg.Wait()

	return nil
}

func Render() {
	if withContent {
		result.RenderWithContent(os.Stdout)
	} else {
		result.RenderFiles(os.Stdout)
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
