package bank

import (
	"context"
	"testing"
	"time"
)

func TestMemoryCache_HitMissExpiry(t *testing.T) {
	ctx := context.Background()
	c := newMemoryCache(50 * time.Millisecond)

	if _, ok := c.GetAccounts(ctx); ok {
		t.Fatal("expected miss on empty cache")
	}

	accounts := []Account{{ID: "a1", Name: "Main", Currency: "EUR"}}
	c.SetAccounts(ctx, accounts)

	got, ok := c.GetAccounts(ctx)
	if !ok || len(got) != 1 || got[0].ID != "a1" {
		t.Fatalf("expected hit with 1 account, got ok=%v %+v", ok, got)
	}

	time.Sleep(70 * time.Millisecond)
	if _, ok := c.GetAccounts(ctx); ok {
		t.Error("expected miss after TTL expiry")
	}
}

func TestNoopCache_AlwaysMisses(t *testing.T) {
	ctx := context.Background()
	c := noopCache{}
	c.SetAccounts(ctx, []Account{{ID: "x"}})
	if _, ok := c.GetAccounts(ctx); ok {
		t.Error("noop cache must never return a hit")
	}
	c.SetDetail(ctx, "x", AccountDetail{})
	if _, ok := c.GetDetail(ctx, "x"); ok {
		t.Error("noop cache must never return a hit")
	}
}

func TestNewCache_Factory(t *testing.T) {
	for _, typ := range []string{"", "memory", "none"} {
		c, err := NewCache(CacheOptions{Type: typ, TTL: time.Minute})
		if err != nil {
			t.Fatalf("type %q: %v", typ, err)
		}
		_ = c.Close()
	}
	if _, err := NewCache(CacheOptions{Type: "bogus"}); err == nil {
		t.Error("unknown cache type should error")
	}
}
