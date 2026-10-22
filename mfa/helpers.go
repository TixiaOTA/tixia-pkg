package mfa

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strings"
	"unicode"
)

const (
	recoveryAlphabet  = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	recoveryGroupSize = 5
)

func normalizeCode(code string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || r == '-' {
			return -1
		}

		return r
	}, code)
}

func generateRecoveryCode(length int) (string, error) {
	alphabetSize := big.NewInt(int64(len(recoveryAlphabet)))
	var builder strings.Builder

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, alphabetSize)
		if err != nil {
			return "", err
		}

		builder.WriteByte(recoveryAlphabet[index.Int64()])
	}

	return formatRecoveryCode(builder.String()), nil
}

func formatRecoveryCode(code string) string {
	var builder strings.Builder

	for i, character := range code {
		if i > 0 && i%recoveryGroupSize == 0 {
			builder.WriteByte('-')
		}

		builder.WriteRune(character)
	}

	return builder.String()
}

func hashRecoveryCode(pepper, code string) string {
	mac := hmac.New(sha256.New, []byte(pepper))
	mac.Write([]byte(strings.ToUpper(normalizeCode(code))))

	return hex.EncodeToString(mac.Sum(nil))
}
