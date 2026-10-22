package mfa

import "errors"

var (
	ErrMissingIssuer             = errors.New("mfa: issuer is required")
	ErrNilLogger                 = errors.New("mfa: logger is required")
	ErrNilStore                  = errors.New("mfa: store is required")
	ErrInvalidPeriod             = errors.New("mfa: period must be greater than zero")
	ErrInvalidDigits             = errors.New("mfa: digits must be 6 or 8")
	ErrInvalidAlgorithm          = errors.New("mfa: algorithm must be SHA1, SHA256, or SHA512")
	ErrInvalidSkew               = errors.New("mfa: skew must be between 0 and 5")
	ErrInvalidSecretSize         = errors.New("mfa: secret size must be greater than zero")
	ErrInvalidQRCodeSize         = errors.New("mfa: QR code size must not be negative")
	ErrMissingAccountName        = errors.New("mfa: account name is required")
	ErrMissingUserID             = errors.New("mfa: user ID is required")
	ErrMissingSecret             = errors.New("mfa: secret is required")
	ErrMissingCode               = errors.New("mfa: code is required")
	ErrInvalidCode               = errors.New("mfa: invalid code")
	ErrCodeReused                = errors.New("mfa: code has already been used")
	ErrMissingRecoveryCodePepper = errors.New("mfa: recovery code pepper is required")
	ErrWeakRecoveryCodePepper    = errors.New("mfa: recovery code pepper must be at least 32 characters")
	ErrInvalidRecoveryCodeCount  = errors.New("mfa: recovery code count must be between 1 and 100")
	ErrInvalidRecoveryCodeLength = errors.New("mfa: recovery code length must be between 8 and 32")
	ErrInvalidRecoveryCode       = errors.New("mfa: invalid recovery code")
)
