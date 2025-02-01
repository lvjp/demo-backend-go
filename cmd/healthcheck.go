package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var healthCheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Check the health of the server",
	Run: func(*cobra.Command, []string) {
		resp, err := http.Get("http://localhost:8080/api/v0/misc/version")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer resp.Body.Close()

		fmt.Println(resp.Status)

		if resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(healthCheckCmd)
}
