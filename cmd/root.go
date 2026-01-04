package cmd

import (
	"fmt"
	"os"

	"go.lvjp.me/demo-backend-go/pkg/buildinfo"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "demo-backend-go",
	Version: buildinfo.Get().VersionString(),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
