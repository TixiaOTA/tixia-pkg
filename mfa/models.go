package mfa

import (
	"strings"
	"time"
)

// EnrollTOTPRequest represents a TOTP enrollment for an account.
type EnrollTOTPRequest struct {
	AccountName string // Account name shown in the authenticator app, for example an email address (required)
}

// Validate checks that the request has the fields required to enroll TOTP.
func (r EnrollTOTPRequest) Validate() error {
	if strings.TrimSpace(r.AccountName) == "" {
		return ErrMissingAccountName
	}

	return nil
}

// EnrollTOTPResponse is returned after generating a TOTP secret.
type EnrollTOTPResponse struct {
	Secret    string // Base32 secret to store for the account (required)
	URL       string // otpauth:// URL for authenticator apps
	QRCodePNG []byte // PNG-encoded QR code for the otpauth URL
}

// VerifyTOTPRequest represents a TOTP code verification.
type VerifyTOTPRequest struct {
	UserID string // Stable identifier of the account (required)
	Secret string // Base32 secret stored for the account (required)
	Code   string // Code entered by the user (required)
}

// Validate checks that the request has the fields required to verify TOTP.
func (r VerifyTOTPRequest) Validate() error {
	if strings.TrimSpace(r.UserID) == "" {
		return ErrMissingUserID
	}

	if strings.TrimSpace(r.Secret) == "" {
		return ErrMissingSecret
	}

	if strings.TrimSpace(r.Code) == "" {
		return ErrMissingCode
	}

	return nil
}

// VerifyTOTPResponse is returned after a successful TOTP verification.
type VerifyTOTPResponse struct {
	Valid       bool      // Always true when verification succeeds
	Counter     uint64    // TOTP time step that matched
	ValidatedAt time.Time // Time used for verification
}

// GenerateRecoveryCodesRequest represents a recovery code generation.
type GenerateRecoveryCodesRequest struct {
	Count int // Number of codes to generate (optional, default: configured count)
}

// GenerateRecoveryCodesResponse contains the plaintext codes and their hashes.
type GenerateRecoveryCodesResponse struct {
	Codes  []string // Plaintext codes, shown to the user once
	Hashes []string // Hashes to store for later verification
}

// VerifyRecoveryCodeRequest represents a recovery code verification.
type VerifyRecoveryCodeRequest struct {
	UserID string   // Stable identifier of the account (required)
	Code   string   // Code entered by the user (required)
	Hashes []string // Stored hashes for the account (required)
}

// Validate checks that the request has the fields required to verify a recovery code.
func (r VerifyRecoveryCodeRequest) Validate() error {
	if strings.TrimSpace(r.UserID) == "" {
		return ErrMissingUserID
	}

	if strings.TrimSpace(r.Code) == "" {
		return ErrMissingCode
	}

	if len(r.Hashes) == 0 {
		return ErrInvalidRecoveryCode
	}

	return nil
}

// VerifyRecoveryCodeResponse is returned after a successful recovery code verification.
type VerifyRecoveryCodeResponse struct {
	Valid           bool     // Always true when verification succeeds
	RemainingHashes []string // Hashes to persist after removing the used code
}
