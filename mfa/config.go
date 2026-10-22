package mfa

import (
	"strings"
	"time"

	"github.com/pquerna/otp"
)

const (
	defaultPeriod             = 30
	defaultSkew               = 1
	defaultSecretSize         = 20
	defaultRecoveryCodeCount  = 10
	defaultRecoveryCodeLength = 10
	defaultQRCodeSize         = 256

	maxSkew                   = 5
	maxRecoveryCodeCount      = 100
	minRecoveryCodeLength     = 8
	maxRecoveryCodeLength     = 32
	minRecoveryCodePepperSize = 32
)

// Config represents the MFA manager configuration.
type Config struct {
	Issuer             string           // Issuer shown in authenticator apps (required)
	Period             uint             // TOTP period in seconds (optional, default: 30)
	Digits             otp.Digits       // TOTP code length (optional, default: 6)
	Algorithm          otp.Algorithm    // HMAC algorithm (optional, default: SHA1)
	Skew               *uint            // Accepted TOTP periods before and after now (optional, default: 1)
	SecretSize         uint             // TOTP secret size in bytes (optional, default: 20)
	RecoveryCodeCount  int              // Recovery codes per set (optional, default: 10)
	RecoveryCodeLength int              // Characters per recovery code, excluding separators (optional, default: 10)
	RecoveryCodePepper string           // Secret used to hash recovery codes (required for recovery codes, min 32 characters)
	QRCodeSize         int              // QR code image size in pixels (optional, default: 256)
	Store              Store            // Replay protection store (required)
	Now                func() time.Time // Clock used for TOTP validation (optional, default: time.Now)
}

// Validate checks that all required configuration fields are set.
func (c Config) Validate() error {
	c.applyDefaults()

	if strings.TrimSpace(c.Issuer) == "" {
		return ErrMissingIssuer
	}

	if c.Store == nil {
		return ErrNilStore
	}

	if c.Period == 0 {
		return ErrInvalidPeriod
	}

	if c.Digits != otp.DigitsSix && c.Digits != otp.DigitsEight {
		return ErrInvalidDigits
	}

	if c.Algorithm != otp.AlgorithmSHA1 && c.Algorithm != otp.AlgorithmSHA256 && c.Algorithm != otp.AlgorithmSHA512 {
		return ErrInvalidAlgorithm
	}

	if c.Skew != nil && *c.Skew > maxSkew {
		return ErrInvalidSkew
	}

	if c.SecretSize == 0 {
		return ErrInvalidSecretSize
	}

	if c.QRCodeSize < 0 {
		return ErrInvalidQRCodeSize
	}

	if c.RecoveryCodeCount < 1 || c.RecoveryCodeCount > maxRecoveryCodeCount {
		return ErrInvalidRecoveryCodeCount
	}

	if c.RecoveryCodeLength < minRecoveryCodeLength || c.RecoveryCodeLength > maxRecoveryCodeLength {
		return ErrInvalidRecoveryCodeLength
	}

	return nil
}

func (c *Config) applyDefaults() {
	if c.Period == 0 {
		c.Period = defaultPeriod
	}

	if c.Digits == 0 {
		c.Digits = otp.DigitsSix
	}

	if c.Skew == nil {
		skew := uint(defaultSkew)
		c.Skew = &skew
	}

	if c.SecretSize == 0 {
		c.SecretSize = defaultSecretSize
	}

	if c.RecoveryCodeCount == 0 {
		c.RecoveryCodeCount = defaultRecoveryCodeCount
	}

	if c.RecoveryCodeLength == 0 {
		c.RecoveryCodeLength = defaultRecoveryCodeLength
	}

	if c.QRCodeSize == 0 {
		c.QRCodeSize = defaultQRCodeSize
	}

	if c.Now == nil {
		c.Now = time.Now
	}
}

func (c Config) skew() uint {
	if c.Skew == nil {
		return defaultSkew
	}

	return *c.Skew
}
