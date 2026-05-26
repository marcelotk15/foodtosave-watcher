package watcher

import (
	"testing"
	"time"
)

func gondola(id string, qty int) CachedGondola {
	return CachedGondola{
		GondolaID:         id,
		Quantity:          qty,
		BagDescription:    "Test bag",
		BagCategory:       "Misc",
		BagPrice:          10,
		AvailabilityEndAt: time.Date(2026, 5, 26, 13, 0, 0, 0, time.UTC),
	}
}

func TestDiff_NewBag(t *testing.T) {
	prev := map[string]CachedGondola{}
	cur := map[string]CachedGondola{
		"a": gondola("a", 5),
	}
	got := Diff(prev, cur)
	if len(got) != 1 || got[0].Type != EventNewBag || got[0].GondolaID != "a" {
		t.Fatalf("expected NewBag for gondola a; got %#v", got)
	}
}

func TestDiff_QuantityDecreased(t *testing.T) {
	prev := map[string]CachedGondola{"a": gondola("a", 5)}
	cur := map[string]CachedGondola{"a": gondola("a", 3)}
	got := Diff(prev, cur)
	if len(got) != 1 || got[0].Type != EventQuantityDecreased {
		t.Fatalf("expected QuantityDecreased; got %#v", got)
	}
	if got[0].OldQuantity != 5 || got[0].NewQuantity != 3 {
		t.Fatalf("quantities are wrong: %#v", got[0])
	}
}

func TestDiff_QuantityEqual_NoEvent(t *testing.T) {
	prev := map[string]CachedGondola{"a": gondola("a", 5)}
	cur := map[string]CachedGondola{"a": gondola("a", 5)}
	got := Diff(prev, cur)
	if len(got) != 0 {
		t.Fatalf("expected no events; got %#v", got)
	}
}

func TestDiff_QuantityIncreased_NoEvent(t *testing.T) {
	prev := map[string]CachedGondola{"a": gondola("a", 3)}
	cur := map[string]CachedGondola{"a": gondola("a", 7)}
	got := Diff(prev, cur)
	if len(got) != 0 {
		t.Fatalf("expected no events; got %#v", got)
	}
}

func TestDiff_BagSoldOut(t *testing.T) {
	prev := map[string]CachedGondola{"a": gondola("a", 4)}
	cur := map[string]CachedGondola{}
	got := Diff(prev, cur)
	if len(got) != 1 || got[0].Type != EventBagSoldOut {
		t.Fatalf("expected BagSoldOut; got %#v", got)
	}
	if got[0].OldQuantity != 4 {
		t.Fatalf("OldQuantity is wrong: %d", got[0].OldQuantity)
	}
}

func TestDiff_AfterSoldOut_NoRepeat(t *testing.T) {
	// Após esgotar, cache vazio; próxima Diff com API vazia não gera novo sold out.
	prev := map[string]CachedGondola{}
	cur := map[string]CachedGondola{}
	got := Diff(prev, cur)
	if len(got) != 0 {
		t.Fatalf("expected no events; got %#v", got)
	}
}

func TestDiff_SameBagSecondRun_NoNewBag(t *testing.T) {
	prev := map[string]CachedGondola{}
	cur1 := map[string]CachedGondola{"a": gondola("a", 5)}
	first := Diff(prev, cur1)
	if len(first) != 1 || first[0].Type != EventNewBag {
		t.Fatalf("first execution failed: %#v", first)
	}
	prev2 := cur1
	cur2 := map[string]CachedGondola{"a": gondola("a", 5)}
	second := Diff(prev2, cur2)
	if len(second) != 0 {
		t.Fatalf("second execution should not generate NewBag; got %#v", second)
	}
}
