package mfa

import (
	"bytes"
	"context"
	"image/png"
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTOTPTestManager(t *testing.T, now time.Time) *Manager {
	t.Helper()

	manager, err := New(Config{
		Issuer: "Tixia",
		Store:  NewMemoryStore(),
		Now:    func() time.Time { return now },
	}, testLogger())
	require.NoError(t, err)

	return manager
}

func generateTOTPCode(t *testing.T, secret string, at time.Time, digits otp.Digits) string {
	t.Helper()

	code, err := totp.GenerateCodeCustom(secret, at, totp.ValidateOpts{
		Period:    30,
		Skew:      0,
		Digits:    digits,
		Algorithm: otp.AlgorithmSHA1,
	})
	require.NoError(t, err)

	return code
}

func TestEnrollTOTP(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	manager := newTOTPTestManager(t, now)

	response, err := manager.EnrollTOTP(context.Background(), EnrollTOTPRequest{
		AccountName: "user@example.com",
	})
	require.NoError(t, err)

	assert.NotEmpty(t, response.Secret)
	assert.Contains(t, response.URL, "otpauth://totp/")
	assert.Contains(t, response.URL, "issuer=Tixia")

	_, err = png.Decode(bytes.NewReader(response.QRCodePNG))
	assert.NoError(t, err)
}

func TestEnrollTOTP_Validation(t *testing.T) {
	manager := newTOTPTestManager(t, time.Now())

	_, err := manager.EnrollTOTP(context.Background(), EnrollTOTPRequest{})
	assert.ErrorIs(t, err, ErrMissingAccountName)
}

func TestVerifyTOTP_Valid(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	manager := newTOTPTestManager(t, now)

	enrollment, err := manager.EnrollTOTP(context.Background(), EnrollTOTPRequest{
		AccountName: "user@example.com",
	})
	require.NoError(t, err)

	code := generateTOTPCode(t, enrollment.Secret, now, otp.DigitsSix)

	response, err := manager.VerifyTOTP(context.Background(), VerifyTOTPRequest{
		UserID: "user-1",
		Secret: enrollment.Secret,
		Code:   code,
	})
	require.NoError(t, err)

	assert.True(t, response.Valid)
	assert.Equal(t, uint64(now.Unix())/30, response.Counter)
	assert.Equal(t, now, response.ValidatedAt)
}

func TestVerifyTOTP_AcceptsFormattedCode(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	manager := newTOTPTestManager(t, now)

	enrollment, err := manager.EnrollTOTP(context.Background(), EnrollTOTPRequest{
		AccountName: "user@example.com",
	})
	require.NoError(t, err)

	code := generateTOTPCode(t, enrollment.Secret, now, otp.DigitsSix)
	formatted := code[:3] + " " + code[3:]

	_, err = manager.VerifyTOTP(context.Background(), VerifyTOTPRequest{
		UserID: "user-1",
		Secret: enrollment.Secret,
		Code:   formatted,
	})
	assert.NoError(t, err)
}

func TestVerifyTOTP_RejectsReplay(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	manager := newTOTPTestManager(t, now)

	enrollment, err := manager.EnrollTOTP(context.Background(), EnrollTOTPRequest{
		AccountName: "user@example.com",
	})
	require.NoError(t, err)

	code := generateTOTPCode(t, enrollment.Secret, now, otp.DigitsSix)
	request := VerifyTOTPRequest{
		UserID: "user-1",
		Secret: enrollment.Secret,
		Code:   code,
	}

	_, err = manager.VerifyTOTP(context.Background(), request)
	require.NoError(t, err)

	_, err = manager.VerifyTOTP(context.Background(), request)
	assert.ErrorIs(t, err, ErrCodeReused)
}

func TestVerifyTOTP_RejectsInvalidCode(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	manager := newTOTPTestManager(t, now)

	enrollment, err := manager.EnrollTOTP(context.Background(), EnrollTOTPRequest{
		AccountName: "user@example.com",
	})
	require.NoError(t, err)

	code := generateTOTPCode(t, enrollment.Secret, now, otp.DigitsSix)
	invalid := "000000"
	if invalid == code {
		invalid = "000001"
	}

	_, err = manager.VerifyTOTP(context.Background(), VerifyTOTPRequest{
		UserID: "user-1",
		Secret: enrollment.Secret,
		Code:   invalid,
	})
	assert.ErrorIs(t, err, ErrInvalidCode)
}

func TestVerifyTOTP_Skew(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()

	t.Run("accepts previous period", func(t *testing.T) {
		manager := newTOTPTestManager(t, now)

		enrollment, err := manager.EnrollTOTP(context.Background(), EnrollTOTPRequest{
			AccountName: "user@example.com",
		})
		require.NoError(t, err)

		code := generateTOTPCode(t, enrollment.Secret, now.Add(30*time.Second), otp.DigitsSix)

		_, err = manager.VerifyTOTP(context.Background(), VerifyTOTPRequest{
			UserID: "user-1",
			Secret: enrollment.Secret,
			Code:   code,
		})
		assert.NoError(t, err)
	})

	t.Run("rejects code outside skew", func(t *testing.T) {
		manager := newTOTPTestManager(t, now)

		enrollment, err := manager.EnrollTOTP(context.Background(), EnrollTOTPRequest{
			AccountName: "user@example.com",
		})
		require.NoError(t, err)

		code := generateTOTPCode(t, enrollment.Secret, now.Add(90*time.Second), otp.DigitsSix)

		_, err = manager.VerifyTOTP(context.Background(), VerifyTOTPRequest{
			UserID: "user-1",
			Secret: enrollment.Secret,
			Code:   code,
		})
		assert.ErrorIs(t, err, ErrInvalidCode)
	})
}

func TestVerifyTOTP_EightDigits(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()

	manager, err := New(Config{
		Issuer: "Tixia",
		Digits: otp.DigitsEight,
		Store:  NewMemoryStore(),
		Now:    func() time.Time { return now },
	}, testLogger())
	require.NoError(t, err)

	enrollment, err := manager.EnrollTOTP(context.Background(), EnrollTOTPRequest{
		AccountName: "user@example.com",
	})
	require.NoError(t, err)

	code := generateTOTPCode(t, enrollment.Secret, now, otp.DigitsEight)

	_, err = manager.VerifyTOTP(context.Background(), VerifyTOTPRequest{
		UserID: "user-1",
		Secret: enrollment.Secret,
		Code:   code,
	})
	assert.NoError(t, err)
}

func TestVerifyTOTP_RequestValidation(t *testing.T) {
	manager := newTOTPTestManager(t, time.Now())

	tests := []struct {
		name    string
		request VerifyTOTPRequest
		err     error
	}{
		{
			name:    "missing user ID",
			request: VerifyTOTPRequest{Secret: "secret", Code: "123456"},
			err:     ErrMissingUserID,
		},
		{
			name:    "missing secret",
			request: VerifyTOTPRequest{UserID: "user-1", Code: "123456"},
			err:     ErrMissingSecret,
		},
		{
			name:    "missing code",
			request: VerifyTOTPRequest{UserID: "user-1", Secret: "secret"},
			err:     ErrMissingCode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := manager.VerifyTOTP(context.Background(), test.request)
			assert.ErrorIs(t, err, test.err)
		})
	}
}

func TestVerifyTOTP_InvalidSecret(t *testing.T) {
	manager := newTOTPTestManager(t, time.Now())

	_, err := manager.VerifyTOTP(context.Background(), VerifyTOTPRequest{
		UserID: "user-1",
		Secret: "not-valid-base32!",
		Code:   "123456",
	})
	assert.Error(t, err)
	assert.NotErrorIs(t, err, ErrInvalidCode)
}
