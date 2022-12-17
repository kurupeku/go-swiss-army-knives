/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"curl/client"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var (
	method              string
	withQueryParamsFlag bool
)
var data = make([]string, 0)
var customHeaders = make([]string, 0)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "curl",
	Short: "curl is http/https client command.",
	Long: `curl is http/https client command.
- Available HTTP Methods: GET, POST, PUT, DELETE, PATCH
- Available Content-Type: application/json, application/x-www-form-urlencoded`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("this command requires at least one argument")
		}

		builder := client.NewHttpClientBuilder(args[0], method, data, withQueryParamsFlag, customHeaders)

		if err := builder.Validate(); err != nil {
			return err
		}

		c, err := builder.Build()
		if err != nil {
			return err
		}

		return c.Execute()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&method, "request", "X", "GET", "HTTP method")
	rootCmd.Flags().StringArrayVarP(&data, "data", "d", []string{}, "HTTP Post, Put, Patch Data")
	rootCmd.Flags().BoolVarP(&withQueryParamsFlag, "get", "G", false, "Put the post data in the URL and use GET")
	rootCmd.Flags().StringArrayVarP(&customHeaders, "header", "H", []string{}, "Pass custom header(s) to server")
}
