package cmd

import (
	"fmt"
	"os"

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
	Run: func(*cobra.Command, []string) {
		if err := server.Run(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}
