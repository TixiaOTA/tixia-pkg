package mfa

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newRecoveryTestManager(t *testing.T) *Manager {
	t.Helper()

	manager, err := New(Config{
		Issuer:             "Tixia",
		Store:              NewMemoryStore(),
		RecoveryCodePepper: strings.Repeat("p", 32),
		Now:                func() time.Time { return time.Unix(1700000000, 0).UTC() },
	}, testLogger())
	require.NoError(t, err)

	return manager
}

func TestGenerateRecoveryCodes(t *testing.T) {
	manager := newRecoveryTestManager(t)

	response, err := manager.GenerateRecoveryCodes(context.Background(), GenerateRecoveryCodesRequest{
		Count: 5,
	})
	require.NoError(t, err)

	assert.Len(t, response.Codes, 5)
	assert.Len(t, response.Hashes, 5)

	for index, code := range response.Codes {
		assert.Contains(t, code, "-")
		assert.Len(t, normalizeCode(code), manager.config.RecoveryCodeLength)
		assert.Len(t, response.Hashes[index], 64)
	}
}

func TestGenerateRecoveryCodes_UsesConfiguredCount(t *testing.T) {
	manager := newRecoveryTestManager(t)

	response, err := manager.GenerateRecoveryCodes(context.Background(), GenerateRecoveryCodesRequest{})
	require.NoError(t, err)

	assert.Len(t, response.Codes, manager.config.RecoveryCodeCount)
}

func TestVerifyRecoveryCode(t *testing.T) {
	manager := newRecoveryTestManager(t)

	generated, err := manager.GenerateRecoveryCodes(context.Background(), GenerateRecoveryCodesRequest{
		Count: 3,
	})
	require.NoError(t, err)

	response, err := manager.VerifyRecoveryCode(context.Background(), VerifyRecoveryCodeRequest{
		UserID: "user-1",
		Code:   generated.Codes[1],
		Hashes: generated.Hashes,
	})
	require.NoError(t, err)

	assert.True(t, response.Valid)
	assert.Len(t, response.RemainingHashes, 2)
	assert.NotContains(t, response.RemainingHashes, generated.Hashes[1])
}

func TestVerifyRecoveryCode_NormalizesInput(t *testing.T) {
	manager := newRecoveryTestManager(t)

	generated, err := manager.GenerateRecoveryCodes(context.Background(), GenerateRecoveryCodesRequest{
		Count: 1,
	})
	require.NoError(t, err)

	normalized := strings.ToLower(strings.ReplaceAll(generated.Codes[0], "-", " "))

	_, err = manager.VerifyRecoveryCode(context.Background(), VerifyRecoveryCodeRequest{
		UserID: "user-1",
		Code:   normalized,
		Hashes: generated.Hashes,
	})
	assert.NoError(t, err)
}

func TestVerifyRecoveryCode_RejectsUsedCode(t *testing.T) {
	manager := newRecoveryTestManager(t)

	generated, err := manager.GenerateRecoveryCodes(context.Background(), GenerateRecoveryCodesRequest{
		Count: 2,
	})
	require.NoError(t, err)

	response, err := manager.VerifyRecoveryCode(context.Background(), VerifyRecoveryCodeRequest{
		UserID: "user-1",
		Code:   generated.Codes[0],
		Hashes: generated.Hashes,
	})
	require.NoError(t, err)
	assert.Len(t, response.RemainingHashes, 1)

	_, err = manager.VerifyRecoveryCode(context.Background(), VerifyRecoveryCodeRequest{
		UserID: "user-1",
		Code:   generated.Codes[0],
		Hashes: response.RemainingHashes,
	})
	assert.ErrorIs(t, err, ErrInvalidRecoveryCode)
}

func TestVerifyRecoveryCode_Invalid(t *testing.T) {
	manager := newRecoveryTestManager(t)

	generated, err := manager.GenerateRecoveryCodes(context.Background(), GenerateRecoveryCodesRequest{
		Count: 1,
	})
	require.NoError(t, err)

	_, err = manager.VerifyRecoveryCode(context.Background(), VerifyRecoveryCodeRequest{
		UserID: "user-1",
		Code:   "AAAAA-AAAAA",
		Hashes: generated.Hashes,
	})
	assert.ErrorIs(t, err, ErrInvalidRecoveryCode)
}

func TestRecoveryCode_PepperValidation(t *testing.T) {
	t.Run("missing pepper", func(t *testing.T) {
		manager, err := New(Config{
			Issuer: "Tixia",
			Store:  NewMemoryStore(),
		}, testLogger())
		require.NoError(t, err)

		_, err = manager.GenerateRecoveryCodes(context.Background(), GenerateRecoveryCodesRequest{})
		assert.ErrorIs(t, err, ErrMissingRecoveryCodePepper)
	})

	t.Run("weak pepper", func(t *testing.T) {
		manager, err := New(Config{
			Issuer:             "Tixia",
			Store:              NewMemoryStore(),
			RecoveryCodePepper: "short",
		}, testLogger())
		require.NoError(t, err)

		_, err = manager.GenerateRecoveryCodes(context.Background(), GenerateRecoveryCodesRequest{})
		assert.ErrorIs(t, err, ErrWeakRecoveryCodePepper)
	})

	t.Run("invalid count", func(t *testing.T) {
		manager := newRecoveryTestManager(t)

		_, err := manager.GenerateRecoveryCodes(context.Background(), GenerateRecoveryCodesRequest{
			Count: 101,
		})
		assert.ErrorIs(t, err, ErrInvalidRecoveryCodeCount)
	})
}
