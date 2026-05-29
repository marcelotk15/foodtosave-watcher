package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultGotifyPriority = 5
	markdownContentType   = "text/markdown"
)

type gotifyExtras struct {
	ClientDisplay struct {
		ContentType string `json:"contentType"`
	} `json:"client::display"`
}

type gotifyPayload struct {
	Title    string       `json:"title"`
	Message  string       `json:"message"`
	Priority int          `json:"priority"`
	Extras   gotifyExtras `json:"extras"`
}

func newMarkdownPayload(title, body string) gotifyPayload {
	p := gotifyPayload{
		Title:    title,
		Message:  body,
		Priority: defaultGotifyPriority,
	}
	p.Extras.ClientDisplay.ContentType = markdownContentType
	return p
}

// GotifyNotifier sends messages to a Gotify application.
type GotifyNotifier struct {
	serverURL string
	appToken  string
	client    *http.Client
}

func NewGotify(serverURL, appToken string) *GotifyNotifier {
	return &GotifyNotifier{
		serverURL: strings.TrimRight(serverURL, "/"),
		appToken:  appToken,
		client:    &http.Client{Timeout: 15 * time.Second},
	}
}

func (g *GotifyNotifier) Send(ctx context.Context, title, body string) error {
	payload, err := json.Marshal(newMarkdownPayload(title, body))
	if err != nil {
		return fmt.Errorf("marshal gotify payload: %w", err)
	}

	reqURL := g.serverURL + "/message"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create gotify request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Gotify-Key", g.appToken)

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("send gotify: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("gotify response status %d: token inválido — use o token de uma Application (não o token de Client): %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("gotify response status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
