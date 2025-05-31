package main

import (
	"fmt"
	"os"

	"github.com/tsuna-can/s3-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
