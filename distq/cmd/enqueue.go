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

func init() {
	rootCmd.AddCommand(enqueueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enqueueCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enqueueCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
