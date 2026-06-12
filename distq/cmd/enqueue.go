package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

// enqueueCmd represents the enqueue command
var enqueueCmd = &cobra.Command{
	Use:   "enqueue",
	Short: "Enqueue a job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !json.Valid([]byte(enqueuePayload)) {
			return fmt.Errorf("payload must be valid JSON")
		}

		req := map[string]any{
			"type":         enqueueType,
			"payload":      json.RawMessage(enqueuePayload),
			"max_attempts": enqueueMaxAttempts,
		}

		body, err := json.Marshal(req)
		if err != nil {
			return err
		}

		httpReq, err := http.NewRequestWithContext(
			cmd.Context(),
			http.MethodPost,
			enqueueServerURL+"/jobs",
			bytes.NewReader(body),
		)
		if err != nil {
			return err
		}

		httpReq.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			return err
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("enqueue failed: status=%s body=%s", resp.Status, string(respBody))
		}

		cmd.Println(string(respBody))
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

	enqueueCmd.Flags().StringVar(&enqueueServerURL, "server", "http://localhost:8080", "queue server URL")
	enqueueCmd.Flags().StringVar(&enqueueType, "type", "sleep", "job type")
	enqueueCmd.Flags().StringVar(&enqueuePayload, "payload", `{"seconds": 3}`, "job payload JSON")
	enqueueCmd.Flags().IntVar(&enqueueMaxAttempts, "max-attempts", 3, "max attempts")
}
