// Package secret stores and retrieves sensitive values in the operating
// system's native keychain (macOS Keychain, Windows Credential Manager, Linux
// Secret Service). It is OPTIONAL and intended for local runs only: in
// Kubernetes and other headless deployments, supply secrets via env vars or
// mounted Secret files instead — the keychain is never touched unless a config
// explicitly references it.
package secret

import (
	"errors"

	"github.com/zalando/go-keyring"
)

// Service is the keychain service name under which fin-mcp secrets are grouped.
const Service = "fin-mcp"

// ErrNotFound is returned when no secret exists for the account.
var ErrNotFound = keyring.ErrNotFound

// Set stores value under account in the OS keychain.
func Set(account, value string) error {
	if account == "" {
		return errors.New("secret: empty account")
	}
	return keyring.Set(Service, account, value)
}

// Get retrieves the secret stored under account.
func Get(account string) (string, error) {
	if account == "" {
		return "", errors.New("secret: empty account")
	}
	return keyring.Get(Service, account)
}

// Delete removes the secret stored under account.
func Delete(account string) error {
	if account == "" {
		return errors.New("secret: empty account")
	}
	return keyring.Delete(Service, account)
}
