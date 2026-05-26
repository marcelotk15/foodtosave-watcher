package watcher

import (
	"sort"
)

type EventType string

const (
	EventNewBag            EventType = "new_bag"
	EventQuantityDecreased EventType = "quantity_decreased"
	EventBagSoldOut        EventType = "bag_sold_out"
)

type Event struct {
	Type        EventType
	GondolaID   string
	OldQuantity int
	NewQuantity int
	Snapshot    CachedGondola
}

func Diff(previous, current map[string]CachedGondola) []Event {
	ids := unionKeys(previous, current)
	sort.Strings(ids)

	var out []Event
	for _, id := range ids {
		prev, okPrev := previous[id]
		cur, okCur := current[id]

		switch {
		case !okPrev && okCur:
			out = append(out, Event{
				Type:        EventNewBag,
				GondolaID:   id,
				OldQuantity: 0,
				NewQuantity: cur.Quantity,
				Snapshot:    cur,
			})
		case okPrev && !okCur:
			out = append(out, Event{
				Type:        EventBagSoldOut,
				GondolaID:   id,
				OldQuantity: prev.Quantity,
				NewQuantity: 0,
				Snapshot:    prev,
			})
		case okPrev && okCur && cur.Quantity < prev.Quantity:
			ev := CachedGondola{
				GondolaID:         cur.GondolaID,
				Quantity:          cur.Quantity,
				BagDescription:    cur.BagDescription,
				BagCategory:       cur.BagCategory,
				BagPrice:          cur.BagPrice,
				AvailabilityEndAt: cur.AvailabilityEndAt,
			}
			out = append(out, Event{
				Type:        EventQuantityDecreased,
				GondolaID:   id,
				OldQuantity: prev.Quantity,
				NewQuantity: cur.Quantity,
				Snapshot:    ev,
			})
		default:
			// quantity equal or increase — no event
		}
	}
	return out
}

func unionKeys(a, b map[string]CachedGondola) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	ids := make([]string, 0, len(seen))
	for k := range seen {
		ids = append(ids, k)
	}
	return ids
}
