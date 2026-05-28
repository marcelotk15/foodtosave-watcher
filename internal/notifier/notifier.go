package notifier

import "context"

// Notifier sends push notifications to an external service.
type Notifier interface {
	Send(ctx context.Context, title, body string) error
}
