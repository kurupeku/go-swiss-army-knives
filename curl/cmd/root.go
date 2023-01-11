/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"curl/client"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var method, data string
var customHeaders = make([]string, 0)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "curl [URL]",
	Short: "curl is http/https client command.",
	Long: `curl is http/https client command.
- Args: URL
- Available HTTP Methods: GET, POST, PUT, DELETE, PATCH
- Available Content-Type: application/json(only for POST, PUT, PATCH)`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("%s: You must set only URL", err.Error())
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.ValidateFlags(args[0], method, data, customHeaders); err != nil {
			return err
		}

		c, err := client.NewHttpClient(args[0], method, data, customHeaders)
		if err != nil {
			return err
		}

		req, res, err := c.Execute()
		if err != nil {
			return err
		}

		fmt.Println(req)
		fmt.Println(res)

		return nil
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
	rootCmd.Flags().StringVarP(&data, "data", "d", "", "HTTP Post, Put, Patch Data")
	rootCmd.Flags().StringArrayVarP(&customHeaders, "header", "H", []string{}, "Pass custom header(s) to server")
}
