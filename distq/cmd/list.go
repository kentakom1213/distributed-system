package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var (
	listServerURL string
	listStatus    string
)

type listJobResponse struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Staus       string          `json:"status"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	LeaseID     *string         `json:"lease_id,omitempty"`
	LeasedBy    *string         `json:"leased_by,omitempty"`
	LeaseUntil  *time.Time      `json:"lease_until,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	LastError   *string         `json:"last_error,omitempty"`
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List jobs",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, err := url.Parse(listServerURL + "/jobs")
		if err != nil {
			return err
		}

		q := u.Query()
		if listStatus != "" {
			q.Set("status", listStatus)
		}
		u.RawQuery = q.Encode()

		req, err := http.NewRequestWithContext(
			cmd.Context(),
			http.MethodGet,
			u.String(),
			nil,
		)
		if err != nil {
			return err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("list failed: status=%s body=%s", resp.Status, string(body))
		}

		var jobs []listJobResponse
		if err := json.Unmarshal(body, &jobs); err != nil {
			return err
		}

		printJobs(jobs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listServerURL, "server", "http://localhost:8080", "queue server URL")
	listCmd.Flags().StringVar(&listStatus, "status", "", "filter by status")
}

func printJobs(jobs []listJobResponse) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "ID\tSTATUS\tATTEMPTS\tTYPE\tLEASED_BY\tCREATED_AT")

	for _, job := range jobs {
		leasedBy := "-"
		if job.LeasedBy != nil {
			leasedBy = *job.LeasedBy
		}

		fmt.Fprintf(
			w,
			"%s\t%s\t%d/%d\t%s\t%s\t%s\n",
			job.ID,
			job.Staus,
			job.Attempts,
			job.MaxAttempts,
			job.Type,
			leasedBy,
			job.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	_ = w.Flush()
}
