package cmd

import (
	"github.com/spf13/cobra"
)

// enqueueCmd represents the enqueue command
var enqueueCmd = &cobra.Command{
	Use:   "enqueue",
	Short: "Enqueue a job",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("enqueue is not implemented yet")
		return nil
	},
}

var (
	enqueueServerURL   string
	enqueueType        string
	enqueuePayload     string
	enqueueMaxAttempts int
)

func init() {
	rootCmd.AddCommand(enqueueCmd)

	enqueueCmd.Flags().StringVar(&enqueueServerURL, "server", "https://localhost:8080", "queue server URL")
	enqueueCmd.Flags().StringVar(&enqueueType, "type", "sleep", "job type")
	enqueueCmd.Flags().StringVar(&enqueuePayload, "payload", `{"seconds": 3}`, "job payload JSON")
	enqueueCmd.Flags().IntVar(&enqueueMaxAttempts, "max-attempts", 3, "max attempts")
}
