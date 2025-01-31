package cmd

import (
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
		os.Exit(1)
	}
}
