package cmd

import (
	"github.com/spf13/cobra"
)

// workerCmd represents the worker command
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start a worker",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("worker is not implemented yet")
		return nil
	},
}

var (
	workerID  string
	serverURL string
)

func init() {
	rootCmd.AddCommand(workerCmd)

	workerCmd.Flags().StringVar(&workerID, "id", "worker-1", "worker id")
	workerCmd.Flags().StringVar(&serverURL, "server", "http://localhost:8080", "queue server URL")
}
