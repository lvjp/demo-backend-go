package cmd

import (
	"os"

	"go.lvjp.me/demo-backend-go/pkg/buildinfo"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "demo-backend-go",
	Version:      buildinfo.Get().VersionString(),
	SilenceUsage: true,
}

func Execute() {
	// Error handling is done inside cobra beacause SilenceErrors is false by default.
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
