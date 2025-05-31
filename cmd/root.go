package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tsuna-can/s3-cli/internal/ui"
)

var outputDir string
var profile string
var debugMode bool

var rootCmd = &cobra.Command{
	Use:   "s3-cli",
	Short: "Interactive AWS S3 CLI tool",
	Long:  `An interactive CLI tool for browsing and downloading files from AWS S3 buckets.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.StartUI(outputDir, profile, debugMode)
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
}
