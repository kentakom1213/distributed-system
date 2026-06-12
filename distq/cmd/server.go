package cmd

import (
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the queue server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("server is not implemented yet")
		return nil
	},
}

var (
	serverAddr string
	dbPath     string
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&serverAddr, "addr", ":8080", "server address")
	serverCmd.Flags().StringVar(&dbPath, "db", "distq.db", "sqlite database path")
}
