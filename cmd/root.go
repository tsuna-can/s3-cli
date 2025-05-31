package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tsuna-can/s3-cli/internal/ui"
)

var outputDir string
var profile string
var debugMode bool
var endpointURL string

var rootCmd = &cobra.Command{
	Use:   "s3-cli",
	Short: "Interactive AWS S3 CLI tool",
	Long:  `An interactive CLI tool for browsing and downloading files from AWS S3 buckets.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// --endpoint-urlフラグが指定されているか確認
		if endpointURL == "" {
			return fmt.Errorf("--endpoint-url フラグは必須です")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ui.StartUI(outputDir, profile, endpointURL, debugMode) // 引数にendpointURLを追加
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&outputDir, "output-dir", "", "Directory to save downloaded files (default is current directory)")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "AWS profile to use (default: default)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug mode")

	// エンドポイントURLフラグを追加（必須）
	rootCmd.PersistentFlags().StringVar(&endpointURL, "endpoint-url", "", "AWS S3 endpoint URL (required)")
	rootCmd.MarkPersistentFlagRequired("endpoint-url")
}
