package secret

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestRoundTrip(t *testing.T) {
	account := fmt.Sprintf("fin-mcp-test-%d", time.Now().UnixNano())

	if err := Set(account, "pem-data"); err != nil {
		// No usable keychain (common on headless CI/Linux runners): skip.
		t.Skipf("OS keychain unavailable: %v", err)
	}
	t.Cleanup(func() { _ = Delete(account) })

	got, err := Get(account)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != "pem-data" {
		t.Errorf("Get = %q, want %q", got, "pem-data")
	}

	if err := Delete(account); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := Get(account); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get after delete = %v, want ErrNotFound", err)
	}
}

func TestEmptyAccount(t *testing.T) {
	if err := Set("", "x"); err == nil {
		t.Error("Set with empty account should error")
	}
	if _, err := Get(""); err == nil {
		t.Error("Get with empty account should error")
	}
}
