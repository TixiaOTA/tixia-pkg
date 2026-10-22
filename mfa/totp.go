package mfa

import (
	"bytes"
	"context"
	"crypto/subtle"
	"fmt"
	"image/png"
	"time"

	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"

	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
)

// EnrollTOTP generates a TOTP secret and the data required to add it to an
// authenticator app.
func (m *Manager) EnrollTOTP(ctx context.Context, req EnrollTOTPRequest) (*EnrollTOTPResponse, error) {
	ctx, span, spanID := otelpkg.StartSpanWithID(ctx, "utility", "mfa", "EnrollTOTP")
	defer span.End()

	if err := req.Validate(); err != nil {
		m.logger.WithTrace(ctx, spanID).Error("invalid TOTP enrollment request", zap.Error(err))
		return nil, err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      m.config.Issuer,
		AccountName: req.AccountName,
		Period:      m.config.Period,
		SecretSize:  m.config.SecretSize,
		Digits:      m.config.Digits,
		Algorithm:   m.config.Algorithm,
	})
	if err != nil {
		return nil, fmt.Errorf("mfa: generate TOTP key: %w", err)
	}

	image, err := key.Image(m.config.QRCodeSize, m.config.QRCodeSize)
	if err != nil {
		return nil, fmt.Errorf("mfa: generate QR code: %w", err)
	}

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, image); err != nil {
		return nil, fmt.Errorf("mfa: encode QR code: %w", err)
	}

	m.logger.WithTrace(ctx, spanID).Info(
		"TOTP enrollment created",
		zap.String("account_name", req.AccountName),
	)

	return &EnrollTOTPResponse{
		Secret:    key.Secret(),
		URL:       key.URL(),
		QRCodePNG: buffer.Bytes(),
	}, nil
}

// VerifyTOTP verifies a TOTP code and records it so it cannot be replayed.
func (m *Manager) VerifyTOTP(ctx context.Context, req VerifyTOTPRequest) (*VerifyTOTPResponse, error) {
	ctx, span, spanID := otelpkg.StartSpanWithID(ctx, "utility", "mfa", "VerifyTOTP")
	defer span.End()

	if err := req.Validate(); err != nil {
		m.logger.WithTrace(ctx, spanID).Error("invalid TOTP verification request", zap.Error(err))
		return nil, err
	}

	now := m.now()

	counter, matched, err := m.matchTOTPCounter(req.Secret, req.Code, now)
	if err != nil {
		return nil, err
	}

	if !matched {
		m.logger.WithTrace(ctx, spanID).Warn(
			"invalid TOTP code",
			zap.String("user_id", req.UserID),
		)
		return nil, ErrInvalidCode
	}

	key := fmt.Sprintf("totp:used:%s:%d", req.UserID, counter)

	created, err := m.config.Store.MarkCodeUsed(ctx, key, m.replayTTL())
	if err != nil {
		return nil, fmt.Errorf("mfa: record used TOTP code: %w", err)
	}

	if !created {
		m.logger.WithTrace(ctx, spanID).Warn(
			"TOTP code replay rejected",
			zap.String("user_id", req.UserID),
			zap.Uint64("counter", counter),
		)
		return nil, ErrCodeReused
	}

	return &VerifyTOTPResponse{
		Valid:       true,
		Counter:     counter,
		ValidatedAt: now,
	}, nil
}

func (m *Manager) matchTOTPCounter(secret, code string, now time.Time) (uint64, bool, error) {
	code = normalizeCode(code)

	if len(code) != m.config.Digits.Length() {
		return 0, false, nil
	}

	period := uint64(m.config.Period)
	counter := uint64(now.Unix()) / period
	skew := uint64(m.config.skew())

	options := totp.ValidateOpts{
		Period:    m.config.Period,
		Skew:      0,
		Digits:    m.config.Digits,
		Algorithm: m.config.Algorithm,
	}

	for offset := -int64(skew); offset <= int64(skew); offset++ {
		if offset < 0 && uint64(-offset) > counter {
			continue
		}

		candidateCounter := counter + uint64(offset)
		candidateTime := time.Unix(int64(candidateCounter*period), 0).UTC()

		candidate, err := totp.GenerateCodeCustom(secret, candidateTime, options)
		if err != nil {
			return 0, false, fmt.Errorf("mfa: generate TOTP code: %w", err)
		}

		if subtle.ConstantTimeCompare([]byte(candidate), []byte(code)) == 1 {
			return candidateCounter, true, nil
		}
	}

	return 0, false, nil
}
