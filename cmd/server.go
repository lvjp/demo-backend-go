package cmd

import (
	"go.lvjp.me/demo-backend-go/internal/app/appcontext"
	"go.lvjp.me/demo-backend-go/internal/app/cmd/server"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the embedded web server",
	Long:  "Serve the API on an integrated webserver",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, err := appcontext.NewFromCommand(cmd)
		if err != nil {
			return err
		}

		return server.Run(ctx)
	},
}
