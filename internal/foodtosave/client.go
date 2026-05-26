package foodtosave

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.foodtosave.com.br"

type Client struct {
	httpClient *http.Client
}

// NewClient creates a client with a default timeout of 15 seconds.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// GetGondolas retrieves all available gondolas for the merchant.
func (c *Client) GetGondolas(ctx context.Context, merchantID string) ([]Gondola, error) {
	reqURL := fmt.Sprintf("%s/api/v1/merchants/%s/gondolas", baseURL, merchantID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, truncateBytes(body, 512))
	}
	var gondolas []Gondola
	if err := json.Unmarshal(body, &gondolas); err != nil {
		return nil, fmt.Errorf("decode JSON: %w", err)
	}
	return gondolas, nil
}

func truncateBytes(b []byte, max int) string {
	if len(b) <= max {
		return string(b)
	}
	return string(b[:max]) + "…"
}
