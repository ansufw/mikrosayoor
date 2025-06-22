package cmd

import (
	"github.com/spf13/cobra"

	"notification-service/internal/app"
)

var startCmd = &cobra.Command{
	Use:   "notif-start",
	Short: "start",
	Long:  "start",
	Run: func(cmd *cobra.Command, args []string) {
		// Call Func Route API
		app.RunServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
