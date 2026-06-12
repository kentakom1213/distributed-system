package cmd

import (
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("list is not implemented yet")
		return nil
	},
}

var listServerURL string

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listServerURL, "server", "http://localhost:8080", "queue server URL")
}
