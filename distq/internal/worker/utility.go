package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) postJSON(ctx context.Context, path string, body []byte) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.serverURL+path,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed: status=%s body=%s", resp.Status, string(respBody))
	}

	return nil
}
