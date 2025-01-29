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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
