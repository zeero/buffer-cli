package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	apiToken       string
	jsonOutput     bool
	orgFlag        string
	bufferEndpoint = "https://api.buffer.com"
)

var rootCmd = &cobra.Command{
	Use:   "buffer",
	Short: "Buffer API CLI",
	Long:  "Buffer APIをターミナルから操作するCLIツールです。",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiToken, "token", os.Getenv("BUFFER_TOKEN"), "Buffer API Token")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "JSON形式で出力する（AIエージェント向け）")
	rootCmd.PersistentFlags().StringVar(&orgFlag, "org", "", "Organization ID または名前（省略時は最初の1件）")
}
