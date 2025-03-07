package cmd

import (
	"go.lvjp.me/demo-backend-go/internal/app/cmd/healthcheck"

	"github.com/spf13/cobra"
)

var healthCheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Check the health of the server",
	Run: func(*cobra.Command, []string) {
		healthcheck.Run()
	},
}

func init() {
	rootCmd.AddCommand(healthCheckCmd)
}
