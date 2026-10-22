# MFA Package

A Go library for time-based one-time passwords (TOTP) and recovery codes. The package generates secrets that work with Google Authenticator, Authy, 1Password, Microsoft Authenticator, and any other RFC 6238 compatible app.

## Overview

The MFA package enables you to:

- Enroll accounts with authenticator apps using an `otpauth://` URL and QR code
- Verify TOTP codes with configurable period, digits, algorithm, and clock skew
- Prevent code replay with a pluggable store
- Generate and verify single-use recovery codes
- Use the included Redis-backed store or provide your own `Store`
- Trace operations with OpenTelemetry and structured logging

The package is stateless. The TOTP secret and recovery code hashes belong in your account storage, and the `Store` only holds short-lived replay markers.

## Installation

```bash
go get github.com/TixiaOTA/tixia-pkg/mfa
```

For the Redis-backed replay store:

```bash
go get github.com/TixiaOTA/tixia-pkg/mfa/store/redis
```

## Configuration

```go
import (
    "time"

    "github.com/TixiaOTA/tixia-pkg/logger"
    "github.com/TixiaOTA/tixia-pkg/mfa"
)

manager, err := mfa.New(mfa.Config{
    Issuer: "Tixia",
    Store:  mfa.NewMemoryStore(),
}, logger.Default())
if err != nil {
    log.Fatal(err)
}
```

For multiple instances, use the Redis-backed store so a code used on one instance is rejected by all of them:

```go
import (
    "os"
    "time"

    "github.com/TixiaOTA/tixia-pkg/logger"
    "github.com/TixiaOTA/tixia-pkg/mfa"
    redisstore "github.com/TixiaOTA/tixia-pkg/mfa/store/redis"
    "github.com/TixiaOTA/tixia-pkg/redis/cache"
    redisclient "github.com/TixiaOTA/tixia-pkg/redis/client"
)

cacheManager, err := cache.NewCacheManager(cache.CacheConfig{
    BaseRedisConfig: redisclient.BaseRedisConfig{
        URL: "redis://localhost:6379",
    },
    DefaultTTL: time.Minute,
    MaxTTL:     24 * time.Hour,
    KeyPrefix:  "tixia:",
}, logger.Default())
if err != nil {
    log.Fatal(err)
}

store, err := redisstore.New(cacheManager)
if err != nil {
    log.Fatal(err)
}

manager, err := mfa.New(mfa.Config{
    Issuer:             "Tixia",
    Store:              store,
    RecoveryCodePepper: os.Getenv("MFA_RECOVERY_CODE_PEPPER"),
}, logger.Default())
if err != nil {
    log.Fatal(err)
}
```

### Configuration Fields

- **Issuer** (string, required): Organization name shown in the authenticator app.
- **Period** (uint, optional): TOTP period in seconds. Defaults to `30`.
- **Digits** (otp.Digits, optional): Code length, `otp.DigitsSix` or `otp.DigitsEight`. Defaults to `otp.DigitsSix`.
- **Algorithm** (otp.Algorithm, optional): HMAC algorithm, `otp.AlgorithmSHA1`, `otp.AlgorithmSHA256`, or `otp.AlgorithmSHA512`. Defaults to `otp.AlgorithmSHA1`, which is required by most authenticator apps.
- **Skew** (*uint, optional): Periods accepted before and after the current time. Defaults to `1` (plus or minus 30 seconds). A pointer is used so `0` (no skew) is distinguishable from the default.
- **SecretSize** (uint, optional): TOTP secret size in bytes. Defaults to `20` (160 bits).
- **RecoveryCodeCount** (int, optional): Codes generated per set. Defaults to `10`. Must be between `1` and `100`.
- **RecoveryCodeLength** (int, optional): Characters per code, excluding separators. Defaults to `10`. Must be between `8` and `32`.
- **RecoveryCodePepper** (string, optional): Secret used to hash recovery codes. Required by the recovery code methods and must be at least 32 characters. Store it outside the database.
- **QRCodeSize** (int, optional): QR code image size in pixels. Defaults to `256`.
- **Store** (Store, required): Replay protection store.
- **Now** (func() time.Time, optional): Clock used for TOTP validation. Defaults to `time.Now`.

## Usage

### Enroll an Authenticator App

Enrollment returns the base32 secret to store, the `otpauth://` URL, and a PNG QR code. Show the QR code or the secret to the user once.

```go
enrollment, err := manager.EnrollTOTP(ctx, mfa.EnrollTOTPRequest{
    AccountName: "user@example.com",
})
if err != nil {
    log.Fatal(err)
}

// Persist the secret with the account.
account.TOTPSecret = enrollment.Secret

// Render the QR code in your enrollment screen.
w.Header().Set("Content-Type", "image/png")
w.Write(enrollment.QRCodePNG)
```

### Verify a TOTP Code

```go
_, err := manager.VerifyTOTP(ctx, mfa.VerifyTOTPRequest{
    UserID: account.ID,
    Secret: account.TOTPSecret,
    Code:   "123456",
})
if err != nil {
    switch {
    case errors.Is(err, mfa.ErrInvalidCode):
        // Wrong code, ask the user to try again.
    case errors.Is(err, mfa.ErrCodeReused):
        // The code was already used within its validity window.
    default:
        // Store or configuration failure.
    }
}
```

A code can be verified only once. The matched TOTP counter is recorded in the store for the duration of the acceptance window.

### Generate Recovery Codes

```go
recovery, err := manager.GenerateRecoveryCodes(ctx, mfa.GenerateRecoveryCodesRequest{
    Count: 10,
})
if err != nil {
    log.Fatal(err)
}

// Show recovery.Codes to the user once.
// Store recovery.Hashes with the account.
account.RecoveryCodeHashes = recovery.Hashes
```

### Verify a Recovery Code

```go
result, err := manager.VerifyRecoveryCode(ctx, mfa.VerifyRecoveryCodeRequest{
    UserID: account.ID,
    Code:   "ABCDE-FGHJK",
    Hashes: account.RecoveryCodeHashes,
})
if err != nil {
    log.Fatal(err)
}

// Persist the remaining hashes atomically so the code cannot be used again.
account.RecoveryCodeHashes = result.RemainingHashes
```

Codes are case-insensitive and ignore spaces and dashes.

## API Reference

### New

Creates a new MFA manager and validates the configuration.

```go
func New(config Config, logger *logger.Logger) (*Manager, error)
```

**Parameters:**

- `config` (Config, required): Manager configuration
- `logger` (*logger.Logger, required): Logger instance

**Return Values:**

- `*Manager`: Manager instance
- `error`: Configuration or dependency error

### EnrollTOTP

Generates a TOTP secret and enrollment data.

```go
func (m *Manager) EnrollTOTP(ctx context.Context, req EnrollTOTPRequest) (*EnrollTOTPResponse, error)
```

**Request Parameters:**

- `AccountName` (string, required): Account name shown in the authenticator app

**Response:**

- `Secret` (string): Base32 secret to store
- `URL` (string): `otpauth://` URL
- `QRCodePNG` ([]byte): PNG-encoded QR code

### VerifyTOTP

Verifies a TOTP code and records it as used.

```go
func (m *Manager) VerifyTOTP(ctx context.Context, req VerifyTOTPRequest) (*VerifyTOTPResponse, error)
```

**Request Parameters:**

- `UserID` (string, required): Stable account identifier
- `Secret` (string, required): Base32 secret stored for the account
- `Code` (string, required): Code entered by the user

**Response:**

- `Valid` (bool): `true` when verification succeeds
- `Counter` (uint64): TOTP time step that matched
- `ValidatedAt` (time.Time): Time used for verification

### GenerateRecoveryCodes

Generates plaintext recovery codes and their hashes.

```go
func (m *Manager) GenerateRecoveryCodes(ctx context.Context, req GenerateRecoveryCodesRequest) (*GenerateRecoveryCodesResponse, error)
```

**Request Parameters:**

- `Count` (int, optional): Number of codes to generate. Uses the configured count when zero.

**Response:**

- `Codes` ([]string): Plaintext codes to show once
- `Hashes` ([]string): Hashes to store

### VerifyRecoveryCode

Verifies a recovery code and returns the hashes without the used code.

```go
func (m *Manager) VerifyRecoveryCode(ctx context.Context, req VerifyRecoveryCodeRequest) (*VerifyRecoveryCodeResponse, error)
```

**Request Parameters:**

- `UserID` (string, required): Stable account identifier
- `Code` (string, required): Code entered by the user
- `Hashes` ([]string, required): Stored hashes for the account

**Response:**

- `Valid` (bool): `true` when verification succeeds
- `RemainingHashes` ([]string): Hashes to persist after removing the used code

### Store

```go
type Store interface {
    MarkCodeUsed(ctx context.Context, key string, ttl time.Duration) (bool, error)
}
```

`MarkCodeUsed` must be atomic. It returns `true` when this call recorded the code and `false` when it was already recorded.

### MemoryStore

An in-process `Store` for tests and single-instance deployments.

```go
func NewMemoryStore() *MemoryStore
```

## Error Handling

The package returns sentinel errors that can be checked with `errors.Is`:

- `ErrMissingIssuer`: `Issuer` is empty.
- `ErrNilLogger`: no logger was provided.
- `ErrNilStore`: no store was provided.
- `ErrInvalidPeriod`: `Period` is zero.
- `ErrInvalidDigits`: `Digits` is not 6 or 8.
- `ErrInvalidAlgorithm`: `Algorithm` is not SHA1, SHA256, or SHA512.
- `ErrInvalidSkew`: `Skew` is greater than 5.
- `ErrInvalidSecretSize`: `SecretSize` is zero.
- `ErrInvalidQRCodeSize`: `QRCodeSize` is negative.
- `ErrMissingAccountName`: enrollment request has no account name.
- `ErrMissingUserID`: verification request has no user ID.
- `ErrMissingSecret`: verification request has no secret.
- `ErrMissingCode`: verification request has no code.
- `ErrInvalidCode`: the TOTP code is not valid for the current window.
- `ErrCodeReused`: the TOTP code was already used.
- `ErrMissingRecoveryCodePepper`: recovery code pepper is empty.
- `ErrWeakRecoveryCodePepper`: recovery code pepper is shorter than 32 characters.
- `ErrInvalidRecoveryCodeCount`: recovery code count is outside 1 to 100.
- `ErrInvalidRecoveryCodeLength`: recovery code length is outside 8 to 32.
- `ErrInvalidRecoveryCode`: the recovery code does not match any stored hash.

## Data Models

- `EnrollTOTPRequest`, `EnrollTOTPResponse`
- `VerifyTOTPRequest`, `VerifyTOTPResponse`
- `GenerateRecoveryCodesRequest`, `GenerateRecoveryCodesResponse`
- `VerifyRecoveryCodeRequest`, `VerifyRecoveryCodeResponse`

See `models.go` for complete definitions.

## Security Notes

- Store TOTP secrets encrypted at rest. The package returns the secret once and never persists it.
- Never log the secret or a code. The built-in logs include only account and user identifiers.
- Use a shared `Store` such as Redis when more than one instance verifies codes. `MemoryStore` only protects a single process.
- When using the Redis store, set the cache `MaxTTL` to at least the TOTP acceptance window (`period * (2 * skew + 2)`). `redis/cache` clamps longer TTLs to `MaxTTL`, and a shorter window would allow replay.
- Persist `RemainingHashes` in the same atomic write that marks the recovery code as used. Two concurrent verifications using the same stored hashes can otherwise both succeed.
- Keep `RecoveryCodePepper` in a secret manager, not in the database. Changing it invalidates all existing recovery codes.
- Apply rate limiting to verification endpoints. This package prevents replay but does not throttle guessing.

## Examples

### Complete Enrollment and Verification Flow

```go
package main

import (
    "context"
    "errors"
    "log"

    "github.com/TixiaOTA/tixia-pkg/logger"
    "github.com/TixiaOTA/tixia-pkg/mfa"
)

func main() {
    manager, err := mfa.New(mfa.Config{
        Issuer: "Tixia",
        Store:  mfa.NewMemoryStore(),
    }, logger.Default())
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    enrollment, err := manager.EnrollTOTP(ctx, mfa.EnrollTOTPRequest{
        AccountName: "user@example.com",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("secret: %s", enrollment.Secret)
    log.Printf("otpauth URL: %s", enrollment.URL)

    _, err = manager.VerifyTOTP(ctx, mfa.VerifyTOTPRequest{
        UserID: "user-1",
        Secret: enrollment.Secret,
        Code:   "123456",
    })
    if errors.Is(err, mfa.ErrInvalidCode) {
        log.Println("the code did not match")
    } else if err != nil {
        log.Fatal(err)
    }
}
```

## License

This package is part of the TixiaOTA project.
