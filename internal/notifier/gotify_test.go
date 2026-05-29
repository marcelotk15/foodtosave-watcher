package notifier

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGotifyNotifierSend(t *testing.T) {
	var got gotifyPayload
	var contentType string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/message" {
			t.Errorf("path = %q, want /message", r.URL.Path)
		}
		if gotToken := r.Header.Get("X-Gotify-Key"); gotToken != "test-token" {
			t.Errorf("X-Gotify-Key = %q, want test-token", gotToken)
		}

		contentType = r.Header.Get("Content-Type")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("unmarshal payload: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewGotify(srv.URL, "test-token")
	if err := n.Send(t.Context(), "Titulo", "Corpo da mensagem"); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}
	if got.Title != "Titulo" {
		t.Errorf("title = %q, want Titulo", got.Title)
	}
	if got.Message != "Corpo da mensagem" {
		t.Errorf("message = %q, want Corpo da mensagem", got.Message)
	}
	if got.Priority != defaultGotifyPriority {
		t.Errorf("priority = %d, want %d", got.Priority, defaultGotifyPriority)
	}
	if got.Extras.ClientDisplay.ContentType != markdownContentType {
		t.Errorf("extras.client::display.contentType = %q, want %q", got.Extras.ClientDisplay.ContentType, markdownContentType)
	}
}

func TestGotifyNotifierSendErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer srv.Close()

	n := NewGotify(srv.URL, "test-token")
	err := n.Send(t.Context(), "Titulo", "Corpo")
	if err == nil {
		t.Fatal("Send() expected error, got nil")
	}
}
