package watcher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"foodtosave-notify/internal/config"
	"foodtosave-notify/internal/foodtosave"
	"foodtosave-notify/internal/notifier"
)

const (
	titleNewBag       = "Nova sacola disponível"
	titleQtyDecreased = "Sacolas diminuíram"
	titleSoldOut      = "Sacola esgotada"
)

// Watcher orchestrates API query, diff, and ntfy sending.
type Watcher struct {
	cfg    *config.Config
	api    *foodtosave.Client
	notify *notifier.Notifier
	cache  *Cache
	log    *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Watcher {
	if log == nil {
		log = slog.Default()
	}
	return &Watcher{
		cfg:    cfg,
		api:    foodtosave.NewClient(),
		notify: notifier.New(cfg.Ntfy.ServerURL, cfg.Ntfy.Topic),
		cache:  NewCache(cfg.CachePath),
		log:    log,
	}
}

func (w *Watcher) Tick(ctx context.Context) {
	w.log.Info("starting monitoring routine")
	for _, m := range w.cfg.Merchants {
		mid, mname := m.ID, m.Name
		w.log.Info("querying merchant", "merchant_id", mid, "merchant_name", mname)
		list, err := w.api.GetGondolas(ctx, mid)
		if err != nil {
			w.log.Error("Food To Save API error", "merchant_id", mid, "err", err)
			continue
		}
		w.log.Info("gondolas retrieved", "merchant_id", mid, "count", len(list))
		cur := snapshotFromGondolas(list)
		prev := w.cache.Snapshot(mid)
		events := Diff(prev, cur)
		w.sendEvents(ctx, mname, events)
		if err := w.cache.Replace(mid, cur); err != nil {
			w.log.Error("failed to persist cache to disk", "merchant_id", mid, "err", err)
		}
	}
}

func (w *Watcher) sendEvents(ctx context.Context, merchantName string, events []Event) {
	for _, ev := range events {
		title, body := formatNotification(merchantName, ev)
		if err := w.notify.Send(ctx, title, body); err != nil {
			w.log.Error("failed to send ntfy", "event", ev.Type, "gondola_id", ev.GondolaID, "err", err)
			continue
		}
		w.log.Info("notification sent", "event", ev.Type, "title", title, "gondola_id", ev.GondolaID)
	}
}

func snapshotFromGondolas(list []foodtosave.Gondola) map[string]CachedGondola {
	out := make(map[string]CachedGondola, len(list))
	for _, g := range list {
		id := g.ID
		if id == "" {
			continue
		}
		out[id] = CachedGondola{
			GondolaID:         id,
			Quantity:          g.Quantity,
			BagDescription:    g.Bag.Description,
			BagCategory:       g.Bag.Category,
			BagPrice:          g.Bag.Price,
			AvailabilityEndAt: g.AvailabilityEndAt,
		}
	}
	return out
}

func formatNotification(merchantName string, ev Event) (title string, body string) {
	s := ev.Snapshot
	switch ev.Type {
	case EventNewBag:
		title = titleNewBag
		body = fmt.Sprintf(
			`Nova sacola disponível em %s

Sacola: %s
Categoria: %s
Preço: R$ %.2f
Quantidade: %d
Disponível até: %s`,
			merchantName,
			s.BagDescription,
			s.BagCategory,
			s.BagPrice,
			s.Quantity,
			formatAvailEndAt(s.AvailabilityEndAt),
		)
	case EventQuantityDecreased:
		title = titleQtyDecreased
		body = fmt.Sprintf(
			`Sacolas diminuíram em %s

Sacola: %s
Quantidade anterior: %d
Quantidade atual: %d`,
			merchantName,
			s.BagDescription,
			ev.OldQuantity,
			ev.NewQuantity,
		)
	case EventBagSoldOut:
		title = titleSoldOut
		body = fmt.Sprintf(
			`Sacola esgotada em %s

Sacola: %s
Última quantidade vista: %d`,
			merchantName,
			s.BagDescription,
			ev.OldQuantity,
		)
	default:
		title = "Food To Save"
		body = ""
	}
	return title, body
}

func formatAvailEndAt(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format(time.RFC3339)
}
