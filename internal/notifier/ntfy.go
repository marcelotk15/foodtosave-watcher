package notifier

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Notifier sends text messages to an ntfy topic.
type Notifier struct {
	serverURL string
	topic     string
	client    *http.Client
}

func New(serverURL, topic string) *Notifier {
	return &Notifier{
		serverURL: strings.TrimRight(serverURL, "/"),
		topic:     topic,
		client:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (n *Notifier) Send(ctx context.Context, title, body string) error {
	reqURL := fmt.Sprintf("%s/%s", n.serverURL, n.topic)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("create ntfy request: %w", err)
	}
	req.Header.Set("Title", title)
	req.Header.Set("Priority", "default")
	req.Header.Set("Tags", "shopping_bags")
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("send ntfy: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ntfy response status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
