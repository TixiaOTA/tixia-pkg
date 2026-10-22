package mfa

import (
	"context"
	"crypto/subtle"
	"fmt"
	"strings"

	"go.uber.org/zap"

	otelpkg "github.com/TixiaOTA/tixia-pkg/otel"
)

// GenerateRecoveryCodes generates one-time recovery codes and the hashes that
// must be stored for later verification. The plaintext codes are returned only
// once.
func (m *Manager) GenerateRecoveryCodes(ctx context.Context, req GenerateRecoveryCodesRequest) (*GenerateRecoveryCodesResponse, error) {
	ctx, span, spanID := otelpkg.StartSpanWithID(ctx, "utility", "mfa", "GenerateRecoveryCodes")
	defer span.End()

	if err := m.validateRecoveryConfig(); err != nil {
		m.logger.WithTrace(ctx, spanID).Error("invalid recovery code configuration", zap.Error(err))
		return nil, err
	}

	count := req.Count
	if count == 0 {
		count = m.config.RecoveryCodeCount
	}

	if count < 1 || count > maxRecoveryCodeCount {
		return nil, ErrInvalidRecoveryCodeCount
	}

	codes := make([]string, 0, count)
	hashes := make([]string, 0, count)

	for i := 0; i < count; i++ {
		code, err := generateRecoveryCode(m.config.RecoveryCodeLength)
		if err != nil {
			return nil, fmt.Errorf("mfa: generate recovery code: %w", err)
		}

		codes = append(codes, code)
		hashes = append(hashes, hashRecoveryCode(m.config.RecoveryCodePepper, code))
	}

	m.logger.WithTrace(ctx, spanID).Info(
		"recovery codes generated",
		zap.Int("count", count),
	)

	return &GenerateRecoveryCodesResponse{
		Codes:  codes,
		Hashes: hashes,
	}, nil
}

// VerifyRecoveryCode verifies a recovery code against the stored hashes and
// returns the remaining hashes. The caller must persist the remaining hashes
// atomically to keep the code single-use.
func (m *Manager) VerifyRecoveryCode(ctx context.Context, req VerifyRecoveryCodeRequest) (*VerifyRecoveryCodeResponse, error) {
	ctx, span, spanID := otelpkg.StartSpanWithID(ctx, "utility", "mfa", "VerifyRecoveryCode")
	defer span.End()

	if err := req.Validate(); err != nil {
		m.logger.WithTrace(ctx, spanID).Error("invalid recovery code verification request", zap.Error(err))
		return nil, err
	}

	if err := m.validateRecoveryConfig(); err != nil {
		m.logger.WithTrace(ctx, spanID).Error("invalid recovery code configuration", zap.Error(err))
		return nil, err
	}

	digest := hashRecoveryCode(m.config.RecoveryCodePepper, req.Code)

	for index, hash := range req.Hashes {
		if subtle.ConstantTimeCompare([]byte(hash), []byte(digest)) != 1 {
			continue
		}

		remaining := make([]string, 0, len(req.Hashes)-1)
		remaining = append(remaining, req.Hashes[:index]...)
		remaining = append(remaining, req.Hashes[index+1:]...)

		m.logger.WithTrace(ctx, spanID).Info(
			"recovery code verified",
			zap.String("user_id", req.UserID),
			zap.Int("remaining_codes", len(remaining)),
		)

		return &VerifyRecoveryCodeResponse{
			Valid:           true,
			RemainingHashes: remaining,
		}, nil
	}

	m.logger.WithTrace(ctx, spanID).Warn(
		"invalid recovery code",
		zap.String("user_id", req.UserID),
	)

	return nil, ErrInvalidRecoveryCode
}

func (m *Manager) validateRecoveryConfig() error {
	if strings.TrimSpace(m.config.RecoveryCodePepper) == "" {
		return ErrMissingRecoveryCodePepper
	}

	if len(m.config.RecoveryCodePepper) < minRecoveryCodePepperSize {
		return ErrWeakRecoveryCodePepper
	}

	return nil
}
