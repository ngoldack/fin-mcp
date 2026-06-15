package bank

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcvalkey "github.com/testcontainers/testcontainers-go/modules/valkey"
)

// TestValkeyCache_Integration spins a real Valkey via testcontainers and
// exercises the valkey backend end to end. It skips when Docker is unavailable.
func TestValkeyCache_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping valkey integration test in -short mode")
	}
	ctx := context.Background()

	container, err := tcvalkey.Run(ctx, "valkey/valkey:7.2.5")
	if err != nil {
		t.Skipf("could not start valkey container (Docker unavailable?): %v", err)
	}
	testcontainers.CleanupContainer(t, container)

	addr, err := container.Endpoint(ctx, "")
	if err != nil {
		t.Fatalf("container endpoint: %v", err)
	}

	c, err := NewCache(CacheOptions{
		Type:   "valkey",
		TTL:    time.Minute,
		Valkey: ValkeyOptions{Address: addr},
	})
	if err != nil {
		t.Fatalf("NewCache: %v", err)
	}
	defer func() { _ = c.Close() }()

	accounts := []Account{{ID: "v1", Name: "Valkey", IBAN: "DE89370400440532013000", Currency: "EUR"}}
	c.SetAccounts(ctx, accounts)
	got, ok := c.GetAccounts(ctx)
	if !ok || len(got) != 1 || got[0].ID != "v1" {
		t.Fatalf("accounts miss: ok=%v %+v", ok, got)
	}

	c.SetDetail(ctx, "v1", AccountDetail{Account: Account{ID: "v1"}})
	if _, ok := c.GetDetail(ctx, "v1"); !ok {
		t.Fatal("detail miss")
	}

	c.Clear(ctx)
	if _, ok := c.GetAccounts(ctx); ok {
		t.Error("accounts should be gone after Clear")
	}
}
