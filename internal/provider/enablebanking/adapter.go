// Package enablebanking adapts the Enable Banking SDK to the provider.Provider
// port. It owns session handling and the mapping from raw SDK models to the
// application's domain types.
package enablebanking

import (
	"context"
	"fmt"
	"time"

	"github.com/ngoldack/fin-mcp/internal/bank"
	"github.com/ngoldack/fin-mcp/internal/config"
	"github.com/ngoldack/fin-mcp/internal/secret"
	eb "github.com/ngoldack/fin-mcp/pkg/enablebanking"
)

// Adapter implements provider.Provider on top of the Enable Banking SDK.
type Adapter struct {
	name    string
	client  eb.APIClient
	cfg     *config.EnableBankingConfig
	persist func() // saves the owning application config (e.g. refreshed consent)
}

// New builds the adapter and its underlying SDK client from a provider config.
// If PrivateKeyKeyring is set, the PEM is read from the OS keychain (local only).
func New(name string, cfg *config.EnableBankingConfig, persist func()) (*Adapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("enable-banking provider %q: missing enable_banking config", name)
	}

	keyContent := cfg.PrivateKeyContent
	if cfg.PrivateKeyKeyring != "" {
		v, err := secret.Get(cfg.PrivateKeyKeyring)
		if err != nil {
			return nil, fmt.Errorf("read private key from keychain account %q: %w", cfg.PrivateKeyKeyring, err)
		}
		keyContent = v
	}

	client := eb.NewClient(cfg.AppID, cfg.PrivateKeyPath, keyContent, cfg.Environment)
	return NewWithClient(name, client, cfg, persist), nil
}

// NewWithClient injects an SDK client (used in tests).
func NewWithClient(name string, client eb.APIClient, cfg *config.EnableBankingConfig, persist func()) *Adapter {
	if persist == nil {
		persist = func() {}
	}
	return &Adapter{name: name, client: client, cfg: cfg, persist: persist}
}

func (a *Adapter) Name() string { return a.name }

func (a *Adapter) Info() bank.ProviderInfo {
	return bank.ProviderInfo{
		Name:              a.name,
		Environment:       a.cfg.Environment,
		BankName:          a.cfg.BankName,
		BankCountry:       a.cfg.BankCountry,
		SessionRef:        a.cfg.SessionID,
		ConsentValidUntil: a.cfg.ConsentValidUntil,
	}
}

func (a *Adapter) VerifyConnection(ctx context.Context) (bank.ConnectionStatus, error) {
	if a.cfg.SessionID == "" {
		return bank.ConnectionStatus{}, fmt.Errorf("no active bank session found; run setup to link your bank account")
	}

	sess, err := a.client.GetSession(ctx, a.cfg.SessionID)
	if err != nil {
		return bank.ConnectionStatus{}, fmt.Errorf("failed to verify bank session: %w; your session may have been invalidated", err)
	}

	status := bank.ConnectionStatus{
		Authorized:        sess.Status == "AUTHORIZED",
		Status:            sess.Status,
		ConsentValidUntil: sess.Access.ValidUntil,
	}

	if !sess.Access.ValidUntil.IsZero() && !sess.Access.ValidUntil.Equal(a.cfg.ConsentValidUntil) {
		a.cfg.ConsentValidUntil = sess.Access.ValidUntil
		a.persist()
	}

	if !status.Authorized {
		return status, fmt.Errorf("bank session status is %s, expected AUTHORIZED; re-run setup to refresh", sess.Status)
	}
	return status, nil
}

func (a *Adapter) ListAccounts(ctx context.Context) ([]bank.Account, error) {
	sess, err := a.client.GetSession(ctx, a.cfg.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve session details: %w", err)
	}

	var accounts []bank.Account
	for _, accID := range sess.Accounts {
		details, err := a.client.GetAccountDetails(ctx, accID)
		if err != nil {
			continue
		}
		accounts = append(accounts, mapAccount(*details, a.cfg.BankName))
	}
	if len(accounts) == 0 {
		return nil, fmt.Errorf("no accounts linked or accessible in this session")
	}
	return accounts, nil
}

func (a *Adapter) GetBalances(ctx context.Context, accountID string) (bank.Balances, error) {
	raw, err := a.client.GetBalances(ctx, accountID)
	if err != nil {
		return bank.Balances{}, err
	}
	items, available, booked := mapBalances(raw)
	return bank.Balances{Items: items, Available: available, Booked: booked}, nil
}

func (a *Adapter) GetTransactions(ctx context.Context, accountID, dateFrom, dateTo string) ([]bank.Transaction, error) {
	raw, err := a.client.GetTransactions(ctx, accountID, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	return mapTransactions(raw), nil
}

func (a *Adapter) InitiateTransfer(ctx context.Context, req bank.TransferRequest) (*bank.TransferResult, error) {
	state := fmt.Sprintf("pay-%d", time.Now().UnixNano())
	resp, err := a.client.CreatePayment(
		ctx, req.DebtorIBAN, req.CreditorIBAN, req.CreditorName,
		req.Amount, req.Currency, req.PaymentType, state, a.cfg.RedirectURL,
	)
	if err != nil {
		return nil, err
	}
	return &bank.TransferResult{PaymentID: resp.PaymentID, Status: resp.Status, AuthURL: resp.URL}, nil
}

func (a *Adapter) PaymentStatus(ctx context.Context, paymentID string) (*bank.TransferResult, error) {
	p, err := a.client.GetPayment(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	return &bank.TransferResult{PaymentID: p.PaymentID, Status: p.Status}, nil
}

func (a *Adapter) SubmitTransfer(ctx context.Context, paymentID string) (*bank.TransferResult, error) {
	r, err := a.client.SubmitPayment(ctx, paymentID)
	if err != nil {
		return nil, err
	}
	res := &bank.TransferResult{PaymentID: paymentID, Status: r.Status}
	if r.Message != "" && res.Status == "" {
		res.Status = r.Message
	}
	return res, nil
}
