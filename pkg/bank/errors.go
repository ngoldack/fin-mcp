package bank

import (
	"strings"
)

type CleanError struct {
	Title       string
	Description string
}

func (e CleanError) Error() string {
	return e.Title + ": " + e.Description
}

func FriendlyError(err error) CleanError {
	if err == nil {
		return CleanError{}
	}

	errStr := err.Error()

	switch {
	case strings.Contains(errStr, "dial tcp") || strings.Contains(errStr, "no such host") || strings.Contains(errStr, "timeout") || strings.Contains(errStr, "connection refused"):
		return CleanError{
			Title:       "Network Connection Error",
			Description: "Unable to connect to the bank API. Please check your internet connection and try again.",
		}
	case strings.Contains(errStr, "401") || strings.Contains(errStr, "Unauthorized") || strings.Contains(errStr, "invalid_client"):
		return CleanError{
			Title:       "Authentication Error",
			Description: "Your Application ID or Private Key is incorrect. Double-check your Enable Banking Dashboard settings.",
		}
	case strings.Contains(errStr, "expired") || strings.Contains(errStr, "consent") || strings.Contains(errStr, "Consent"):
		return CleanError{
			Title:       "Bank Consent Expired",
			Description: "Your 90-day bank authorization has expired. Please run the setup wizard to re-authorize.",
		}
	case strings.Contains(errStr, "no active bank session") || strings.Contains(errStr, "session_id"):
		return CleanError{
			Title:       "Bank Session Required",
			Description: "No active bank session found. Please run the setup wizard first to link your bank account.",
		}
	case strings.Contains(errStr, "all fields are required"):
		return CleanError{
			Title:       "Missing Information",
			Description: "All fields are required to initiate the transfer. Please fill out IBAN, Name, and Amount.",
		}
	case strings.Contains(errStr, "invalid amount format"):
		return CleanError{
			Title:       "Invalid Amount",
			Description: "The amount must be a positive decimal number (e.g., 10.50).",
		}
	case strings.Contains(errStr, "REDIRECT_URI_NOT_ALLOWED") || strings.Contains(errStr, "redirect_url") || strings.Contains(errStr, "redirect_uri"):
		return CleanError{
			Title:       "Redirect URL Mismatch",
			Description: "The redirect URL does not match the allowed redirect URIs configured in your Enable Banking Developer Dashboard.",
		}
	default:
		return CleanError{
			Title:       "Unexpected Error",
			Description: errStr,
		}
	}
}
