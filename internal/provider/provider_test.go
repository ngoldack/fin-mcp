package provider

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/ngoldack/fin-mcp/internal/config"
)

func TestFromConfig_MockProvider(t *testing.T) {
	cfg := &config.Config{
		Providers: []config.ProviderConfig{{Type: "mock", Name: "m1"}},
	}
	reg, err := FromConfig(cfg, filepath.Join(t.TempDir(), "c.json"))
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}
	p, ok := reg.Default()
	if !ok {
		t.Fatal("no default provider")
	}
	if p.Name() != "m1" {
		t.Errorf("name = %q, want m1", p.Name())
	}
	accs, err := p.ListAccounts(context.Background())
	if err != nil || len(accs) != 1 {
		t.Errorf("ListAccounts = %v, %v", len(accs), err)
	}
}

func TestFromConfig_MultipleProviders(t *testing.T) {
	cfg := &config.Config{
		Providers: []config.ProviderConfig{
			{Type: "mock", Name: "primary"},
			{Type: "mock", Name: "secondary"},
		},
	}
	reg, err := FromConfig(cfg, filepath.Join(t.TempDir(), "c.json"))
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}
	if got := len(reg.All()); got != 2 {
		t.Fatalf("registry size = %d, want 2", got)
	}
	if d, _ := reg.Default(); d.Name() != "primary" {
		t.Errorf("default = %q, want primary", d.Name())
	}
	if _, ok := reg.Get("secondary"); !ok {
		t.Error("secondary not registered")
	}
}

func TestFromConfig_Empty(t *testing.T) {
	if _, err := FromConfig(&config.Config{}, "x"); err == nil {
		t.Error("expected error for no providers")
	}
}

func TestFromConfig_UnknownType(t *testing.T) {
	cfg := &config.Config{Providers: []config.ProviderConfig{{Type: "wells-fargo"}}}
	if _, err := FromConfig(cfg, "x"); err == nil {
		t.Error("expected error for unknown provider type")
	}
}
