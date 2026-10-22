package mfa

import (
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TixiaOTA/tixia-pkg/logger"
)

func testLogger() *logger.Logger {
	return logger.New(logger.Config{
		Level:       "error",
		Environment: "development",
	})
}

func TestNew_AppliesDefaults(t *testing.T) {
	manager, err := New(Config{
		Issuer: "Tixia",
		Store:  NewMemoryStore(),
	}, testLogger())
	require.NoError(t, err)

	assert.Equal(t, uint(30), manager.config.Period)
	assert.Equal(t, otp.DigitsSix, manager.config.Digits)
	assert.Equal(t, otp.AlgorithmSHA1, manager.config.Algorithm)
	assert.Equal(t, uint(1), manager.config.skew())
	assert.Equal(t, uint(20), manager.config.SecretSize)
	assert.Equal(t, 10, manager.config.RecoveryCodeCount)
	assert.Equal(t, 10, manager.config.RecoveryCodeLength)
	assert.Equal(t, 256, manager.config.QRCodeSize)
	assert.NotNil(t, manager.config.Now)
}

func TestNew_ValidationErrors(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		err    error
	}{
		{
			name:   "missing issuer",
			config: Config{Store: NewMemoryStore()},
			err:    ErrMissingIssuer,
		},
		{
			name:   "nil store",
			config: Config{Issuer: "Tixia"},
			err:    ErrNilStore,
		},
		{
			name:   "invalid digits",
			config: Config{Issuer: "Tixia", Store: NewMemoryStore(), Digits: otp.Digits(4)},
			err:    ErrInvalidDigits,
		},
		{
			name:   "invalid algorithm",
			config: Config{Issuer: "Tixia", Store: NewMemoryStore(), Algorithm: otp.AlgorithmMD5},
			err:    ErrInvalidAlgorithm,
		},
		{
			name:   "invalid skew",
			config: Config{Issuer: "Tixia", Store: NewMemoryStore(), Skew: uintPtr(9)},
			err:    ErrInvalidSkew,
		},
		{
			name:   "invalid recovery code count",
			config: Config{Issuer: "Tixia", Store: NewMemoryStore(), RecoveryCodeCount: -1},
			err:    ErrInvalidRecoveryCodeCount,
		},
		{
			name:   "invalid recovery code length",
			config: Config{Issuer: "Tixia", Store: NewMemoryStore(), RecoveryCodeLength: 4},
			err:    ErrInvalidRecoveryCodeLength,
		},
		{
			name:   "negative QR code size",
			config: Config{Issuer: "Tixia", Store: NewMemoryStore(), QRCodeSize: -1},
			err:    ErrInvalidQRCodeSize,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.config, testLogger())
			assert.ErrorIs(t, err, test.err)
		})
	}
}

func TestNew_NilLogger(t *testing.T) {
	_, err := New(Config{
		Issuer: "Tixia",
		Store:  NewMemoryStore(),
	}, nil)

	assert.ErrorIs(t, err, ErrNilLogger)
}

func TestConfig_ZeroSkewIsPreserved(t *testing.T) {
	manager, err := New(Config{
		Issuer: "Tixia",
		Store:  NewMemoryStore(),
		Skew:   uintPtr(0),
		Now:    func() time.Time { return time.Unix(1700000000, 0).UTC() },
	}, testLogger())
	require.NoError(t, err)

	assert.Equal(t, uint(0), manager.config.skew())
}

func uintPtr(value uint) *uint {
	return &value
}
