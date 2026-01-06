package cmd

import (
	"go.lvjp.me/demo-backend-go/internal/app/appcontext"
	"go.lvjp.me/demo-backend-go/internal/app/cmd/healthcheck"

	"github.com/spf13/cobra"
)

var healthCheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Check the health of the server",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, err := appcontext.NewFromCommand(cmd)
		if err != nil {
			return err
		}

		return healthcheck.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(healthCheckCmd)
}
